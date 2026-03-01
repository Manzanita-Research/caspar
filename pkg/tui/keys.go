package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit     key.Binding
	Posts    key.Binding
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Escape   key.Binding
	Tab      key.Binding
	Filter   key.Binding
	NextPage key.Binding
	PrevPage key.Binding
	Edit     key.Binding
	Open     key.Binding
	Toggle   key.Binding
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
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Open: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "draft/publish"),
	),
}

// dashboardKeys is the help.KeyMap for the dashboard view.
type dashboardKeys struct{}

func (dashboardKeys) ShortHelp() []key.Binding {
	return []key.Binding{keys.Posts, keys.Quit}
}

func (dashboardKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{{keys.Posts, keys.Quit}}
}

// postListKeys is the help.KeyMap for the post list view.
type postListKeys struct {
	escape key.Binding
}

func (k postListKeys) ShortHelp() []key.Binding {
	return []key.Binding{
		keys.Up, keys.Down, keys.Enter, keys.Tab,
		keys.Edit, keys.Open, keys.Toggle,
		keys.Filter, keys.NextPage, k.escape, keys.Quit,
	}
}

func (k postListKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
