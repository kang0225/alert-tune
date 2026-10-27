package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "alert-tune 버전을 출력합니다",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			_, err := fmt.Fprintf(command.OutOrStdout(), "alert-tune %s\n", Version)
			return err
		},
	}
	return cmd
}
