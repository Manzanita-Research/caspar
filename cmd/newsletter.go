package cmd

import (
	"fmt"

	"github.com/manzanita-research/ghostctl/pkg/ghost"
	"github.com/manzanita-research/ghostctl/pkg/output"
	"github.com/spf13/cobra"
)

var newsletterCmd = &cobra.Command{
	Use:   "newsletter",
	Short: "View newsletters",
}

var newsletterListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List newsletters",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		params := parseListParams(cmd)
		newsletters, pag, err := client.ListNewsletters(params)
		if err != nil {
			return err
		}
		return output.Print(jsonOut, map[string]any{"newsletters": newsletters, "meta": map[string]any{"pagination": pag}}, func() {
			if len(newsletters) == 0 {
				fmt.Println("No newsletters.")
				return
			}
			for _, n := range newsletters {
				fmt.Printf("  %-8s %-30s %s\n", n.Status, n.Name, n.Slug)
			}
		})
	},
}

var newsletterGetCmd = &cobra.Command{
	Use:     "get <id>",
	Aliases: []string{"show"},
	Short:   "Get a newsletter by ID",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		fields, _ := cmd.Flags().GetString("fields")
		newsletter, err := client.GetNewsletter(args[0], ghost.ListParams{Fields: fields})
		if err != nil {
			return err
		}
		return output.Print(jsonOut, newsletter, func() {
			output.Title(newsletter.Name)
			output.Field("ID", newsletter.ID)
			output.Field("Slug", newsletter.Slug)
			if newsletter.Description != "" {
				output.Field("Description", newsletter.Description)
			}
			output.Field("Status", newsletter.Status)
		})
	},
}

func init() {
	rootCmd.AddCommand(newsletterCmd)

	addListFlags(newsletterListCmd)
	newsletterCmd.AddCommand(newsletterListCmd)

	newsletterGetCmd.Flags().String("fields", "", "comma-separated fields to include")
	newsletterCmd.AddCommand(newsletterGetCmd)
}
