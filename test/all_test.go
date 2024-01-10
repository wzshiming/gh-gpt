package test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/wzshiming/gh-gpt/pkg/api"
	"github.com/wzshiming/gh-gpt/pkg/auth"
	"github.com/wzshiming/gh-gpt/pkg/cache"
)

func TestA(t *testing.T) {
	err := example()
	if err != nil {
		t.Fatal(err)
	}
}

func example() error {
	ctx := context.Background()
	hosts := auth.Hosts()

	oauth, err := hosts.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get oauth token: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}

	tokenCachePath := filepath.Join(home, ".gh-gpt/token.json")
	tokenCache := cache.NewFileCache(tokenCachePath)
	cli := api.NewClient(
		api.WithTokenCache(tokenCache),
	)

	token, err := cli.TokenWishCache(ctx, oauth)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	req := api.ChatRequest{
		Model: "gpt-4",
		Messages: []api.Message{
			{Role: "system", Content: "You are hosts helpful assistant."},
			{Role: "user", Content: "Who are you?"},
		},
	}

	fn := func(resp api.ChatResponse) error {
		for _, choice := range resp.Choices {
			if choice.Delta.Content != "" {
				fmt.Print(choice.Delta.Content)
			}
			if choice.Message.Content != "" {
				fmt.Println(choice.Message.Content)
			}
		}
		return nil
	}

	err = cli.ChatCompletions(ctx, token, &req, fn)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return fmt.Errorf("failed to chat: %w", err)
	}
	return nil
}
