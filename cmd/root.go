package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	jsonOut bool
)

var rootCmd = &cobra.Command{
	Use:   "ghostctl",
	Short: "A Claude-first CLI for Ghost CMS",
	Long:  "A clean Go CLI for Ghost CMS. Agents are the primary user — humans are welcome too.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output as JSON")
	rootCmd.Version = version
}
