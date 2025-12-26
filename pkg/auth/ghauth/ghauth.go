package ghauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

const (
	keyringService = "datatug-app"
	keyringUser    = "github-oauth-token"
)

// SaveToken stores the GitHub OAuth token in the system keyring.
func SaveToken(token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}
	if err := keyring.Set(keyringService, keyringUser, string(data)); err != nil {
		return fmt.Errorf("failed to save token to keyring: %w", err)
	}
	return nil
}

// GetToken retrieves the GitHub OAuth token from the system keyring.
func GetToken() (*oauth2.Token, error) {
	data, err := keyring.Get(keyringService, keyringUser)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from keyring: %w", err)
	}
	var token oauth2.Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}
	return &token, nil
}

// DeleteToken removes the GitHub OAuth token from the system keyring.
func DeleteToken() error {
	if err := keyring.Delete(keyringService, keyringUser); err != nil {
		return fmt.Errorf("failed to delete token from keyring: %w", err)
	}
	return nil
}

// DeviceCodeResponse represents the response from GitHub's device code request.
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// RequestDeviceCode requests a device code from GitHub.
func RequestDeviceCode(ctx context.Context, clientID string) (*DeviceCodeResponse, error) {
	data := map[string]string{
		"client_id": clientID,
		"scope":     "repo",
	}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/device/code", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &res, nil
}

// PollForToken polls GitHub for an access token.
func PollForToken(ctx context.Context, clientID, clientSecret, deviceCode string, interval int, onAttempt func(attempt int)) (*oauth2.Token, error) {
	if interval == 0 {
		interval = 5
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	attempt := 0
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			attempt++
			if onAttempt != nil {
				onAttempt(attempt)
			}
			token, err := requestToken(ctx, clientID, clientSecret, deviceCode)
			if err != nil {
				if errors.Is(err, errAuthorizationPending) {
					continue
				}
				if errors.Is(err, errSlowDown) {
					interval += 5
					ticker.Reset(time.Duration(interval) * time.Second)
					continue
				}
				return nil, err
			}
			return token, nil
		}
	}
}

var (
	errAuthorizationPending = errors.New("authorization_pending")
	errSlowDown             = errors.New("slow_down")
)

func requestToken(ctx context.Context, clientID, clientSecret, deviceCode string) (*oauth2.Token, error) {
	data := map[string]string{
		"client_id":   clientID,
		"device_code": deviceCode,
		"grant_type":  "urn:ietf:params:oauth:grant-type:device_code",
	}
	if clientSecret != "" {
		data["client_secret"] = clientSecret
	}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var errRes struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	if err = json.Unmarshal(respBody, &errRes); err == nil && errRes.Error != "" {
		switch errRes.Error {
		case "authorization_pending":
			return nil, errAuthorizationPending
		case "slow_down":
			return nil, errSlowDown
		default:
			return nil, fmt.Errorf("github error: %s (%s)", errRes.Error, errRes.ErrorDescription)
		}
	}

	var token oauth2.Token
	if err = json.Unmarshal(respBody, &token); err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	return &token, nil
}
