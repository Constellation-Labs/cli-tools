package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "cli-tools",
		Short:   "DAG Command Line Utility Tools",
		Version: "v1.0.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	Verbose bool
)

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose output")
}
