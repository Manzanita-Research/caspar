# caspar

Go CLI for Ghost CMS. Claude-first — agents are the primary user.

## Dev

```sh
go build -o caspar .   # build
go test ./...           # test
make build              # same as above, plus install
```

## Structure

- `cmd/` — cobra commands. `resource.go` has shared post/page logic.
- `pkg/ghost/` — API client, JWT auth, types for each resource.
- `pkg/config/` — `~/.caspar.json` management.
- `pkg/output/` — JSON vs lipgloss-styled output.
- `pkg/tui/` — interactive bubbletea TUI (dashboard + post list).

## TUI Libraries

Prefer charmbracelet for all terminal UI work:
- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [bubbles](https://github.com/charmbracelet/bubbles) — components (help, textinput, etc.)
- [lipgloss](https://github.com/charmbracelet/lipgloss) — styling and layout
- [huh](https://github.com/charmbracelet/huh) — forms and prompts
- [glamour](https://github.com/charmbracelet/glamour) — markdown rendering

## Conventions

- `--json` flag on every command for structured output.
- `--fields` for token-efficient responses.
- Aliases: `ls`, `show`, `new`, `edit`, `rm`.
- Ghost IDs are 24-char hex strings. Slugs are everything else.
- Updates always fetch current `updated_at` first (Ghost 409 conflict resolution).
- HTML content uses `?source=html` query param.
