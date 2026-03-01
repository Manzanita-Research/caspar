package cmd

import (
	"fmt"

	"github.com/manzanita-research/caspar/pkg/ghost"
	"github.com/manzanita-research/caspar/pkg/output"
	"github.com/spf13/cobra"
)

var memberCmd = &cobra.Command{
	Use:   "member",
	Short: "Manage members",
}

var memberListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List members",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		params := parseListParams(cmd)
		members, pag, err := client.ListMembers(params)
		if err != nil {
			return err
		}
		return output.Print(jsonOut, map[string]any{"members": members, "meta": map[string]any{"pagination": pag}}, func() {
			if len(members) == 0 {
				fmt.Println("No members.")
				return
			}
			for _, m := range members {
				name := m.Name
				if name == "" {
					name = "(no name)"
				}
				fmt.Printf("  %-8s %-30s %s\n", m.Status, name, m.Email)
			}
		})
	},
}

var memberGetCmd = &cobra.Command{
	Use:     "get <id>",
	Aliases: []string{"show"},
	Short:   "Get a member by ID",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		fields, _ := cmd.Flags().GetString("fields")
		member, err := client.GetMember(args[0], ghost.ListParams{Fields: fields})
		if err != nil {
			return err
		}
		return output.Print(jsonOut, member, func() {
			name := member.Name
			if name == "" {
				name = member.Email
			}
			output.Title(name)
			output.Field("ID", member.ID)
			output.Field("Email", member.Email)
			output.Field("Status", member.Status)
			if len(member.Labels) > 0 {
				for _, l := range member.Labels {
					output.Field("Label", l.Name)
				}
			}
		})
	},
}

var memberCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a member",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		email, _ := cmd.Flags().GetString("email")
		name, _ := cmd.Flags().GetString("name")
		labels, _ := cmd.Flags().GetStringSlice("label")

		member, err := client.CreateMember(ghost.CreateMemberInput{
			Email:  email,
			Name:   name,
			Labels: labels,
		})
		if err != nil {
			return err
		}
		return output.Print(jsonOut, member, func() {
			output.Success(fmt.Sprintf("Created member: %s", member.Email))
			output.Field("ID", member.ID)
		})
	},
}

var memberUpdateCmd = &cobra.Command{
	Use:     "update <id>",
	Aliases: []string{"edit"},
	Short:   "Update a member",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		input := ghost.UpdateMemberInput{}
		if cmd.Flags().Changed("email") {
			v, _ := cmd.Flags().GetString("email")
			input.Email = &v
		}
		if cmd.Flags().Changed("name") {
			v, _ := cmd.Flags().GetString("name")
			input.Name = &v
		}
		if cmd.Flags().Changed("label") {
			v, _ := cmd.Flags().GetStringSlice("label")
			input.Labels = v
		}

		member, err := client.UpdateMember(args[0], input)
		if err != nil {
			return err
		}
		return output.Print(jsonOut, member, func() {
			output.Success(fmt.Sprintf("Updated member: %s", member.Email))
		})
	},
}

func init() {
	rootCmd.AddCommand(memberCmd)

	addListFlags(memberListCmd)
	memberCmd.AddCommand(memberListCmd)

	memberGetCmd.Flags().String("fields", "", "comma-separated fields to include")
	memberCmd.AddCommand(memberGetCmd)

	memberCreateCmd.Flags().String("email", "", "member email (required)")
	memberCreateCmd.Flags().String("name", "", "member name")
	memberCreateCmd.Flags().StringSlice("label", nil, "member label (repeatable)")
	_ = memberCreateCmd.MarkFlagRequired("email")
	memberCmd.AddCommand(memberCreateCmd)

	memberUpdateCmd.Flags().String("email", "", "new email")
	memberUpdateCmd.Flags().String("name", "", "new name")
	memberUpdateCmd.Flags().StringSlice("label", nil, "member label (repeatable, replaces existing)")
	memberCmd.AddCommand(memberUpdateCmd)
}
