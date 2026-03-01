package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Manzanita brand palette — adaptive for light and dark terminals.
// AdaptiveColor picks Light value on light bg, Dark value on dark bg.
var (
	accentColor = lipgloss.AdaptiveColor{Light: "#C2714F", Dark: "#D4886A"} // terracotta
	titleColor  = lipgloss.AdaptiveColor{Light: "#6B3A2A", Dark: "#D4886A"} // bark / terracotta
	textColor   = lipgloss.AdaptiveColor{Light: "#2C2C2C", Dark: "#E8E5DF"} // warm black / fog
	mutedColor  = lipgloss.AdaptiveColor{Light: "#8C8478", Dark: "#9B9590"} // stone
	sageColor   = lipgloss.AdaptiveColor{Light: "#5A7A4A", Dark: "#8B9E7E"} // published
	ochreColor  = lipgloss.AdaptiveColor{Light: "#A07D2E", Dark: "#C49A3C"} // scheduled
	duskColor   = lipgloss.AdaptiveColor{Light: "#6B6B6B", Dark: "#8C8C8C"} // draft
	rustColor   = lipgloss.AdaptiveColor{Light: "#A0522D", Dark: "#C2714F"} // errors
	lavColor    = lipgloss.AdaptiveColor{Light: "#7B6E88", Dark: "#B0A3BD"} // tags
	fogColor    = lipgloss.AdaptiveColor{Light: "#D5D0C8", Dark: "#4A4A4A"} // dividers
)

const maxWidth = 80

// Status indicators.
const (
	statusPublished = "●"
	statusDraft     = "○"
	statusScheduled = "◐"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(titleColor)
	subtitleStyle = lipgloss.NewStyle().Foreground(mutedColor)
	labelStyle    = lipgloss.NewStyle().Foreground(mutedColor)
	valueStyle    = lipgloss.NewStyle().Foreground(textColor)

	publishedStyle = lipgloss.NewStyle().Foreground(sageColor)
	draftStyle     = lipgloss.NewStyle().Foreground(duskColor)
	scheduledStyle = lipgloss.NewStyle().Foreground(ochreColor)

	selectedRowStyle = lipgloss.NewStyle().Bold(true).Foreground(accentColor)
	normalRowStyle   = lipgloss.NewStyle().Foreground(textColor)

	tagStyle = lipgloss.NewStyle().Foreground(lavColor)

	dividerStyle = lipgloss.NewStyle().Foreground(fogColor)
	errStyle     = lipgloss.NewStyle().Bold(true).Foreground(rustColor)
	countStyle   = lipgloss.NewStyle().Bold(true).Foreground(accentColor)

	filterPromptStyle = lipgloss.NewStyle().Foreground(accentColor)

	detailLabelStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	detailValueStyle = lipgloss.NewStyle().Foreground(textColor)
)

func statusIndicator(status string) string {
	switch status {
	case "published":
		return publishedStyle.Render(statusPublished)
	case "draft":
		return draftStyle.Render(statusDraft)
	case "scheduled":
		return scheduledStyle.Render(statusScheduled)
	default:
		return draftStyle.Render(statusDraft)
	}
}

func indent(s string) string {
	return "  " + s
}

func divider(width int) string {
	if width > maxWidth {
		width = maxWidth
	}
	return dividerStyle.Render(strings.Repeat("─", width))
}
