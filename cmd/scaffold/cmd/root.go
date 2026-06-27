package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "scaffold",
		Short: "Code generation tool for ApexSudo backend services",
	}
	root.AddCommand(newDBCmd())

	return root
}

func Execute() error {
	if err := newRootCmd().Execute(); err != nil {
		return fmt.Errorf("scaffold: %w", err)
	}

	return nil
}
