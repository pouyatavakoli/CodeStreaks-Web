package codeforces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pouyatavakoli/CodeStreaks-web/internal/domain"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetUserSubmissions fetches submissions for a user with a specified count
func (c *Client) GetUserSubmissions(handle string, count int) ([]domain.CodeforcesSubmission, error) {
	url := fmt.Sprintf("%s/user.status?handle=%s&from=1&count=%d", c.baseURL, handle, count)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch submissions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp domain.CodeforcesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API returned non-OK status: %s", apiResp.Status)
	}

	return apiResp.Result, nil
}

// GetUserInfo fetches user information from Codeforces
func (c *Client) GetUserInfo(handle string) (*domain.CodeforcesUserInfo, error) {
	url := fmt.Sprintf("%s/user.info?handles=%s", c.baseURL, handle)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp domain.CodeforcesUserAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResp.Status != "OK" || len(apiResp.Result) == 0 {
		return nil, fmt.Errorf("user not found or API error")
	}

	return &apiResp.Result[0], nil
}

// ValidateHandle checks if a Codeforces handle exists
func (c *Client) ValidateHandle(handle string) (bool, error) {
	_, err := c.GetUserInfo(handle)
	if err != nil {
		return false, err
	}
	return true, nil
}
