package run

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gh-gpt/pkg/api"
	"github.com/wzshiming/gh-gpt/pkg/auth"
	"github.com/wzshiming/gh-gpt/pkg/cache"
	"golang.org/x/term"
)

type runOptions struct {
	Model          string
	System         string
	Content        string
	Stream         bool
	TokenCachePath string
}

func NewCommand() *cobra.Command {
	opts := runOptions{
		Model:          "gpt-4",
		System:         "You are a helpful assistant.",
		Stream:         true,
		TokenCachePath: "~/.gh-gpt/token.json",
	}
	cmd := &cobra.Command{
		Use:   "run [content...]",
		Short: "Run the gpt",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Content = strings.Join(args, " ")
			} else if !term.IsTerminal(int(os.Stdin.Fd())) {
				in, err := io.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				opts.Content = string(in)
			}

			if opts.Content == "" {
				return fmt.Errorf("no content")
			}

			return run(cmd, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Model, "model", opts.Model, "model")
	cmd.Flags().StringVar(&opts.System, "system", opts.System, "system")
	cmd.Flags().BoolVar(&opts.Stream, "stream", opts.Stream, "stream")
	cmd.Flags().StringVar(&opts.TokenCachePath, "token-cache-path", opts.TokenCachePath, "token cache path")

	return cmd
}

func run(cmd *cobra.Command, opts runOptions) error {
	ctx := cmd.Context()
	hosts := auth.Hosts()

	oauth, err := hosts.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get oauth token: %w", err)
	}

	// expand the '~' for opts.TokenCachePath
	if strings.HasPrefix(opts.TokenCachePath, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home dir: %w", err)
		}
		opts.TokenCachePath = filepath.Join(home, opts.TokenCachePath[1:])
	}
	opts.TokenCachePath, err = filepath.Abs(opts.TokenCachePath)
	if err != nil {
		return fmt.Errorf("failed to get token cache path: %w", err)
	}

	tokenCache := cache.NewFileCache(opts.TokenCachePath)
	cli := api.NewClient(
		api.WithTokenCache(tokenCache),
	)

	token, err := cli.TokenWishCache(ctx, oauth)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	req := api.ChatRequest{
		Model: opts.Model,
		Messages: []api.Message{
			{
				Role:    "system",
				Content: opts.System,
			},
			{
				Role:    "user",
				Content: opts.Content,
			},
		},
		Stream: opts.Stream,
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
