package run

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gh-gpt/pkg/api"
	"github.com/wzshiming/gh-gpt/pkg/auth"
	pkgrun "github.com/wzshiming/gh-gpt/pkg/run"
	"github.com/wzshiming/gh-gpt/pkg/utils"
	"golang.org/x/term"
)

type runOptions struct {
	Model            string
	System           string
	Content          string
	Stream           bool
	TokenCachePath   string
	GHTokenCachePath string
}

func NewCommand() *cobra.Command {
	opts := runOptions{
		Model:            "gpt-4",
		System:           "You are a helpful assistant.",
		Stream:           true,
		TokenCachePath:   "~/.gh-gpt/token.json",
		GHTokenCachePath: "~/.gh-gpt/gh-token.json",
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
	cmd.Flags().StringVar(&opts.GHTokenCachePath, "gh-token-cache-path", opts.GHTokenCachePath, "github token cache path")
	return cmd
}

func run(cmd *cobra.Command, opts runOptions) error {
	ctx := cmd.Context()

	auths := auth.Auths{}
	if opts.GHTokenCachePath != "" {
		tokenCachePath, err := utils.ExpandPath(opts.GHTokenCachePath)
		if err != nil {
			return err
		}
		auths = append(auths, auth.DeviceSession(tokenCachePath))
	}
	auths = append(auths, auth.Hosts(), auth.Envs())

	runOpts := []pkgrun.Option{
		pkgrun.WithModel(opts.Model),
		pkgrun.WithMessages([]api.Message{
			{Role: "system", Content: opts.System},
		}),
		pkgrun.WithTokenCachePath(opts.TokenCachePath),
		pkgrun.WithAuth(auths),
	}
	if opts.Stream {
		return pkgrun.RunStream(ctx, opts.Content, os.Stdout, runOpts...)
	}

	resp, err := pkgrun.Run(ctx, opts.Content, runOpts...)
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}
