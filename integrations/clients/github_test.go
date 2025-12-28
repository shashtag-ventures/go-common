package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubClient_ListRepositories(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "token test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/vnd.github.v3+json", r.Header.Get("Accept"))
		assert.Contains(t, r.URL.Path, "/user/repos")

		// Return mock response
		repos := []map[string]any{
			{
				"name":      "repo1",
				"full_name": "user/repo1",
				"html_url":  "https://github.com/user/repo1",
				"private":   false,
			},
		}
		json.NewEncoder(w).Encode(repos)
	}))
	defer server.Close()

	client := NewGitHubClient()
	client.BaseURL = server.URL // Point to mock server

	repos, err := client.ListRepositories(context.Background(), "test-token")
	assert.NoError(t, err)
	assert.Len(t, repos, 1)
	assert.Equal(t, "repo1", repos[0].Name)
	assert.Equal(t, "https://github.com/user/repo1", repos[0].URL)
}

func TestGitHubClient_ListNamespaces(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "token test-token", r.Header.Get("Authorization"))

		if r.URL.Path == "/user" {
			user := map[string]any{
				"login":      "testuser",
				"avatar_url": "https://avatar.com/u",
			}
			json.NewEncoder(w).Encode(user)
			return
		}

		if r.URL.Path == "/user/installations" {
			inst := map[string]any{
				"installations": []map[string]any{
					{
						"account": map[string]any{
							"login":      "org1",
							"avatar_url": "https://avatar.com/o",
							"type":       "Organization",
						},
					},
				},
			}
			json.NewEncoder(w).Encode(inst)
			return
		}
	}))
	defer server.Close()

	client := NewGitHubClient()
	client.BaseURL = server.URL

	namespaces, err := client.ListNamespaces(context.Background(), "test-token")
	assert.NoError(t, err)
	assert.Len(t, namespaces, 2)
	assert.Equal(t, "testuser", namespaces[0].Name)
	assert.Equal(t, "User", namespaces[0].Type)
	assert.Equal(t, "org1", namespaces[1].Name)
	assert.Equal(t, "Organization", namespaces[1].Type)
}

func TestGitHubClient_Errors(t *testing.T) {
	t.Run("API Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		client := NewGitHubClient()
		client.BaseURL = server.URL

		_, err := client.ListRepositories(context.Background(), "bad-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "returned status: 401")
	})

	t.Run("Malformed JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{invalid json`))
		}))
		defer server.Close()

		client := NewGitHubClient()
		client.BaseURL = server.URL

		_, err := client.ListRepositories(context.Background(), "token")
		assert.Error(t, err)
	})
}
