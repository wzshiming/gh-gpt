package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wzshiming/gh-gpt/pkg/cmd/ping"
	"github.com/wzshiming/gh-gpt/pkg/cmd/run"
	"github.com/wzshiming/gh-gpt/pkg/cmd/server"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.NoArgs,
		Use:  "gh-gpt [command]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		run.NewCommand(),
		server.NewCommand(),
		ping.NewCommand(),
	)
	return cmd
}
