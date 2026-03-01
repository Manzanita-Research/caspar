package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/manzanita-research/ghostctl/pkg/config"
	"github.com/manzanita-research/ghostctl/pkg/ghost"
	"github.com/manzanita-research/ghostctl/pkg/output"
	"github.com/spf13/cobra"
)

// resourceKind distinguishes posts from pages in shared command logic.
type resourceKind string

const (
	kindPost resourceKind = "post"
	kindPage resourceKind = "page"
)

func newClient() (*ghost.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	return ghost.NewClient(cfg.URL, cfg.AdminAPIKey), nil
}

// buildResourceCommands creates the full set of list|get|create|update|delete
// subcommands for a given resource kind (post or page).
func buildResourceCommands(kind resourceKind) *cobra.Command {
	name := string(kind)

	root := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Manage %ss", name),
	}

	// --- list ---
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   fmt.Sprintf("List %ss", name),
		RunE:    makeListFn(kind),
	}
	addListFlags(listCmd)
	root.AddCommand(listCmd)

	// --- get ---
	getCmd := &cobra.Command{
		Use:     "get <id-or-slug>",
		Aliases: []string{"show"},
		Short:   fmt.Sprintf("Get a %s by ID or slug", name),
		Args:    cobra.ExactArgs(1),
		RunE:    makeGetFn(kind),
	}
	getCmd.Flags().String("fields", "", "comma-separated fields to include")
	getCmd.Flags().String("include", "", "related resources to include (e.g. tags,authors)")
	root.AddCommand(getCmd)

	// --- create ---
	createCmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"new"},
		Short:   fmt.Sprintf("Create a %s", name),
		RunE:    makeCreateFn(kind),
	}
	addCreateFlags(createCmd)
	root.AddCommand(createCmd)

	// --- update ---
	updateCmd := &cobra.Command{
		Use:     "update <id-or-slug>",
		Aliases: []string{"edit"},
		Short:   fmt.Sprintf("Update a %s", name),
		Args:    cobra.ExactArgs(1),
		RunE:    makeUpdateFn(kind),
	}
	addUpdateFlags(updateCmd)
	root.AddCommand(updateCmd)

	// --- delete ---
	deleteCmd := &cobra.Command{
		Use:     "delete <id>",
		Aliases: []string{"rm"},
		Short:   fmt.Sprintf("Delete a %s", name),
		Args:    cobra.ExactArgs(1),
		RunE:    makeDeleteFn(kind),
	}
	root.AddCommand(deleteCmd)

	return root
}

func addListFlags(cmd *cobra.Command) {
	cmd.Flags().Int("limit", 15, "number of results")
	cmd.Flags().Int("page", 0, "page number")
	cmd.Flags().String("filter", "", "Ghost NQL filter expression")
	cmd.Flags().String("order", "", "sort order (e.g. published_at desc)")
	cmd.Flags().String("fields", "", "comma-separated fields to include")
	cmd.Flags().String("include", "", "related resources to include (e.g. tags,authors)")
}

func addCreateFlags(cmd *cobra.Command) {
	cmd.Flags().String("title", "", "title (required)")
	cmd.Flags().String("html", "", "HTML content")
	cmd.Flags().String("lexical", "", "Lexical JSON content")
	cmd.Flags().String("status", "", "status (draft, published, scheduled)")
	cmd.Flags().String("slug", "", "custom slug")
	cmd.Flags().StringSlice("tag", nil, "tag name (repeatable)")
	cmd.Flags().Bool("featured", false, "mark as featured")
	cmd.Flags().Bool("stdin", false, "read HTML content from stdin")
	_ = cmd.MarkFlagRequired("title")
}

func addUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().String("title", "", "new title")
	cmd.Flags().String("html", "", "new HTML content")
	cmd.Flags().String("lexical", "", "new Lexical JSON content")
	cmd.Flags().String("status", "", "new status (draft, published, scheduled)")
	cmd.Flags().String("slug", "", "new slug")
	cmd.Flags().StringSlice("tag", nil, "tag name (repeatable, replaces existing)")
	cmd.Flags().Bool("featured", false, "mark as featured")
	cmd.Flags().Bool("no-featured", false, "unmark as featured")
	cmd.Flags().Bool("stdin", false, "read HTML content from stdin")
}

