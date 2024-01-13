package api

import (
	"context"
	"io"
	"net/http"
)

const defaultPingURI = "https://api.githubcopilot.com/_ping"

// Ping sends a ping request to the GitHub Copilot API.
func (c *Client) Ping(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultPingURI, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", statusError{
			StatusCode:   resp.StatusCode,
			Status:       resp.Status,
			ErrorMessage: string(body),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
