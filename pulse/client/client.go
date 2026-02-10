// Package client provides an HTTP client for the pulse server API.
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// TokenData holds authentication token data with expiration info.
type TokenData struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	SavedAt   time.Time `json:"saved_at"`
}

// Client is an HTTP client for the pulse server.
type Client struct {
	ServerURL  string
	httpClient *http.Client
	token      *TokenData
}

// New creates a new pulse client for the given server URL.
func New(serverURL string) *Client {
	return &Client{
		ServerURL:  serverURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	if home == "" {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".config", "pulse", "token.json"), nil
}

// LoadToken reads the saved authentication token from disk.
func (c *Client) LoadToken() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("not logged in, run 'pulse login' first")
		}
		return fmt.Errorf("read token file: %w", err)
	}

	var token TokenData
	if err := json.Unmarshal(data, &token); err != nil {
		return fmt.Errorf("parse token file: %w", err)
	}

	c.token = &token
	return nil
}

// SaveToken persists the authentication token to disk.
func (c *Client) SaveToken(token string, expiresAt time.Time) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data := TokenData{
		Token:     token,
		ExpiresAt: expiresAt,
		SavedAt:   time.Now(),
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal token: %w", err)
	}

	if err := os.WriteFile(path, jsonData, 0o600); err != nil {
		return fmt.Errorf("write token file: %w", err)
	}

	c.token = &data
	return nil
}

// ShouldRefresh reports whether the token needs refreshing.
func (c *Client) ShouldRefresh() bool {
	if c.token == nil {
		return true
	}
	if time.Now().After(c.token.ExpiresAt) {
		return true
	}
	if time.Since(c.token.SavedAt) > 24*time.Hour {
		return true
	}
	return false
}

// RefreshToken exchanges the current token for a new one.
func (c *Client) RefreshToken() error {
	if c.token == nil {
		return errors.New("no token loaded")
	}

	req, err := http.NewRequest("POST", c.ServerURL+"/api/user/token/refresh", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("refresh request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("refresh failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	expiresAt := time.Unix(result.ExpiresAt, 0)
	return c.SaveToken(result.Token, expiresAt)
}

// EnsureToken loads and refreshes the token if needed.
func (c *Client) EnsureToken() error {
	if err := c.LoadToken(); err != nil {
		return err
	}
	if c.ShouldRefresh() {
		if err := c.RefreshToken(); err != nil {
			return fmt.Errorf("refresh token: %w", err)
		}
	}
	return nil
}

// SendPulse submits a keystroke count to the server.
func (c *Client) SendPulse(count int64, hostname string) error {
	if c.token == nil {
		return errors.New("not authenticated")
	}

	payload := struct {
		Count    int64  `json:"count"`
		Hostname string `json:"hostname"`
	}{Count: count, Hostname: hostname}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.ServerURL+"/api/pulse/ingest", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("send failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// Token returns the current authentication token string.
func (c *Client) Token() string {
	if c.token == nil {
		return ""
	}
	return c.token.Token
}
