package logout

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wzshiming/gh-gpt/pkg/auth"
	"github.com/wzshiming/gh-gpt/pkg/utils"
)

type runOptions struct {
	GHTokenCachePath string
}

func NewCommand() *cobra.Command {
	opts := runOptions{
		GHTokenCachePath: "~/.gh-gpt/gh-token.json",
	}
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout the copilot",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, opts)
		},
	}

	cmd.Flags().StringVar(&opts.GHTokenCachePath, "gh-token-cache-path", opts.GHTokenCachePath, "github token cache path")
	return cmd
}

func run(cmd *cobra.Command, opts runOptions) error {
	ctx := cmd.Context()

	tokenCachePath, err := utils.ExpandPath(opts.GHTokenCachePath)
	if err != nil {
		return err
	}

	err = auth.DeviceLogout(ctx, tokenCachePath)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}
