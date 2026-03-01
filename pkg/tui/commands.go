package tui

import (
	"fmt"

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
