package clients

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shashtag-ventures/go-common/integrations/types"
)

type GitHubClient struct {
	HTTPClient   *http.Client
	BaseURL      string
	ClientID     string
	ClientSecret string
	AppID        string
	PrivateKey   string
}

func NewGitHubClient(clientID, clientSecret string) *GitHubClient {
	return &GitHubClient{
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		BaseURL:      "https://api.github.com",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

func (c *GitHubClient) WithAppAuth(appID, privateKey string) *GitHubClient {
	c.AppID = appID
	c.PrivateKey = privateKey
	return c
}

func extractNextPageURL(linkHeader string) string {
	if linkHeader == "" {
		return ""
	}
	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) >= 2 && strings.Contains(parts[1], `rel="next"`) {
			urlPart := strings.TrimSpace(parts[0])
			if strings.HasPrefix(urlPart, "<") && strings.HasSuffix(urlPart, ">") {
				return urlPart[1 : len(urlPart)-1]
			}
		}
	}
	return ""
}


func (c *GitHubClient) generateJWT() (string, error) {
	if c.AppID == "" || c.PrivateKey == "" {
		return "", fmt.Errorf("GitHub App ID or Private Key not configured")
	}

	// Parse RSA Private Key
	block, _ := pem.Decode([]byte(c.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block from private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 if PKCS1 fails
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("failed to parse private key: %v", err)
		}
		var ok bool
		privKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return "", fmt.Errorf("not an RSA private key")
		}
	}

	// Create JWT
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    c.AppID,
		IssuedAt:  jwt.NewNumericDate(now.Add(-60 * time.Second)),
		ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privKey)
}

func (c *GitHubClient) GenerateInstallationToken(ctx context.Context, installationID string) (string, error) {
	if installationID == "" {
		return "", fmt.Errorf("installation ID is required")
	}
	var iid int64
	_, err := fmt.Sscanf(installationID, "%d", &iid)
	if err != nil {
		return "", fmt.Errorf("invalid installation ID: %v", err)
	}

	signedJWT, err := c.generateJWT()
	if err != nil {
		return "", err
	}

	// Request Installation Token
	urlStr := fmt.Sprintf("%s/app/installations/%d/access_tokens", c.BaseURL, iid)
	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+signedJWT)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get installation token (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}

// ListRepositories lists repositories. If token is provided, it uses it as a bearer token.
// If not, and an installation ID is known (fetched elsewhere), it would need that.
// For now, we update ListRepositories to be clearer about its usage.
func (c *GitHubClient) ListRepositories(ctx context.Context, token string, installationID string) ([]types.Repository, error) {
	if installationID != "" && token == "" {
		var err error
		token, err = c.GenerateInstallationToken(ctx, installationID)
		if err != nil {
			return nil, err
		}
	}
	var allRepos []types.Repository
	// If the token is an installation token, we might want to use a different endpoint
	// but /user/repos also works for OAuth. For installation tokens, we technically use /installation/repositories
	// however, if we are in a "User Context", /user/repos is correct.

	urlStr := c.BaseURL + "/user/repos?sort=updated&per_page=100&visibility=all&affiliation=owner,collaborator,organization_member"

	// NOTE: Installation tokens often use /installation/repositories instead of /user/repos
	// because /user/repos requires a USER context, whereas installation tokens have an INSTALLATION context.
	// If the token starts with 'ghs_', it's likely an installation token.
	if strings.HasPrefix(token, "ghs_") {
		urlStr = c.BaseURL + "/installation/repositories?per_page=100"
	}

	for urlStr != "" {
		req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("github api returned status: %s", resp.Status)
		}

        var githubRepos []struct {
            Name      string    `json:"name"`
            FullName  string    `json:"full_name"`
            HTMLURL   string    `json:"html_url"`
            Private   bool      `json:"private"`
            UpdatedAt time.Time `json:"updated_at"`
        }

        // Installation repositories are returned in a 'repositories' field
        if strings.HasPrefix(token, "ghs_") {
            var wrapper struct {
                Repositories []struct {
                    Name      string    `json:"name"`
                    FullName  string    `json:"full_name"`
                    HTMLURL   string    `json:"html_url"`
                    Private   bool      `json:"private"`
                    UpdatedAt time.Time `json:"updated_at"`
                } `json:"repositories"`
            }
            if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
                resp.Body.Close()
                return nil, err
            }
            for _, gr := range wrapper.Repositories {
                allRepos = append(allRepos, types.Repository{
                    Name:      gr.Name,
                    FullName:  gr.FullName,
                    URL:       gr.HTMLURL,
                    Private:   gr.Private,
                    UpdatedAt: gr.UpdatedAt,
                })
            }
        } else {
            if err := json.NewDecoder(resp.Body).Decode(&githubRepos); err != nil {
                resp.Body.Close()
                return nil, err
            }
            for _, gr := range githubRepos {
                allRepos = append(allRepos, types.Repository{
                    Name:      gr.Name,
                    FullName:  gr.FullName,
                    URL:       gr.HTMLURL,
                    Private:   gr.Private,
                    UpdatedAt: gr.UpdatedAt,
                })
            }
        }

		urlStr = extractNextPageURL(resp.Header.Get("Link"))
		resp.Body.Close()
	}

	return allRepos, nil
}