func parseListParams(cmd *cobra.Command) ghost.ListParams {
	limit, _ := cmd.Flags().GetInt("limit")
	page, _ := cmd.Flags().GetInt("page")
	filter, _ := cmd.Flags().GetString("filter")
	order, _ := cmd.Flags().GetString("order")
	fields, _ := cmd.Flags().GetString("fields")
	include, _ := cmd.Flags().GetString("include")
	return ghost.ListParams{
		Limit:   limit,
		Page:    page,
		Filter:  filter,
		Order:   order,
		Fields:  fields,
		Include: include,
	}
}

func makeListFn(kind resourceKind) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		params := parseListParams(cmd)

		if kind == kindPost {
			posts, pag, err := client.ListPosts(params)
			if err != nil {
				return err
			}
			return output.Print(jsonOut, map[string]any{"posts": posts, "meta": map[string]any{"pagination": pag}}, func() {
				printPostTable(posts)
			})
		}

		pages, pag, err := client.ListPages(params)
		if err != nil {
			return err
		}
		return output.Print(jsonOut, map[string]any{"pages": pages, "meta": map[string]any{"pagination": pag}}, func() {
			printPostTable(pages)
		})
	}
}

func makeGetFn(kind resourceKind) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		fields, _ := cmd.Flags().GetString("fields")
		include, _ := cmd.Flags().GetString("include")
		params := ghost.ListParams{Fields: fields, Include: include}

		if kind == kindPost {
			post, err := client.GetPost(args[0], params)
			if err != nil {
				return err
			}
			return output.Print(jsonOut, post, func() {
				printPostDetail(post)
			})
		}

		page, err := client.GetPage(args[0], params)
		if err != nil {
			return err
		}
		return output.Print(jsonOut, page, func() {
			printPostDetail(page)
		})
	}
}

func makeCreateFn(kind resourceKind) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		html, _ := cmd.Flags().GetString("html")
		lexical, _ := cmd.Flags().GetString("lexical")
		status, _ := cmd.Flags().GetString("status")
		slug, _ := cmd.Flags().GetString("slug")
		tags, _ := cmd.Flags().GetStringSlice("tag")
		featured, _ := cmd.Flags().GetBool("featured")
		useStdin, _ := cmd.Flags().GetBool("stdin")

		if useStdin {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
			html = string(data)
		}

		input := ghost.CreatePostInput{
			Title:    title,
			HTML:     html,
			Lexical:  lexical,
			Status:   status,
			Slug:     slug,
			Tags:     tags,
			Featured: featured,
		}

		useHTML := html != ""

		if kind == kindPost {
			post, err := client.CreatePost(input, useHTML)
			if err != nil {
				return err
			}
			return output.Print(jsonOut, post, func() {
				output.Success(fmt.Sprintf("Created post: %s (%s)", post.Title, post.Status))
				output.Field("ID", post.ID)
				output.Field("Slug", post.Slug)
			})
		}

		page, err := client.CreatePage(input, useHTML)
		if err != nil {
			return err
		}
		return output.Print(jsonOut, page, func() {
			output.Success(fmt.Sprintf("Created page: %s (%s)", page.Title, page.Status))
			output.Field("ID", page.ID)
			output.Field("Slug", page.Slug)
		})
	}
}

