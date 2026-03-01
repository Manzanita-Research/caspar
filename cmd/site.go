package cmd

import (
	"github.com/manzanita-research/ghostctl/pkg/config"
	"github.com/manzanita-research/ghostctl/pkg/ghost"
	"github.com/manzanita-research/ghostctl/pkg/output"
	"github.com/spf13/cobra"
)

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Show site information",
	RunE:  runSite,
}

func init() {
	rootCmd.AddCommand(siteCmd)
}

func runSite(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := ghost.NewClient(cfg.URL, cfg.AdminAPIKey)
	site, err := client.GetSite()
	if err != nil {
		return err
	}

	return output.Print(jsonOut, site, func() {
		output.Title(site.Title)
		if site.Description != "" {
			output.Field("Description", site.Description)
		}
		output.Field("URL", site.URL)
		output.Field("Version", site.Version)
	})
}
