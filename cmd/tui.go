package cmd

import (
	"github.com/manzanita-research/caspar/pkg/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui [resource]",
	Short: "Interactive terminal UI",
	Long:  "Launch the interactive TUI. Optionally pass a resource name (e.g. posts) to jump straight to that view.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		startView := ""
		if len(args) > 0 {
			startView = args[0]
		}

		return tui.Run(client, startView)
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
