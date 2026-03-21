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

func (c *GitHubClient) ListRepositories(ctx context.Context, token string) ([]types.Repository, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/user/repos?sort=updated&per_page=100", nil)
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
	instReq, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/user/installations?per_page=100", nil)
	if err != nil {
		return nil, err
	}
	instReq.Header.Set("Authorization", "Bearer "+token)
	instReq.Header.Set("Accept", "application/vnd.github.v3+json")

	instResp, err := c.HTTPClient.Do(instReq)
	if err != nil {
		return nil, err
	}
	defer instResp.Body.Close()

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
	} else {
		// Log error or handle failure
		body, _ := io.ReadAll(instResp.Body)
		// We'll just ignore errors from installations if it's not a GitHub app and proceed.
		// Alternatively, you might want to return here. But it's safer to just log and continue.
		fmt.Printf("failed to fetch installations: %s\n", string(body))
	}

	// 3. Fetch Organizations
	orgsReq, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/user/orgs?per_page=100", nil)
	if err == nil {
		orgsReq.Header.Set("Authorization", "Bearer "+token)
		orgsReq.Header.Set("Accept", "application/vnd.github.v3+json")

		orgsResp, err := c.HTTPClient.Do(orgsReq)
		if err == nil {
			defer orgsResp.Body.Close()

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
			}
		}
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
