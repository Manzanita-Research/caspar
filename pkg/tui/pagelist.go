package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func pageListUpdate(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case pagesLoadedMsg:
		m.pages = msg.pages
		m.pagePag = msg.pagination
		m.loading = false
		m.pageCursor = 0
		m.pageExpanded = -1
		return m, nil

	case tea.KeyMsg:
		// When filtering, route keys to the text input.
		if m.pageFiltering {
			switch {
			case key.Matches(msg, keys.Escape):
				m.pageFiltering = false
				m.filterInput.SetValue("")
				m.filterInput.Blur()
				return m, nil
			case msg.String() == "enter":
				m.pageFiltering = false
				m.filterInput.Blur()
				m.loading = true
				return m, loadPages(m.client, 1, m.pageStatusFilter(), m.filterInput.Value())
			default:
				var cmd tea.Cmd
				m.filterInput, cmd = m.filterInput.Update(msg)
				return m, cmd
			}
		}

		switch {
		case key.Matches(msg, keys.Escape):
			if m.pageExpanded >= 0 {
				m.pageExpanded = -1
				return m, nil
			}
			m.currentView = viewDashboard
			m.filterInput.SetValue("")
			return m, nil

		case key.Matches(msg, keys.Down):
			if len(m.pages) > 0 && m.pageCursor < len(m.pages)-1 {
				m.pageCursor++
				m.pageExpanded = -1
			}
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.pageCursor > 0 {
				m.pageCursor--
				m.pageExpanded = -1
			}
			return m, nil

		case key.Matches(msg, keys.Enter):
			if len(m.pages) > 0 {
				if m.pageExpanded == m.pageCursor {
					m.pageExpanded = -1
				} else {
					m.pageExpanded = m.pageCursor
				}
			}
			return m, nil

		case key.Matches(msg, keys.Tab):
			m.pageStatusIdx = (m.pageStatusIdx + 1) % len(statusFilters)
			m.loading = true
			return m, loadPages(m.client, 1, m.pageStatusFilter(), m.filterInput.Value())

		case key.Matches(msg, keys.Filter):
			m.pageFiltering = true
			m.filterInput.Focus()
			return m, nil

		case msg.String() == "n":
			if m.pagePag != nil && m.pagePag.Next != nil {
				m.loading = true
				return m, loadPages(m.client, *m.pagePag.Next, m.pageStatusFilter(), m.filterInput.Value())
			}
			return m, nil

		case msg.String() == "N":
			if m.pagePag != nil && m.pagePag.Prev != nil {
				m.loading = true
				return m, loadPages(m.client, *m.pagePag.Prev, m.pageStatusFilter(), m.filterInput.Value())
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) pageStatusFilter() string {
	return statusFilters[m.pageStatusIdx]
}

func pageListView(m model) string {
	var b strings.Builder
	w := m.contentWidth()

	b.WriteString("\n")
	b.WriteString(siteHeader(m.site))
	b.WriteString(tabBar(viewPageList))
	b.WriteString("\n\n")

	// Status filter label.
	if sf := m.pageStatusFilter(); sf != "all" {
		b.WriteString(indent(labelStyle.Render("(" + sf + ")")))
		b.WriteString("\n")
	}

	// Active filter display.
	if m.filterInput.Value() != "" {
		b.WriteString(indent(labelStyle.Render("filter: ") + valueStyle.Render(m.filterInput.Value())))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	if m.loading {
		b.WriteString(indent(labelStyle.Render("loading...")))
		b.WriteString("\n")
		return b.String()
	}

	if len(m.pages) == 0 {
		b.WriteString(indent(labelStyle.Render("no pages found")))
		b.WriteString("\n")
	} else {
		for i, p := range m.pages {
			title := p.Title
			maxTitle := w - 8
			if maxTitle < 20 {
				maxTitle = 20
			}
			if len(title) > maxTitle {
				title = title[:maxTitle-3] + "..."
			}

			indicator := statusIndicator(p.Status)
			row := fmt.Sprintf("%s %s", indicator, title)

			if i == m.pageCursor {
				b.WriteString(indent(selectedRowStyle.Render(row)))
			} else {
				b.WriteString(indent(normalRowStyle.Render(row)))
			}
			b.WriteString("\n")

			// Expanded detail pane — reuse postDetailPane since Page = Post.
			if i == m.pageExpanded {
				b.WriteString(postDetailPane(p, w))
			}
		}
	}

	b.WriteString("\n")

	// Pagination.
	if m.pagePag != nil {
		pagInfo := fmt.Sprintf("page %d/%d (%d total)", m.pagePag.Page, m.pagePag.Pages, m.pagePag.Total)
		b.WriteString(indent(labelStyle.Render(pagInfo)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(indent(divider(w - 4)))
	b.WriteString("\n\n")

	// Filter input.
	if m.pageFiltering {
		b.WriteString(indent(filterPromptStyle.Render("/ ") + m.filterInput.View()))
		b.WriteString("\n\n")
	}

	// Help.
	b.WriteString(indent(m.help.View(pageListKeys{escape: keys.Escape})))
	b.WriteString("\n")

	return b.String()
}
