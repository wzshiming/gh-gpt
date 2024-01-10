package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gh-gpt/pkg/api"
	"github.com/wzshiming/gh-gpt/pkg/auth"
	"github.com/wzshiming/gh-gpt/pkg/cache"
	"github.com/wzshiming/gh-gpt/pkg/server"
)

type serverOptions struct {
	TokenCachePath string
	Address        string
}

func NewCommand() *cobra.Command {
	opts := serverOptions{
		TokenCachePath: "~/.gh-gpt/token.json",
		Address:        ":8000",
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("Starting server", "address", opts.Address)
			return run(cmd, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Address, "address", opts.Address, "address")
	cmd.Flags().StringVar(&opts.TokenCachePath, "token-cache-path", opts.TokenCachePath, "token cache path")

	return cmd
}

func run(cmd *cobra.Command, opts serverOptions) error {

	hosts := auth.Hosts()

	var err error
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

	svc := server.NewServer(
		server.WithClient(cli),
		server.WithAuth(hosts),
	)

	handler := server.CORS(http.HandlerFunc(svc.ChatCompletions))
	return http.ListenAndServe(opts.Address, handler)
}
