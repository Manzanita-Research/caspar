package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/manzanita-research/caspar/pkg/ghost"
)

func memberListUpdate(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case membersLoadedMsg:
		m.members = msg.members
		m.memberPag = msg.pagination
		m.loading = false
		m.memberCursor = 0
		m.memberExpanded = -1
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape):
			if m.memberExpanded >= 0 {
				m.memberExpanded = -1
				return m, nil
			}
			m.currentView = viewDashboard
			return m, nil

		case key.Matches(msg, keys.Down):
			if len(m.members) > 0 && m.memberCursor < len(m.members)-1 {
				m.memberCursor++
				m.memberExpanded = -1
			}
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.memberCursor > 0 {
				m.memberCursor--
				m.memberExpanded = -1
			}
			return m, nil

		case key.Matches(msg, keys.Enter):
			if len(m.members) > 0 {
				if m.memberExpanded == m.memberCursor {
					m.memberExpanded = -1
				} else {
					m.memberExpanded = m.memberCursor
				}
			}
			return m, nil

		case msg.String() == "n":
			if m.memberPag != nil && m.memberPag.Next != nil {
				m.loading = true
				return m, loadMembers(m.client, *m.memberPag.Next)
			}
			return m, nil

		case msg.String() == "N":
			if m.memberPag != nil && m.memberPag.Prev != nil {
				m.loading = true
				return m, loadMembers(m.client, *m.memberPag.Prev)
			}
			return m, nil
		}
	}
	return m, nil
}

func memberListView(m model) string {
	var b strings.Builder
	w := m.contentWidth()

	b.WriteString("\n")
	b.WriteString(siteHeader(m.site))
	b.WriteString(tabBar(viewMemberList))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(indent(labelStyle.Render("loading...")))
		b.WriteString("\n")
		return b.String()
	}

	if len(m.members) == 0 {
		b.WriteString(indent(labelStyle.Render("no members found")))
		b.WriteString("\n")
	} else {
		for i, mem := range m.members {
			if i == m.memberCursor {
				if mem.Name != "" {
					b.WriteString(indent(selectedRowStyle.Render(mem.Name) + "  " + labelStyle.Render(mem.Email)))
				} else {
					b.WriteString(indent(selectedRowStyle.Render(mem.Email)))
				}
			} else {
				if mem.Name != "" {
					b.WriteString(indent(normalRowStyle.Render(mem.Name) + "  " + labelStyle.Render(mem.Email)))
				} else {
					b.WriteString(indent(normalRowStyle.Render(mem.Email)))
				}
			}
			b.WriteString("\n")

			if i == m.memberExpanded {
				b.WriteString(memberDetailPane(mem))
			}
		}
	}

	b.WriteString("\n")

	if m.memberPag != nil {
		pagInfo := fmt.Sprintf("page %d/%d (%d total)", m.memberPag.Page, m.memberPag.Pages, m.memberPag.Total)
		b.WriteString(indent(labelStyle.Render(pagInfo)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(indent(divider(w - 4)))
	b.WriteString("\n\n")

	b.WriteString(indent(m.help.View(memberListKeys{})))
	b.WriteString("\n")

	return b.String()
}

func memberDetailPane(mem ghost.Member) string {
	var b strings.Builder
	pad := "      "

	b.WriteString("\n")
	b.WriteString(pad + detailLabelStyle.Render("email  ") + detailValueStyle.Render(mem.Email) + "\n")
	if mem.Status != "" {
		b.WriteString(pad + detailLabelStyle.Render("status ") + detailValueStyle.Render(mem.Status) + "\n")
	}
	if len(mem.Labels) > 0 {
		names := make([]string, len(mem.Labels))
		for i, l := range mem.Labels {
			names[i] = l.Name
		}
		b.WriteString(pad + detailLabelStyle.Render("labels ") + tagStyle.Render(strings.Join(names, ", ")) + "\n")
	}
	if mem.CreatedAt != nil {
		b.WriteString(pad + detailLabelStyle.Render("joined ") + detailValueStyle.Render(mem.CreatedAt.Format("2006-01-02")) + "\n")
	}
	b.WriteString("\n")

	return b.String()
}
