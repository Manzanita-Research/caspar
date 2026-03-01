package cmd

import (
	"fmt"

	"github.com/manzanita-research/caspar/pkg/ghost"
	"github.com/manzanita-research/caspar/pkg/output"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
}

var tagListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		params := parseListParams(cmd)
		tags, pag, err := client.ListTags(params)
		if err != nil {
			return err
		}
		return output.Print(jsonOut, map[string]any{"tags": tags, "meta": map[string]any{"pagination": pag}}, func() {
			if len(tags) == 0 {
				fmt.Println("No tags.")
				return
			}
			for _, t := range tags {
				desc := t.Description
				if len(desc) > 50 {
					desc = desc[:47] + "..."
				}
				fmt.Printf("  %-30s %s\n", t.Name, desc)
			}
		})
	},
}

var tagGetCmd = &cobra.Command{
	Use:     "get <id-or-slug>",
	Aliases: []string{"show"},
	Short:   "Get a tag by ID or slug",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		fields, _ := cmd.Flags().GetString("fields")
		include, _ := cmd.Flags().GetString("include")
		tag, err := client.GetTag(args[0], ghost.ListParams{Fields: fields, Include: include})
		if err != nil {
			return err
		}
		return output.Print(jsonOut, tag, func() {
			output.Title(tag.Name)
			output.Field("ID", tag.ID)
			output.Field("Slug", tag.Slug)
			if tag.Description != "" {
				output.Field("Description", tag.Description)
			}
			output.Field("Visibility", tag.Visibility)
		})
	},
}

var tagCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a tag",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		slug, _ := cmd.Flags().GetString("slug")
		description, _ := cmd.Flags().GetString("description")
		visibility, _ := cmd.Flags().GetString("visibility")

		tag, err := client.CreateTag(ghost.CreateTagInput{
			Name:        name,
			Slug:        slug,
			Description: description,
			Visibility:  visibility,
		})
		if err != nil {
			return err
		}
		return output.Print(jsonOut, tag, func() {
			output.Success(fmt.Sprintf("Created tag: %s", tag.Name))
			output.Field("ID", tag.ID)
			output.Field("Slug", tag.Slug)
		})
	},
}

var tagUpdateCmd = &cobra.Command{
	Use:     "update <id-or-slug>",
	Aliases: []string{"edit"},
	Short:   "Update a tag",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		// resolve slug → ID if needed
		tag, err := client.GetTag(args[0], ghost.ListParams{})
		if err != nil {
			return fmt.Errorf("fetching tag: %w", err)
		}

		input := ghost.UpdateTagInput{}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			input.Name = &v
		}
		if cmd.Flags().Changed("slug") {
			v, _ := cmd.Flags().GetString("slug")
			input.Slug = &v
		}
		if cmd.Flags().Changed("description") {
			v, _ := cmd.Flags().GetString("description")
			input.Description = &v
		}
		if cmd.Flags().Changed("visibility") {
			v, _ := cmd.Flags().GetString("visibility")
			input.Visibility = &v
		}

		updated, err := client.UpdateTag(tag.ID, input)
		if err != nil {
			return err
		}
		return output.Print(jsonOut, updated, func() {
			output.Success(fmt.Sprintf("Updated tag: %s", updated.Name))
		})
	},
}

var tagDeleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Aliases: []string{"rm"},
	Short:   "Delete a tag",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		id := args[0]
		if !ghost.IsID(id) {
			return fmt.Errorf("delete requires a tag ID, not a slug — use `caspar tag get <slug>` to find the ID")
		}
		if err := client.DeleteTag(id); err != nil {
			return err
		}
		if jsonOut {
			return output.JSON(map[string]string{"status": "deleted", "id": id})
		}
		output.Success(fmt.Sprintf("Deleted tag %s", id))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)

	addListFlags(tagListCmd)
	tagCmd.AddCommand(tagListCmd)

	tagGetCmd.Flags().String("fields", "", "comma-separated fields to include")
	tagGetCmd.Flags().String("include", "", "related resources to include")
	tagCmd.AddCommand(tagGetCmd)

	tagCreateCmd.Flags().String("name", "", "tag name (required)")
	tagCreateCmd.Flags().String("slug", "", "custom slug")
	tagCreateCmd.Flags().String("description", "", "tag description")
	tagCreateCmd.Flags().String("visibility", "", "visibility (public, internal)")
	_ = tagCreateCmd.MarkFlagRequired("name")
	tagCmd.AddCommand(tagCreateCmd)

	tagUpdateCmd.Flags().String("name", "", "new name")
	tagUpdateCmd.Flags().String("slug", "", "new slug")
	tagUpdateCmd.Flags().String("description", "", "new description")
	tagUpdateCmd.Flags().String("visibility", "", "new visibility")
	tagCmd.AddCommand(tagUpdateCmd)

	tagCmd.AddCommand(tagDeleteCmd)
}
