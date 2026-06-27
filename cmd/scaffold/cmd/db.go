package cmd

import "github.com/spf13/cobra"

func newDBCmd() *cobra.Command {
	dbCmd := &cobra.Command{
		Use:   "dbCmd",
		Short: "Database scaffolding commands",
	}
	dbCmd.AddCommand(newCreateTableCmd())
	dbCmd.AddCommand(newCreateMigrationCmd())

	return dbCmd
}
