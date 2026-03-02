package tui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/manzanita-research/caspar/pkg/ghost"
)

func loadDashboard(client *ghost.Client) tea.Cmd {
	return func() tea.Msg {
		site, err := client.GetSite()
		if err != nil {
			return errMsg{fmt.Errorf("loading site: %w", err)}
		}

		_, postPag, err := client.ListPosts(ghost.ListParams{Limit: 1})
		if err != nil {
			return errMsg{fmt.Errorf("counting posts: %w", err)}
		}

		_, pagePag, err := client.ListPages(ghost.ListParams{Limit: 1})
		if err != nil {
			return errMsg{fmt.Errorf("counting pages: %w", err)}
		}

		_, memberPag, err := client.ListMembers(ghost.ListParams{Limit: 1})
		if err != nil {
			return errMsg{fmt.Errorf("counting members: %w", err)}
		}

		_, tagPag, err := client.ListTags(ghost.ListParams{Limit: 1})
		if err != nil {
			return errMsg{fmt.Errorf("counting tags: %w", err)}
		}

		msg := siteLoadedMsg{site: site}
		if postPag != nil {
			msg.postCount = postPag.Total
		}
		if pagePag != nil {
			msg.pageCount = pagePag.Total
		}
		if memberPag != nil {
			msg.memberCount = memberPag.Total
		}
		if tagPag != nil {
			msg.tagCount = tagPag.Total
		}
		return msg
	}
}

func loadPostDetail(client *ghost.Client, id string) tea.Cmd {
	return func() tea.Msg {
		post, err := client.GetPost(id, ghost.ListParams{
			Include: "tags,authors",
		})
		if err != nil {
			return errMsg{fmt.Errorf("loading post: %w", err)}
		}
		return postDetailMsg{post: *post}
	}
}

func loadPages(client *ghost.Client, page int, statusFilter, nqlFilter string) tea.Cmd {
	return func() tea.Msg {
		filter := ""
		if statusFilter != "all" {
			filter = "status:" + statusFilter
		}
		if nqlFilter != "" {
			if filter != "" {
				filter += "+" + nqlFilter
			} else {
				filter = nqlFilter
			}
		}

		params := ghost.ListParams{
			Limit:   15,
			Page:    page,
			Filter:  filter,
			Include: "tags",
		}

		pages, pag, err := client.ListPages(params)
		if err != nil {
			return errMsg{fmt.Errorf("loading pages: %w", err)}
		}
		return pagesLoadedMsg{pages: pages, pagination: pag}
	}
}

func searchPosts(client *ghost.Client, query string) tea.Cmd {
	return func() tea.Msg {
		filter := "title:~'" + query + "'"

		params := ghost.ListParams{
			Limit:   15,
			Filter:  filter,
			Include: "tags",
		}

		posts, pag, err := client.ListPosts(params)
		if err != nil {
			return errMsg{fmt.Errorf("searching posts: %w", err)}
		}
		return postsLoadedMsg{posts: posts, pagination: pag}
	}
}

func togglePostStatus(client *ghost.Client, postID string) tea.Cmd {
	return func() tea.Msg {
		current, err := client.GetPost(postID, ghost.ListParams{})
		if err != nil {
			return postToggleErrMsg{fmt.Errorf("fetching post: %w", err)}
		}

		if current.Status == "scheduled" {
			return postToggleErrMsg{fmt.Errorf("cannot toggle scheduled posts")}
		}

		var newStatus string
		if current.Status == "published" {
			newStatus = "draft"
		} else {
			newStatus = "published"
		}

		updatedAt := ""
		if current.UpdatedAt != nil {
			updatedAt = current.UpdatedAt.Format("2006-01-02T15:04:05.000Z")
		}

		updated, err := client.UpdatePost(postID, ghost.UpdatePostInput{
			Status:    &newStatus,
			UpdatedAt: updatedAt,
		}, false)
		if err != nil {
			return postToggleErrMsg{fmt.Errorf("updating post: %w", err)}
		}

		return postToggledMsg{post: *updated}
	}
}

func openInBrowser(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", url)
		default:
			cmd = exec.Command("xdg-open", url)
		}
		_ = cmd.Start()
		return nil
	}
}

func openGhostEditor(baseURL, postID string) tea.Cmd {
	editorURL := strings.TrimRight(baseURL, "/") + "/ghost/#/editor/post/" + postID
	return openInBrowser(editorURL)
}

func loadPosts(client *ghost.Client, page int, statusFilter, nqlFilter string) tea.Cmd {
	return func() tea.Msg {
		filter := ""
		if statusFilter != "all" {
			filter = "status:" + statusFilter
		}
		if nqlFilter != "" {
			if filter != "" {
				filter += "+" + nqlFilter
			} else {
				filter = nqlFilter
			}
		}

		params := ghost.ListParams{
			Limit:   15,
			Page:    page,
			Filter:  filter,
			Include: "tags",
		}

		posts, pag, err := client.ListPosts(params)
		if err != nil {
			return errMsg{fmt.Errorf("loading posts: %w", err)}
		}
		return postsLoadedMsg{posts: posts, pagination: pag}
	}
}

func loadTags(client *ghost.Client, page int) tea.Cmd {
	return func() tea.Msg {
		params := ghost.ListParams{
			Limit: 15,
			Page:  page,
		}
		tags, pag, err := client.ListTags(params)
		if err != nil {
			return errMsg{fmt.Errorf("loading tags: %w", err)}
		}
		return tagsLoadedMsg{tags: tags, pagination: pag}
	}
}

func loadMembers(client *ghost.Client, page int) tea.Cmd {
	return func() tea.Msg {
		params := ghost.ListParams{
			Limit: 15,
			Page:  page,
		}
		members, pag, err := client.ListMembers(params)
		if err != nil {
			return errMsg{fmt.Errorf("loading members: %w", err)}
		}
		return membersLoadedMsg{members: members, pagination: pag}
	}
}