func (c *GitHubClient) ListRepositoriesPaginated(ctx context.Context, token string, installationID string, page int, limit int) ([]types.Repository, error) {
	if installationID != "" && token == "" {
		var err error
		token, err = c.GenerateInstallationToken(ctx, installationID)
		if err != nil {
			return nil, err
		}
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 100
	}

	urlStr := fmt.Sprintf("%s/user/repos?sort=updated&page=%d&per_page=%d&visibility=all&affiliation=owner,collaborator,organization_member", c.BaseURL, page, limit)
	if strings.HasPrefix(token, "ghs_") {
		urlStr = fmt.Sprintf("%s/installation/repositories?page=%d&per_page=%d", c.BaseURL, page, limit)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status: %s", resp.Status)
	}

	var githubRepos []struct {
		Name      string    `json:"name"`
		FullName  string    `json:"full_name"`
		HTMLURL   string    `json:"html_url"`
		Private   bool      `json:"private"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	if strings.HasPrefix(token, "ghs_") {
		var wrapper struct {
			Repositories []struct {
				Name      string    `json:"name"`
				FullName  string    `json:"full_name"`
				HTMLURL   string    `json:"html_url"`
				Private   bool      `json:"private"`
				UpdatedAt time.Time `json:"updated_at"`
			} `json:"repositories"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
			return nil, err
		}
		repos := make([]types.Repository, len(wrapper.Repositories))
		for i, gr := range wrapper.Repositories {
			repos[i] = types.Repository{
				Name:      gr.Name,
				FullName:  gr.FullName,
				URL:       gr.HTMLURL,
				Private:   gr.Private,
				UpdatedAt: gr.UpdatedAt,
			}
		}
		return repos, nil
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubRepos); err != nil {
		return nil, err
	}

	repos := make([]types.Repository, len(githubRepos))
	for i, gr := range githubRepos {
		repos[i] = types.Repository{
			Name:      gr.Name,
			FullName:  gr.FullName,
			URL:       gr.HTMLURL,
			Private:   gr.Private,
			UpdatedAt: gr.UpdatedAt,
		}
	}

	return repos, nil
}

