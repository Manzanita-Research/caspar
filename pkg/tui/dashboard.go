package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func dashboardUpdate(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case siteLoadedMsg:
		m.site = msg.site
		m.postCount = msg.postCount
		m.pageCount = msg.pageCount
		m.memberCount = msg.memberCount
		m.tagCount = msg.tagCount
		m.ready = true
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Posts):
			m.currentView = viewPostList
			m.loading = true
			return m, loadPosts(m.client, 1, m.statusFilter, "")
		case key.Matches(msg, keys.Tags):
			m.currentView = viewTagList
			m.loading = true
			return m, loadTags(m.client, 1)
		case key.Matches(msg, keys.Members):
			m.currentView = viewMemberList
			m.loading = true
			return m, loadMembers(m.client, 1)
		}
	}
	return m, nil
}

func dashboardView(m model) string {
	var b strings.Builder
	w := m.contentWidth()

	if m.loading && m.site == nil {
		b.WriteString("\n")
		b.WriteString(indent(labelStyle.Render("loading...")))
		b.WriteString("\n")
		return b.String()
	}

	if m.site == nil {
		return ""
	}

	b.WriteString("\n")
	b.WriteString(indent(titleStyle.Render(m.site.Title)))
	b.WriteString("\n")
	b.WriteString(indent(subtitleStyle.Render(m.site.URL)))
	b.WriteString("\n")
	b.WriteString(indent(subtitleStyle.Render("Ghost "+m.site.Version)))
	b.WriteString("\n\n")

	b.WriteString(indent(divider(w - 4)))
	b.WriteString("\n\n")

	b.WriteString(indent(countStyle.Render(fmt.Sprintf("%d", m.postCount)) + labelStyle.Render(" posts")))
	b.WriteString("\n")
	b.WriteString(indent(countStyle.Render(fmt.Sprintf("%d", m.pageCount)) + labelStyle.Render(" pages")))
	b.WriteString("\n")
	b.WriteString(indent(countStyle.Render(fmt.Sprintf("%d", m.memberCount)) + labelStyle.Render(" members")))
	b.WriteString("\n")
	b.WriteString(indent(countStyle.Render(fmt.Sprintf("%d", m.tagCount)) + labelStyle.Render(" tags")))
	b.WriteString("\n\n")

	b.WriteString(indent(divider(w - 4)))
	b.WriteString("\n\n")

	b.WriteString(indent(m.help.View(dashboardKeys{})))
	b.WriteString("\n")

	return b.String()
}
