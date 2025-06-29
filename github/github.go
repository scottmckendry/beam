// Package github provides a GitHub API client for Beam.
package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Repo represents minimal metadata for a GitHub repository.
type Repo struct {
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
}

// Client is a GitHub API client.
type Client struct {
	Token string
}

// NewClient creates a new GitHub API client with the given token.
func NewClient(token string) *Client {
	return &Client{Token: token}
}

// GetRepo fetches metadata for the given owner/repo.
func (c *Client) GetRepo(owner, repo string) (*Repo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}
	var r Repo
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}
