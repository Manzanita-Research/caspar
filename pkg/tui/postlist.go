package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/manzanita-research/caspar/pkg/ghost"
)

var statusFilters = []string{"all", "published", "draft", "scheduled"}

func postListUpdate(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case postsLoadedMsg:
		m.posts = msg.posts
		m.postPag = msg.pagination
		m.loading = false
		m.cursor = 0
		m.selected = -1
		return m, nil

	case postToggledMsg:
		m.loading = false
		newStatus := msg.post.Status
		if newStatus == "published" {
			m.statusMsg = "post published"
		} else {
			m.statusMsg = "post set to draft"
		}
		m.statusErr = ""
		// Reload posts to reflect the change.
		return m, loadPosts(m.client, 1, m.statusFilter, m.filterInput.Value())

	case postToggleErrMsg:
		m.loading = false
		m.statusErr = msg.err.Error()
		m.statusMsg = ""
		return m, nil

	case tea.KeyMsg:
		// Clear status messages on any key press.
		m.statusMsg = ""
		m.statusErr = ""

		// When filtering, route keys to the text input.
		if m.filtering {
			switch {
			case key.Matches(msg, keys.Escape):
				m.filtering = false
				m.filterInput.SetValue("")
				m.filterInput.Blur()
				return m, nil
			case msg.String() == "enter":
				m.filtering = false
				m.filterInput.Blur()
				m.loading = true
				return m, loadPosts(m.client, 1, m.statusFilter, m.filterInput.Value())
			default:
				var cmd tea.Cmd
				m.filterInput, cmd = m.filterInput.Update(msg)
				return m, cmd
			}
		}

		switch {
		case key.Matches(msg, keys.Escape):
			if m.selected >= 0 {
				m.selected = -1
				return m, nil
			}
			if m.startView == viewPostList {
				return m, tea.Quit
			}
			m.currentView = viewDashboard
			m.filterInput.SetValue("")
			return m, nil

		case key.Matches(msg, keys.Down):
			if len(m.posts) > 0 && m.cursor < len(m.posts)-1 {
				m.cursor++
				m.selected = -1
			}
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
				m.selected = -1
			}
			return m, nil

		case key.Matches(msg, keys.Enter):
			if len(m.posts) > 0 {
				p := m.posts[m.cursor]
				m.currentView = viewPostDetail
				m.loading = true
				return m, loadPostDetail(m.client, p.ID)
			}
			return m, nil

		case key.Matches(msg, keys.Tab):
			idx := 0
			for i, f := range statusFilters {
				if f == m.statusFilter {
					idx = i
					break
				}
			}
			m.statusFilter = statusFilters[(idx+1)%len(statusFilters)]
			m.loading = true
			return m, loadPosts(m.client, 1, m.statusFilter, m.filterInput.Value())

		case key.Matches(msg, keys.Filter):
			m.filtering = true
			m.filterInput.Focus()
			return m, nil

		case msg.String() == "n":
			if m.postPag != nil && m.postPag.Next != nil {
				m.loading = true
				return m, loadPosts(m.client, *m.postPag.Next, m.statusFilter, m.filterInput.Value())
			}
			return m, nil

		case msg.String() == "N":
			if m.postPag != nil && m.postPag.Prev != nil {
				m.loading = true
				return m, loadPosts(m.client, *m.postPag.Prev, m.statusFilter, m.filterInput.Value())
			}
			return m, nil

		case key.Matches(msg, keys.Edit):
			if len(m.posts) > 0 && m.cursor < len(m.posts) {
				p := m.posts[m.cursor]
				m.statusMsg = ""
				m.statusErr = ""
				return m, openGhostEditor(m.client.BaseURL, p.ID)
			}
			return m, nil

		case key.Matches(msg, keys.Open):
			if len(m.posts) > 0 && m.cursor < len(m.posts) {
				p := m.posts[m.cursor]
				if p.URL != "" {
					m.statusMsg = ""
					m.statusErr = ""
					return m, openInBrowser(p.URL)
				}
				m.statusErr = "no URL for this post"
				return m, nil
			}
			return m, nil

		case key.Matches(msg, keys.Toggle):
			if len(m.posts) > 0 && m.cursor < len(m.posts) {
				p := m.posts[m.cursor]
				if p.Status == "scheduled" {
					m.statusErr = "cannot toggle scheduled posts"
					m.statusMsg = ""
					return m, nil
				}
				m.statusMsg = ""
				m.statusErr = ""
				m.loading = true
				return m, togglePostStatus(m.client, p.ID)
			}
			return m, nil
		}
	}
	return m, nil
}

