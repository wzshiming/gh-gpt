package login

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gh-gpt/pkg/auth"
	"github.com/wzshiming/gh-gpt/pkg/utils"
)

type runOptions struct {
	GHTokenCachePath string
	GHClientID       string
}

func NewCommand() *cobra.Command {
	opts := runOptions{
		GHTokenCachePath: "~/.gh-gpt/gh-token.json",
		GHClientID:       "Iv1.b507a08c87ecfe98",
	}
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login the copilot",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, opts)
		},
	}

	cmd.Flags().StringVar(&opts.GHTokenCachePath, "gh-token-cache-path", opts.GHTokenCachePath, "github token cache path")
	cmd.Flags().StringVar(&opts.GHClientID, "gh-client-id", opts.GHClientID, "github client id")
	return cmd
}

func run(cmd *cobra.Command, opts runOptions) error {
	ctx := cmd.Context()

	tokenCachePath, err := utils.ExpandPath(opts.GHTokenCachePath)
	if err != nil {
		return err
	}

	token, err := auth.DeviceLogin(ctx, tokenCachePath, opts.GHClientID)
	if err != nil {
		return fmt.Errorf("failed to get oauth token: %w", err)
	}

	fmt.Println(token)
	return nil
}