func (c *GitHubClient) SearchRepositories(ctx context.Context, token string, query string, namespace string, page int, limit int, installationID string) ([]types.Repository, error) {
	if installationID != "" && token == "" {
		var err error
		token, err = c.GenerateInstallationToken(ctx, installationID)
		if err != nil {
			return nil, err
		}
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 100
	}

	q := query
	if namespace != "" && namespace != "all" {
		q = fmt.Sprintf("%s user:%s", query, namespace)
	}

	urlStr := fmt.Sprintf("%s/search/repositories?q=%s&page=%d&per_page=%d", c.BaseURL, url.QueryEscape(q), page, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status: %s", resp.Status)
	}

	var searchResult struct {
		Items []struct {
			Name      string    `json:"name"`
			FullName  string    `json:"full_name"`
			HTMLURL   string    `json:"html_url"`
			Private   bool      `json:"private"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, err
	}

	repos := make([]types.Repository, len(searchResult.Items))
	for i, gr := range searchResult.Items {
		repos[i] = types.Repository{
			Name:      gr.Name,
			FullName:  gr.FullName,
			URL:       gr.HTMLURL,
			Private:   gr.Private,
			UpdatedAt: gr.UpdatedAt,
		}
	}

	return repos, nil
}

func (c *GitHubClient) ListNamespaces(ctx context.Context, token string, installationID string) ([]types.Namespace, error) {
	if installationID != "" && token == "" {
		var err error
		token, err = c.GenerateInstallationToken(ctx, installationID)
		if err != nil {
			return nil, err
		}
	}
	namespaces := []types.Namespace{}

	// 1. If it's an installation token, we can't fetch "current user" or "installations for user"
	// because there is no USER context. We should just return the namespace for the installation itself.
	if strings.HasPrefix(token, "ghs_") {
		if installationID == "" {
			return namespaces, nil
		}
		var iid int64
		fmt.Sscanf(installationID, "%d", &iid)

		// Get account info from installation ID using App JWT
		signedJWT, err := c.generateJWT()
		if err == nil {
			req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/app/installations/%d", c.BaseURL, iid), nil)
			if err == nil {
				req.Header.Set("Authorization", "Bearer "+signedJWT)
				req.Header.Set("Accept", "application/vnd.github.v3+json")
				resp, err := c.HTTPClient.Do(req)
				if err == nil {
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						var result struct {
							Account struct {
								Login     string `json:"login"`
								AvatarURL string `json:"avatar_url"`
								Type      string `json:"type"`
							} `json:"account"`
						}
						if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
							namespaces = append(namespaces, types.Namespace{
								Name:      result.Account.Login,
								AvatarURL: result.Account.AvatarURL,
								Type:      result.Account.Type,
							})
						}
					}
				}
			}
		}

		// If we couldn't get it via JWT, try via installation token (list repos and get owner)
		if len(namespaces) == 0 {
			req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/installation/repositories?per_page=1", nil)
			if err == nil {
				req.Header.Set("Authorization", "Bearer "+token)
				req.Header.Set("Accept", "application/vnd.github.v3+json")
				resp, err := c.HTTPClient.Do(req)
				if err == nil {
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						var wrapper struct {
							Repositories []struct {
								Owner struct {
									Login     string `json:"login"`
									AvatarURL string `json:"avatar_url"`
									Type      string `json:"type"`
								} `json:"owner"`
							} `json:"repositories"`
						}
						if err := json.NewDecoder(resp.Body).Decode(&wrapper); err == nil && len(wrapper.Repositories) > 0 {
							namespaces = append(namespaces, types.Namespace{
								Name:      wrapper.Repositories[0].Owner.Login,
								AvatarURL: wrapper.Repositories[0].Owner.AvatarURL,
								Type:      wrapper.Repositories[0].Owner.Type,
							})
						}
					}
				}
			}
		}

		return namespaces, nil
	}

	// 1. Fetch User (Personal Account)
	userReq, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/user", nil)
	if err != nil {
		return nil, err
	}
	userReq.Header.Set("Authorization", "Bearer "+token)
	userReq.Header.Set("Accept", "application/vnd.github.v3+json")

	userResp, err := c.HTTPClient.Do(userReq)
	if err != nil {
		return nil, err
	}
	defer userResp.Body.Close()

	if userResp.StatusCode == http.StatusOK {
		var user struct {
			Login     string `json:"login"`
			AvatarURL string `json:"avatar_url"`
		}
		if err := json.NewDecoder(userResp.Body).Decode(&user); err == nil {
			namespaces = append(namespaces, types.Namespace{
				Name:      user.Login,
				AvatarURL: user.AvatarURL,
				Type:      "User",
			})
		}
	}

	// 2. Fetch Installations (The correct way for GitHub Apps)
	instUrlStr := c.BaseURL + "/user/installations?per_page=100"
	for instUrlStr != "" {
		instReq, err := http.NewRequestWithContext(ctx, "GET", instUrlStr, nil)
		if err != nil {
			break
		}
		instReq.Header.Set("Authorization", "Bearer "+token)
		instReq.Header.Set("Accept", "application/vnd.github.v3+json")

		instResp, err := c.HTTPClient.Do(instReq)
		if err != nil {
			break
		}

		if instResp.StatusCode == http.StatusOK {
			var result struct {
				Installations []struct {
					Account struct {
						Login     string `json:"login"`
						AvatarURL string `json:"avatar_url"`
						Type      string `json:"type"`
					} `json:"account"`
				} `json:"installations"`
			}

			if err := json.NewDecoder(instResp.Body).Decode(&result); err == nil {
				for _, inst := range result.Installations {
					exists := false
					for _, existing := range namespaces {
						if existing.Name == inst.Account.Login {
							exists = true
							break
						}
					}
					if !exists {
						namespaces = append(namespaces, types.Namespace{
							Name:      inst.Account.Login,
							AvatarURL: inst.Account.AvatarURL,
							Type:      inst.Account.Type,
						})
					}
				}
			}
			instUrlStr = extractNextPageURL(instResp.Header.Get("Link"))
		} else {
			// Log error or handle failure
			body, _ := io.ReadAll(instResp.Body)
			// We'll just ignore errors from installations if it's not a GitHub app and proceed.
			fmt.Printf("failed to fetch installations: %s\n", string(body))
			instUrlStr = ""
		}
		instResp.Body.Close()
	}

	// 3. Fetch Organizations
	orgsUrlStr := c.BaseURL + "/user/orgs?per_page=100"
	for orgsUrlStr != "" {
		orgsReq, err := http.NewRequestWithContext(ctx, "GET", orgsUrlStr, nil)
		if err != nil {
			break
		}
		orgsReq.Header.Set("Authorization", "Bearer "+token)
		orgsReq.Header.Set("Accept", "application/vnd.github.v3+json")

		orgsResp, err := c.HTTPClient.Do(orgsReq)
		if err != nil {
			break
		}

		if orgsResp.StatusCode == http.StatusOK {
			var orgs []struct {
				Login     string `json:"login"`
				AvatarURL string `json:"avatar_url"`
			}

			if err := json.NewDecoder(orgsResp.Body).Decode(&orgs); err == nil {
				for _, org := range orgs {
					exists := false
					for _, existing := range namespaces {
						if existing.Name == org.Login {
							exists = true
							break
						}
					}
					if !exists {
						namespaces = append(namespaces, types.Namespace{
							Name:      org.Login,
							AvatarURL: org.AvatarURL,
							Type:      "Organization",
						})
					}
				}
			}
			orgsUrlStr = extractNextPageURL(orgsResp.Header.Get("Link"))
		} else {
			orgsUrlStr = ""
		}
		orgsResp.Body.Close()
	}

	return namespaces, nil
}

func (c *GitHubClient) RefreshToken(ctx context.Context, refreshToken string) (*types.TokenRefreshResponse, error) {
	if c.ClientID == "" || c.ClientSecret == "" {
		return nil, fmt.Errorf("github client id or secret not configured")
	}

	data := url.Values{}
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to refresh github token, status: %s", resp.Status)
	}

	var result struct {
		AccessToken          string `json:"access_token"`
		RefreshToken         string `json:"refresh_token"`
		ExpiresIn            int    `json:"expires_in"`
		RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
		Error                string `json:"error"`
		ErrorDescription     string `json:"error_description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Error != "" {
		return nil, fmt.Errorf("github oauth error: %s - %s", result.Error, result.ErrorDescription)
	}

	res := &types.TokenRefreshResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}
	if result.ExpiresIn > 0 {
		res.ExpiresAt = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	}

	return res, nil
}

func (c *GitHubClient) ListContents(ctx context.Context, token string, repoFullName string, path string, installationID string) ([]types.ContentItem, error) {
	if installationID != "" && token == "" {
		var err error
		token, err = c.GenerateInstallationToken(ctx, installationID)
		if err != nil {
			return nil, err
		}
	}
	path = strings.TrimPrefix(path, "./")
	path = strings.TrimPrefix(path, "/")
	urlStr := fmt.Sprintf("%s/repos/%s/contents/%s", c.BaseURL, repoFullName, path)
	if path == "" || path == "." {
		urlStr = fmt.Sprintf("%s/repos/%s/contents", c.BaseURL, repoFullName)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return []types.ContentItem{}, nil
		}
		return nil, fmt.Errorf("github api returned status: %s", resp.Status)
	}

	var githubContents []struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Type string `json:"type"`
		Size int64  `json:"size"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubContents); err != nil {
		return nil, err
	}

	contents := make([]types.ContentItem, len(githubContents))
	for i, gc := range githubContents {
		contents[i] = types.ContentItem{
			Name: gc.Name,
			Path: gc.Path,
			Type: gc.Type,
			Size: gc.Size,
		}
	}

	return contents, nil
}
