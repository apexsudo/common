package cmd

import (
	"github.com/apexsudo/common/cmd/scaffold/internal/storage"
	"github.com/spf13/cobra"
)

func newCreateMigrationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create_migration <name>",
		Short: "Scaffold a blank up/down migration pair",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.ScaffoldMigration(cmd.Context(), args[0])
		},
	}
}
