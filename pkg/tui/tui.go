package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/manzanita-research/caspar/pkg/ghost"
)

type view int

const (
	viewDashboard view = iota
	viewPostList
	viewPostDetail
	viewPageList
)

type model struct {
	currentView view
	startView   view

	client *ghost.Client
	site   *ghost.SiteInfo
	err    error
	width  int
	height int
	ready  bool

	// Dashboard counts.
	postCount   int
	pageCount   int
	memberCount int
	tagCount    int

	// Post list.
	posts        []ghost.Post
	postPag      *ghost.Pagination
	cursor       int
	statusFilter string
	filterInput  textinput.Model
	filtering    bool
	loading      bool
	selected     int

	// Post detail.
	detailPost    *ghost.Post
	detailContent string

	// Page list.
	pages         []ghost.Page
	pagePag       *ghost.Pagination
	pageCursor    int
	pageExpanded  int
	pageStatusIdx int
	pageFilter    string
	pageFiltering bool

	// Help.
	help help.Model
}

func initialModel(client *ghost.Client, startView view) model {
	ti := textinput.New()
	ti.Placeholder = "NQL filter (e.g. featured:true)"
	ti.CharLimit = 256
	ti.Width = 60
	ti.PromptStyle = filterPromptStyle
	ti.TextStyle = lipgloss.NewStyle().Foreground(textColor)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(mutedColor)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(accentColor)

	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Bold(true).Foreground(accentColor)
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(mutedColor)
	h.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(fogColor)

	return model{
		currentView:  startView,
		startView:    startView,
		client:       client,
		statusFilter: "all",
		filterInput:  ti,
		selected:     -1,
		pageExpanded: -1,
		loading:      true,
		help:         h,
	}
}

func (m model) Init() tea.Cmd {
	switch m.startView {
	case viewPostList:
		return loadPosts(m.client, 1, m.statusFilter, "")
	default:
		return loadDashboard(m.client)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width - 4 // account for indent
		return m, nil

	case errMsg:
		m.err = msg.err
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		if !m.filtering && !m.pageFiltering && msg.String() == "q" {
			return m, tea.Quit
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.currentView {
	case viewDashboard:
		return dashboardUpdate(m, msg)
	case viewPostList:
		return postListUpdate(m, msg)
	case viewPostDetail:
		return postDetailUpdate(m, msg)
	case viewPageList:
		return pageListUpdate(m, msg)
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return "\n" + indent(errStyle.Render(fmt.Sprintf("error: %s", m.err))) + "\n\n" +
			indent(m.help.View(dashboardKeys{})) + "\n"
	}

	switch m.currentView {
	case viewDashboard:
		return dashboardView(m)
	case viewPostList:
		return postListView(m)
	case viewPostDetail:
		return postDetailView(m)
	case viewPageList:
		return pageListView(m)
	}

	return ""
}

func (m model) contentWidth() int {
	w := m.width
	if w == 0 || w > maxWidth {
		w = maxWidth
	}
	return w
}

// Run starts the TUI.
func Run(client *ghost.Client, startView string) error {
	sv := viewDashboard
	if startView == "posts" {
		sv = viewPostList
	}

	p := tea.NewProgram(initialModel(client, sv), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
