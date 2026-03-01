package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/manzanita-research/caspar/pkg/ghost"
)

var (
	detailTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(accentColor)

	detailMetaKeyStyle = lipgloss.NewStyle().
				Foreground(accentColor).
				Bold(true).
				Width(10)

	detailMetaValStyle = lipgloss.NewStyle().
				Foreground(textColor)

	detailURLStyle = lipgloss.NewStyle().
			Foreground(mutedColor)
)

func postDetailUpdate(m model, msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case postDetailMsg:
		m.detailPost = &msg.post
		m.loading = false
		m.detailContent = renderPostContent(msg.post, m.contentWidth())
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape), msg.String() == "backspace":
			m.currentView = viewPostList
			m.detailPost = nil
			m.detailContent = ""
			return m, nil
		}
	}
	return m, nil
}

func postDetailView(m model) string {
	var b strings.Builder
	w := m.contentWidth()

	b.WriteString("\n")

	if m.loading {
		b.WriteString(indent(labelStyle.Render("loading post...")))
		b.WriteString("\n")
		return b.String()
	}

	p := m.detailPost
	if p == nil {
		b.WriteString(indent(labelStyle.Render("no post loaded")))
		b.WriteString("\n")
		return b.String()
	}

	// Title.
	b.WriteString(indent(detailTitleStyle.Render(p.Title)))
	b.WriteString("\n\n")

	// Status + date + author line.
	meta := statusIndicator(p.Status) + " " + detailMetaValStyle.Render(p.Status)
	if p.PublishedAt != nil {
		meta += labelStyle.Render("  ") + detailMetaValStyle.Render(p.PublishedAt.Format("2006-01-02 15:04"))
	}
	if len(p.Authors) > 0 {
		names := make([]string, len(p.Authors))
		for i, a := range p.Authors {
			names[i] = a.Name
		}
		meta += labelStyle.Render("  by ") + detailMetaValStyle.Render(strings.Join(names, ", "))
	}
	b.WriteString(indent(meta))
	b.WriteString("\n")

	// Tags.
	if len(p.Tags) > 0 {
		names := make([]string, len(p.Tags))
		for i, t := range p.Tags {
			names[i] = t.Name
		}
		b.WriteString(indent(detailMetaKeyStyle.Render("tags") + tagStyle.Render(strings.Join(names, ", "))))
		b.WriteString("\n")
	}

	// URL.
	if p.URL != "" {
		b.WriteString(indent(detailMetaKeyStyle.Render("url") + detailURLStyle.Render(p.URL)))
		b.WriteString("\n")
	}

	// Slug.
	b.WriteString(indent(detailMetaKeyStyle.Render("slug") + detailMetaValStyle.Render(p.Slug)))
	b.WriteString("\n")

	// Featured.
	if p.Featured {
		b.WriteString(indent(detailMetaKeyStyle.Render("featured") + detailMetaValStyle.Render("yes")))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(indent(divider(w - 4)))
	b.WriteString("\n\n")

	// Content.
	if m.detailContent != "" {
		// Glamour output comes pre-formatted; indent each line.
		for _, line := range strings.Split(m.detailContent, "\n") {
			b.WriteString(indent(line))
			b.WriteString("\n")
		}
	} else {
		excerpt := p.CustomExcerpt
		if excerpt == "" {
			excerpt = p.Excerpt
		}
		if excerpt != "" {
			b.WriteString(indent(detailMetaValStyle.Render(excerpt)))
			b.WriteString("\n")
		} else {
			b.WriteString(indent(labelStyle.Render("(no content available)")))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(indent(divider(w - 4)))
	b.WriteString("\n\n")

	// Help bar.
	b.WriteString(indent(m.help.View(postDetailKeys{})))
	b.WriteString("\n")

	return b.String()
}

func renderPostContent(p ghost.Post, width int) string {
	// Use excerpt/custom_excerpt as the primary content for rendering,
	// since full HTML-to-markdown conversion is complex. Glamour renders
	// markdown beautifully, so we format the available text as markdown.
	var content strings.Builder

	excerpt := p.CustomExcerpt
	if excerpt == "" {
		excerpt = p.Excerpt
	}

	if excerpt != "" {
		content.WriteString(excerpt)
		content.WriteString("\n")
	}

	if p.FeatureImage != "" {
		content.WriteString(fmt.Sprintf("\n![feature image](%s)\n", p.FeatureImage))
	}

	text := content.String()
	if text == "" {
		return ""
	}

	renderWidth := width - 4
	if renderWidth < 40 {
		renderWidth = 40
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(renderWidth),
	)
	if err != nil {
		return text
	}

	rendered, err := r.Render(text)
	if err != nil {
		return text
	}

	return strings.TrimSpace(rendered)
}
