package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wzshiming/gh-gpt/pkg/cmd/run"
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
	)
	return cmd
}
