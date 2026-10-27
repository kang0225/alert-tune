package cli

import "github.com/spf13/cobra"

// NewRootCmd builds the alert-tune command tree.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "alert-tune",
		Short:        "Prometheus 알림 규칙을 백테스트하고 튜닝합니다",
		Version:      Version,
		SilenceUsage: true,
	}
	cmd.SetVersionTemplate("alert-tune {{.Version}}\n")
	cmd.AddCommand(newVersionCmd())
	return cmd
}
