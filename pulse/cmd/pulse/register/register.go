// Package register implements the pulse user registration command.
package register

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
	"github.com/titpetric/platform-app/pulse/cmd/pulse/login"
)

// Name is the command title.
const Name = "Register a new account"

// Options holds register command configuration.
type Options struct {
	Server   string
	Email    string
	Password string
	Username string
}

// Bind registers registration flags with the flag set.
func (o *Options) Bind(flag *cli.FlagSet) {
	defaultServer := os.Getenv("PULSE_SERVER")
	if defaultServer == "" {
		defaultServer = "http://localhost:8080"
	}
	flag.StringVar(&o.Server, "server", defaultServer, "Pulse server URL")
	flag.StringVar(&o.Email, "email", "", "Email address (optional, will prompt if not provided)")
	flag.StringVar(&o.Password, "password", "", "Password (optional, will prompt if not provided)")
	flag.StringVar(&o.Username, "username", "", "Username (optional, will prompt if not provided)")
}

// NewCommand creates a new register command.
func NewCommand() *cli.Command {
	var opts Options

	return &cli.Command{
		Name:  "register",
		Title: Name,
		Bind: func(flag *cli.FlagSet) {
			opts.Bind(flag)
		},
		Run: func(ctx context.Context, args []string) error {
			return Run(opts)
		},
	}
}

// Run executes the registration flow with the given options.
func Run(opts Options) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Full name: ")
	fullName, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read full name: %w", err)
	}
	fullName = strings.TrimSpace(fullName)

	email := opts.Email
	if email == "" {
		fmt.Print("Email: ")
		email, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read email: %w", err)
		}
		email = strings.TrimSpace(email)
	}

	username := opts.Username
	if username == "" {
		fmt.Print("Username: ")
		username, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read username: %w", err)
		}
		username = strings.TrimSpace(username)
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
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		FullName: fullName,
		Email:    email,
		Username: username,
		Password: password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", opts.Server+"/api/user/register", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("register request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		// Username already taken, try to login instead
		fmt.Println("User already exists, attempting login...")
		return login.Run(login.Options{
			Server:   opts.Server,
			Email:    email,
			Password: password,
		})
	}

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		UserID    string `json:"user_id"`
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

	fmt.Printf("Registration successful! User ID: %s\n", result.UserID)
	return nil
}
