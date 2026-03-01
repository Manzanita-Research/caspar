package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit     key.Binding
	Posts    key.Binding
	Pages    key.Binding
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Tab      key.Binding
	Filter   key.Binding
	NextPage key.Binding
	PrevPage key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Posts: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "posts"),
	),
	Pages: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "pages"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "expand"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "status"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n/N", "page"),
	),
	PrevPage: key.NewBinding(
		key.WithKeys("N"),
	),
}

// dashboardKeys is the help.KeyMap for the dashboard view.
type dashboardKeys struct{}

func (dashboardKeys) ShortHelp() []key.Binding {
	return []key.Binding{keys.Posts, keys.Pages, keys.Quit}
}

func (dashboardKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{{keys.Posts, keys.Pages, keys.Quit}}
}

// postListKeys is the help.KeyMap for the post list view.
type postListKeys struct {
	escape key.Binding
}

func (k postListKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		keys.Up, keys.Down, keys.Enter, keys.Tab,
		keys.Filter, keys.NextPage, k.escape, keys.Quit,
	}
}

func (k postListKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// pageListKeys is the help.KeyMap for the page list view.
type pageListKeys struct {
	escape key.Binding
}

func (k pageListKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		keys.Up, keys.Down, keys.Enter, keys.Tab,
		keys.Filter, keys.NextPage, k.escape, keys.Quit,
	}
}

func (k pageListKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