func postListView(m model) string {
	var b strings.Builder
	w := m.contentWidth()

	b.WriteString("\n")

	// Header.
	header := "Posts"
	if m.statusFilter != "all" {
		header += " " + labelStyle.Render("("+m.statusFilter+")")
	}
	b.WriteString(indent(titleStyle.Render(header)))
	b.WriteString("\n")

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

	if len(m.posts) == 0 {
		b.WriteString(indent(labelStyle.Render("no posts found")))
		b.WriteString("\n")
	} else {
		for i, p := range m.posts {
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

			if i == m.cursor {
				b.WriteString(indent(selectedRowStyle.Render(row)))
			} else {
				b.WriteString(indent(normalRowStyle.Render(row)))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	// Pagination.
	if m.postPag != nil {
		pagInfo := fmt.Sprintf("page %d/%d (%d total)", m.postPag.Page, m.postPag.Pages, m.postPag.Total)
		b.WriteString(indent(labelStyle.Render(pagInfo)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(indent(divider(w - 4)))
	b.WriteString("\n\n")

	// Filter input.
	if m.filtering {
		b.WriteString(indent(filterPromptStyle.Render("/ ") + m.filterInput.View()))
		b.WriteString("\n\n")
	}

	// Status feedback.
	if m.statusMsg != "" {
		b.WriteString(indent(publishedStyle.Render(m.statusMsg)))
		b.WriteString("\n\n")
	}
	if m.statusErr != "" {
		b.WriteString(indent(errStyle.Render(m.statusErr)))
		b.WriteString("\n\n")
	}

	// Help — use the appropriate esc label.
	escBinding := keys.Escape
	if m.startView == viewPostList {
		escBinding = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit"))
	}
	b.WriteString(indent(m.help.View(postListKeys{escape: escBinding})))
	b.WriteString("\n")

	return b.String()
}

func postDetailPane(p ghost.Post, width int) string {
	var b strings.Builder
	pad := "      "

	b.WriteString("\n")
	b.WriteString(pad + detailLabelStyle.Render("slug  ") + detailValueStyle.Render(p.Slug) + "\n")
	b.WriteString(pad + detailLabelStyle.Render("status") + " " + statusIndicator(p.Status) + " " + detailValueStyle.Render(p.Status) + "\n")

	if p.PublishedAt != nil {
		b.WriteString(pad + detailLabelStyle.Render("date  ") + detailValueStyle.Render(p.PublishedAt.Format("2006-01-02 15:04")) + "\n")
	}
	if p.URL != "" {
		b.WriteString(pad + detailLabelStyle.Render("url   ") + detailValueStyle.Render(p.URL) + "\n")
	}

	if len(p.Tags) > 0 {
		names := make([]string, len(p.Tags))
		for i, t := range p.Tags {
			names[i] = t.Name
		}
		b.WriteString(pad + detailLabelStyle.Render("tags  ") + tagStyle.Render(strings.Join(names, ", ")) + "\n")
	}

	excerpt := p.CustomExcerpt
	if excerpt == "" {
		excerpt = p.Excerpt
	}
	if excerpt != "" {
		maxLen := width - 10
		if len(excerpt) > maxLen {
			excerpt = excerpt[:maxLen-3] + "..."
		}
		b.WriteString("\n")
		b.WriteString(pad + labelStyle.Render(excerpt) + "\n")
	}

	b.WriteString("\n")

	return b.String()
}
