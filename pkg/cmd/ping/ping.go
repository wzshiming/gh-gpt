package ping

import (
	"github.com/spf13/cobra"
	"github.com/wzshiming/gh-gpt/pkg/api"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping the server to check if it is alive",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd)
		},
	}
	return cmd
}

func run(cmd *cobra.Command) error {
	ctx := cmd.Context()

	cli := api.NewClient()
	resp, err := cli.Ping(ctx)
	if err != nil {
		return err
	}
	cmd.Println(resp)
	return nil
}
