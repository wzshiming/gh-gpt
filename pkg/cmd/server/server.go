package server

import (
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gh-gpt/pkg/api"
	"github.com/wzshiming/gh-gpt/pkg/auth"
	"github.com/wzshiming/gh-gpt/pkg/cache"
	"github.com/wzshiming/gh-gpt/pkg/server"
	"github.com/wzshiming/gh-gpt/pkg/utils"
)

type serverOptions struct {
	Address          string
	TokenCachePath   string
	GHTokenCachePath string
}

func NewCommand() *cobra.Command {
	opts := serverOptions{
		Address:          ":8000",
		TokenCachePath:   "~/.gh-gpt/token.json",
		GHTokenCachePath: "~/.gh-gpt/gh-token.json",
	}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("Starting server", "address", opts.Address)
			return run(cmd, opts)
		},
	}

	cmd.Flags().StringVar(&opts.Address, "address", opts.Address, "address")
	cmd.Flags().StringVar(&opts.TokenCachePath, "token-cache-path", opts.TokenCachePath, "token cache path")
	cmd.Flags().StringVar(&opts.GHTokenCachePath, "gh-token-cache-path", opts.GHTokenCachePath, "github token cache path")
	return cmd
}

func run(cmd *cobra.Command, opts serverOptions) error {
	auths := auth.Auths{}
	if opts.GHTokenCachePath != "" {
		tokenCachePath, err := utils.ExpandPath(opts.GHTokenCachePath)
		if err != nil {
			return err
		}
		auths = append(auths, auth.DeviceSession(tokenCachePath))
	}
	auths = append(auths, auth.Hosts(), auth.Envs())

	var err error
	opts.TokenCachePath, err = utils.ExpandPath(opts.TokenCachePath)
	if err != nil {
		return err
	}

	tokenCache := cache.NewFileCache(opts.TokenCachePath)
	cli := api.NewClient(
		api.WithTokenCache(tokenCache),
	)

	svc := server.NewServer(
		server.WithClient(cli),
		server.WithAuth(auths),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/chat/completions", svc.ChatCompletions)
	mux.HandleFunc("/_ping", svc.Ping)

	return http.ListenAndServe(opts.Address, server.CORS(mux))
}