func makeUpdateFn(kind resourceKind) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		idOrSlug := args[0]

		// fetch current version to get updated_at and resolve slug → ID
		var currentID, updatedAt string
		if kind == kindPost {
			post, err := client.GetPost(idOrSlug, ghost.ListParams{})
			if err != nil {
				return fmt.Errorf("fetching current %s: %w", kind, err)
			}
			currentID = post.ID
			if post.UpdatedAt != nil {
				updatedAt = post.UpdatedAt.Format("2006-01-02T15:04:05.000Z")
			}
		} else {
			page, err := client.GetPage(idOrSlug, ghost.ListParams{})
			if err != nil {
				return fmt.Errorf("fetching current %s: %w", kind, err)
			}
			currentID = page.ID
			if page.UpdatedAt != nil {
				updatedAt = page.UpdatedAt.Format("2006-01-02T15:04:05.000Z")
			}
		}

		input := ghost.UpdatePostInput{UpdatedAt: updatedAt}

		if cmd.Flags().Changed("title") {
			v, _ := cmd.Flags().GetString("title")
			input.Title = &v
		}
		if cmd.Flags().Changed("status") {
			v, _ := cmd.Flags().GetString("status")
			input.Status = &v
		}
		if cmd.Flags().Changed("slug") {
			v, _ := cmd.Flags().GetString("slug")
			input.Slug = &v
		}
		if cmd.Flags().Changed("tag") {
			v, _ := cmd.Flags().GetStringSlice("tag")
			input.Tags = v
		}
		if cmd.Flags().Changed("featured") {
			v := true
			input.Featured = &v
		}
		if cmd.Flags().Changed("no-featured") {
			v := false
			input.Featured = &v
		}

		useStdin, _ := cmd.Flags().GetBool("stdin")
		if useStdin {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
			s := string(data)
			input.HTML = &s
		} else if cmd.Flags().Changed("html") {
			v, _ := cmd.Flags().GetString("html")
			input.HTML = &v
		}

		if cmd.Flags().Changed("lexical") {
			v, _ := cmd.Flags().GetString("lexical")
			input.Lexical = &v
		}

		useHTML := input.HTML != nil

		if kind == kindPost {
			post, err := client.UpdatePost(currentID, input, useHTML)
			if err != nil {
				return err
			}
			return output.Print(jsonOut, post, func() {
				output.Success(fmt.Sprintf("Updated post: %s (%s)", post.Title, post.Status))
			})
		}

		page, err := client.UpdatePage(currentID, input, useHTML)
		if err != nil {
			return err
		}
		return output.Print(jsonOut, page, func() {
			output.Success(fmt.Sprintf("Updated page: %s (%s)", page.Title, page.Status))
		})
	}
}

func makeDeleteFn(kind resourceKind) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		id := args[0]
		if !ghost.IsID(id) {
			return fmt.Errorf("delete requires a %s ID, not a slug — use `ghostctl %s get <slug>` to find the ID", kind, kind)
		}

		if kind == kindPost {
			if err := client.DeletePost(id); err != nil {
				return err
			}
		} else {
			if err := client.DeletePage(id); err != nil {
				return err
			}
		}

		if jsonOut {
			return output.JSON(map[string]string{"status": "deleted", "id": id})
		}
		output.Success(fmt.Sprintf("Deleted %s %s", kind, id))
		return nil
	}
}

// printPostTable prints a human-friendly table of posts/pages.
func printPostTable(posts []ghost.Post) {
	if len(posts) == 0 {
		fmt.Println("No results.")
		return
	}
	for _, p := range posts {
		status := p.Status
		title := p.Title
		if len(title) > 60 {
			title = title[:57] + "..."
		}
		fmt.Printf("  %-8s %-62s %s\n", status, title, p.Slug)
	}
}

// printPostDetail prints a single post/page in human-friendly format.
func printPostDetail(p *ghost.Post) {
	output.Title(p.Title)
	output.Field("ID", p.ID)
	output.Field("Slug", p.Slug)
	output.Field("Status", p.Status)
	if p.URL != "" {
		output.Field("URL", p.URL)
	}
	if len(p.Tags) > 0 {
		names := make([]string, len(p.Tags))
		for i, t := range p.Tags {
			names[i] = t.Name
		}
		output.Field("Tags", strings.Join(names, ", "))
	}
	if p.PublishedAt != nil {
		output.Field("Published", p.PublishedAt.Format("2006-01-02 15:04"))
	}
	if p.Excerpt != "" {
		fmt.Println()
		fmt.Println(p.Excerpt)
	}
}
