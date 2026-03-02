package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/manzanita-research/caspar/pkg/ghost"
)

func tagListUpdate(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tagsLoadedMsg:
		m.tags = msg.tags
		m.tagPag = msg.pagination
		m.loading = false
		m.tagCursor = 0
		m.tagExpanded = -1
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape):
			if m.tagExpanded >= 0 {
				m.tagExpanded = -1
				return m, nil
			}
			m.currentView = viewDashboard
			return m, nil

		case key.Matches(msg, keys.Down):
			if len(m.tags) > 0 && m.tagCursor < len(m.tags)-1 {
				m.tagCursor++
				m.tagExpanded = -1
			}
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.tagCursor > 0 {
				m.tagCursor--
				m.tagExpanded = -1
			}
			return m, nil

		case key.Matches(msg, keys.Enter):
			if len(m.tags) > 0 {
				if m.tagExpanded == m.tagCursor {
					m.tagExpanded = -1
				} else {
					m.tagExpanded = m.tagCursor
				}
			}
			return m, nil

		case msg.String() == "n":
			if m.tagPag != nil && m.tagPag.Next != nil {
				m.loading = true
				return m, loadTags(m.client, *m.tagPag.Next)
			}
			return m, nil

		case msg.String() == "N":
			if m.tagPag != nil && m.tagPag.Prev != nil {
				m.loading = true
				return m, loadTags(m.client, *m.tagPag.Prev)
			}
			return m, nil
		}
	}
	return m, nil
}

func tagListView(m model) string {
	var b strings.Builder
	w := m.contentWidth()

	b.WriteString("\n")
	b.WriteString(indent(titleStyle.Render("Tags")))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(indent(labelStyle.Render("loading...")))
		b.WriteString("\n")
		return b.String()
	}

	if len(m.tags) == 0 {
		b.WriteString(indent(labelStyle.Render("no tags found")))
		b.WriteString("\n")
	} else {
		for i, t := range m.tags {
			name := t.Name
			slug := t.Slug
			maxName := w - len(slug) - 8
			if maxName < 20 {
				maxName = 20
			}
			if len(name) > maxName {
				name = name[:maxName-3] + "..."
			}

			if i == m.tagCursor {
				b.WriteString(indent(selectedRowStyle.Render(name) + "  " + labelStyle.Render(slug)))
			} else {
				b.WriteString(indent(normalRowStyle.Render(name) + "  " + labelStyle.Render(slug)))
			}
			b.WriteString("\n")

			if i == m.tagExpanded {
				b.WriteString(tagDetailPane(t))
			}
		}
	}

	b.WriteString("\n")

	if m.tagPag != nil {
		pagInfo := fmt.Sprintf("page %d/%d (%d total)", m.tagPag.Page, m.tagPag.Pages, m.tagPag.Total)
		b.WriteString(indent(labelStyle.Render(pagInfo)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(indent(divider(w - 4)))
	b.WriteString("\n\n")

	b.WriteString(indent(m.help.View(tagListKeys{})))
	b.WriteString("\n")

	return b.String()
}

func tagDetailPane(t ghost.Tag) string {
	var b strings.Builder
	pad := "      "

	b.WriteString("\n")
	b.WriteString(pad + detailLabelStyle.Render("slug      ") + detailValueStyle.Render(t.Slug) + "\n")
	if t.Visibility != "" {
		b.WriteString(pad + detailLabelStyle.Render("visibility") + " " + detailValueStyle.Render(t.Visibility) + "\n")
	}
	if t.Description != "" {
		b.WriteString(pad + detailLabelStyle.Render("desc      ") + detailValueStyle.Render(t.Description) + "\n")
	}
	b.WriteString("\n")

	return b.String()
}
