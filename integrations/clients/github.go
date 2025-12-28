package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shashtag-ventures/go-common/integrations/types"
)

type GitHubClient struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		BaseURL:    "https://api.github.com",
	}
}

func (c *GitHubClient) ListRepositories(ctx context.Context, token string) ([]types.Repository, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/user/repos?sort=updated&per_page=100", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)
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
	userReq.Header.Set("Authorization", "token "+token)
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
	instReq.Header.Set("Authorization", "token "+token)
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
		return namespaces, fmt.Errorf("failed to fetch installations: %s", string(body))
	}

	return namespaces, nil
}
