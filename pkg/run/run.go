package run

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/wzshiming/gh-gpt/pkg/api"
	"github.com/wzshiming/gh-gpt/pkg/auth"
	"github.com/wzshiming/gh-gpt/pkg/cache"
	"github.com/wzshiming/gh-gpt/pkg/utils"
)

type option struct {
	Model          string
	TokenCachePath string
	Messages       []api.Message
	Stream         bool
	Auth           auth.Auth
}

type Option func(*option)

func WithModel(model string) Option {
	return func(opt *option) {
		opt.Model = model
	}
}

func WithTokenCachePath(path string) Option {
	return func(opt *option) {
		opt.TokenCachePath = path
	}
}

func WithAuth(a auth.Auths) Option {
	return func(opt *option) {
		opt.Auth = a
	}
}

func WithMessages(messages []api.Message) Option {
	return func(opt *option) {
		opt.Messages = messages
	}
}

func Run(ctx context.Context, content string, opts ...Option) (string, error) {
	var buf strings.Builder
	err := run(ctx, content, false, &buf, opts...)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RunStream(ctx context.Context, content string, out io.Writer, opts ...Option) error {
	return run(ctx, content, true, out, opts...)
}

func run(ctx context.Context, content string, stream bool, out io.Writer, opts ...Option) error {
	opt := option{
		Model:          "gpt-4",
		TokenCachePath: "~/.gh-gpt/token.json",
		Auth:           auth.Auths{auth.Hosts(), auth.Envs()},
		Stream:         stream,
	}
	for _, o := range opts {
		o(&opt)
	}

	oauth, err := opt.Auth.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get oauth token: %w", err)
	}

	opt.TokenCachePath, err = utils.ExpandPath(opt.TokenCachePath)
	if err != nil {
		return err
	}

	tokenCache := cache.NewFileCache(opt.TokenCachePath)
	cli := api.NewClient(
		api.WithTokenCache(tokenCache),
	)

	token, err := cli.TokenWishCache(ctx, oauth)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	req := api.ChatRequest{
		Model:    opt.Model,
		Messages: opt.Messages,
		Stream:   opt.Stream,
	}

	req.Messages = append(req.Messages, api.Message{
		Role:    "user",
		Content: content,
	})

	fn := func(resp api.ChatResponse) error {
		for _, choice := range resp.Choices {
			if choice.Delta.Content != "" {
				fmt.Fprint(out, choice.Delta.Content)
			}
			if choice.Message.Content != "" {
				fmt.Fprintln(out, choice.Message.Content)
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
