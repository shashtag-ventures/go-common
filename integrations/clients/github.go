package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/shashtag-ventures/go-common/integrations/types"
)

type GitHubClient struct {
	HTTPClient   *http.Client
	BaseURL      string
	ClientID     string
	ClientSecret string
}

func NewGitHubClient(clientID, clientSecret string) *GitHubClient {
	return &GitHubClient{
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		BaseURL:      "https://api.github.com",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
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

func (c *GitHubClient) ListRepositories(ctx context.Context, token string) ([]types.Repository, error) {
	var allRepos []types.Repository
	// Explicitly request all visibilities and affiliations to ensure private repos are fetched
	// See: https://docs.github.com/en/rest/repos/repos#list-repositories-for-the-authenticated-user
	urlStr := c.BaseURL + "/user/repos?sort=updated&per_page=100&visibility=all&affiliation=owner,collaborator,organization_member"

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

		urlStr = extractNextPageURL(resp.Header.Get("Link"))
		resp.Body.Close()
	}

	return allRepos, nil
}

func (c *GitHubClient) ListRepositoriesPaginated(ctx context.Context, token string, page int, limit int) ([]types.Repository, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 100
	}

	// Explicitly request all visibilities and affiliations for paginated requests
	urlStr := fmt.Sprintf("%s/user/repos?sort=updated&page=%d&per_page=%d&visibility=all&affiliation=owner,collaborator,organization_member", c.BaseURL, page, limit)

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

func (c *GitHubClient) SearchRepositories(ctx context.Context, token string, query string, namespace string, page int, limit int) ([]types.Repository, error) {
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

func (c *GitHubClient) ListNamespaces(ctx context.Context, token string) ([]types.Namespace, error) {
	namespaces := []types.Namespace{}

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

func (c *GitHubClient) ListContents(ctx context.Context, token string, repoFullName string, path string) ([]types.ContentItem, error) {
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
