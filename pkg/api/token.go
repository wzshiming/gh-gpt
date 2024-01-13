package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultTokenURI = "https://api.github.com/copilot_internal/v2/token"

func (c *Client) TokenWishCache(ctx context.Context, oauthToken string) (string, error) {
	token, err := c.tokenCache.Get(oauthToken)
	if err == nil && token != "" {
		return token, nil
	}

	resp, err := c.Token(ctx, oauthToken)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	err = c.tokenCache.PutWithExpires(oauthToken, resp.Token, time.Unix(resp.ExpiresAt, 0))
	if err != nil {
		return "", fmt.Errorf("failed to put token: %w", err)
	}

	return resp.Token, nil
}

// Token retrieves the token from the GitHub Copilot API.
func (c *Client) Token(ctx context.Context, oauthToken string) (*TokenResponse, error) {
	// https://github.blog/2021-04-05-behind-githubs-new-authentication-token-formats/#identifiable-prefixes
	if !strings.HasPrefix(oauthToken, "gh") {
		return nil, errors.New("invalid token format")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultTokenURI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "GitHubCopilotChat/0.8.0")
	req.Header.Set("Authorization", "token "+oauthToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, statusError{
			StatusCode:   resp.StatusCode,
			Status:       resp.Status,
			ErrorMessage: string(body),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResponse *TokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, err
	}

	return tokenResponse, nil
}

type TokenResponse struct {
	ExpiresAt int64  `json:"expires_at"`
	Token     string `json:"token"`

	AnnotationsEnabled bool   `json:"annotations_enabled"`
	ChatEnabled        bool   `json:"chat_enabled"`
	ChatJetbrains      bool   `json:"chat_jetbrains_enabled"`
	CodeQuoteEnabled   bool   `json:"code_quote_enabled"`
	CopilotIDEAgent    bool   `json:"copilot_ide_agent_chat_gpt4_small_prompt"`
	CopilotIgnore      bool   `json:"copilotignore_enabled"`
	IntellijEditor     bool   `json:"intellij_editor_fetcher"`
	Prompt8k           bool   `json:"prompt_8k"`
	PublicSuggestions  string `json:"public_suggestions"`
	RefreshIn          int64  `json:"refresh_in"`
	Sku                string `json:"sku"`
	SnippyLoadTest     bool   `json:"snippy_load_test_enabled"`
	Telemetry          string `json:"telemetry"`
	TrackingID         string `json:"tracking_id"`
	VSCPanel           bool   `json:"vsc_panel_v2"`
}
