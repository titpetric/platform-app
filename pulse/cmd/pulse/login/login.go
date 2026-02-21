// Package login implements the pulse login command.
package login

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/titpetric/cli"
	"golang.org/x/term"

	"github.com/titpetric/platform-app/pulse/client"
)

// Name is the command title.
const Name = "Authenticate with pulse server"

// Options holds login command configuration.
type Options struct {
	Server   string
	Email    string
	Password string
}

// Bind registers login flags with the flag set.
func (o *Options) Bind(flag *cli.FlagSet) {
	defaultServer := os.Getenv("PULSE_SERVER")
	if defaultServer == "" {
		defaultServer = "http://localhost:8080"
	}
	flag.StringVar(&o.Server, "server", defaultServer, "Pulse server URL")
	flag.StringVar(&o.Email, "email", "", "Email address (optional, will prompt if not provided)")
	flag.StringVar(&o.Password, "password", "", "Password (optional, will prompt if not provided)")
}

// NewCommand creates a new login command.
func NewCommand() *cli.Command {
	var opts Options

	return &cli.Command{
		Name:  "login",
		Title: Name,
		Bind: func(flag *cli.FlagSet) {
			opts.Bind(flag)
		},
		Run: func(ctx context.Context, args []string) error {
			return Run(opts)
		},
	}
}

// Run executes the login flow with the given options.
func Run(opts Options) error {
	reader := bufio.NewReader(os.Stdin)

	email := opts.Email
	if email == "" {
		fmt.Print("Email: ")
		var err error
		email, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read email: %w", err)
		}
		email = strings.TrimSpace(email)
	}

	password := opts.Password
	if password == "" {
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("read password: %w", err)
		}
		fmt.Println()
		password = string(passwordBytes)
	}

	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", opts.Server+"/api/user/token/create", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	expiresAt := time.Unix(result.ExpiresAt, 0)
	c := client.New(opts.Server)
	if err := c.SaveToken(result.Token, expiresAt); err != nil {
		return fmt.Errorf("save token: %w", err)
	}

	fmt.Println("Login successful! Token saved.")
	return nil
}
