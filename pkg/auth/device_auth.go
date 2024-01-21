package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/wzshiming/getch"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type deviceSessionAuth struct {
	cachePath string
}

func DeviceSession(cachePath string) Auth {
	return &deviceSessionAuth{
		cachePath: cachePath,
	}
}

func (a *deviceSessionAuth) GetToken(ctx context.Context) (string, error) {
	token, err := loadToken[oauth2.Token](a.cachePath)
	if err != nil {
		return "", err
	}
	if token.Valid() {
		return token.AccessToken, nil
	}
	return "", fmt.Errorf("not login")
}

func DeviceLogout(ctx context.Context, cachePath string) error {
	return os.Remove(cachePath)
}

func DeviceLogin(ctx context.Context, cachePath string, clientID string) (string, error) {
	token, err := loadToken[oauth2.Token](cachePath)
	if err == nil {
		if token.Valid() {
			return token.AccessToken, nil
		}
	}

	config := &oauth2.Config{
		ClientID: clientID,
		Endpoint: github.Endpoint,
		Scopes:   []string{"read:user"},
	}

	resp, err := config.DeviceAuth(ctx)
	if err != nil {
		return "", err
	}

	if resp.VerificationURIComplete != "" {
		fmt.Printf("Please visit %s to authenticate.\n", resp.VerificationURIComplete)
	} else {
		fmt.Printf("Please take this code %q to authenticate at %s.\n", resp.UserCode, resp.VerificationURI)
	}
	fmt.Println("Press 'y' to continue, or any other key to abort.")
	ch, _, err := getch.Getch()
	if err != nil {
		return "", err
	}
	if ch != 'y' && ch != 'Y' {
		return "", fmt.Errorf("aborted")
	}
	fmt.Println("Waiting for authentication...")

	token, err = config.DeviceAccessToken(ctx, resp)
	if err != nil {
		return "", err
	}
	if token.AccessToken != "" {
		err = saveToken[oauth2.Token](cachePath, token)
		if err != nil {
			return "", err
		}
		return token.AccessToken, nil
	}

	return "", fmt.Errorf("failed to get token from %q", cachePath)
}

func saveToken[T any](cachePath string, data *T) error {
	err := os.MkdirAll(path.Dir(cachePath), 0750)
	if err != nil {
		return err
	}

	d, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(cachePath, d, 0600)
	if err != nil {
		return err
	}
	return nil
}

func loadToken[T any](cachePath string) (*T, error) {
	d, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var data *T
	err = json.Unmarshal(d, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
