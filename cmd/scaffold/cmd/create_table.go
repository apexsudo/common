package cmd

import (
	"github.com/apexsudo/common/cmd/scaffold/internal/storage"
	"github.com/spf13/cobra"
)

func newCreateTableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create_table <name>",
		Short: "Scaffold up/down migration files for a new table",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.ScaffoldTable(cmd.Context(), args[0])
		},
	}
}
