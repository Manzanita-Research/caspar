# ghostctl

Go CLI for Ghost CMS. Claude-first — agents are the primary user.

## Dev

```sh
go build -o ghostctl .   # build
go test ./...            # test
make build               # same as above, plus install
```

## Structure

- `cmd/` — cobra commands. `resource.go` has shared post/page logic.
- `pkg/ghost/` — API client, JWT auth, types for each resource.
- `pkg/config/` — `~/.ghostctl.json` management.
- `pkg/output/` — JSON vs lipgloss-styled output.

## Conventions

- `--json` flag on every command for structured output.
- `--fields` for token-efficient responses.
- Aliases: `ls`, `show`, `new`, `edit`, `rm`.
- Ghost IDs are 24-char hex strings. Slugs are everything else.
- Updates always fetch current `updated_at` first (Ghost 409 conflict resolution).
- HTML content uses `?source=html` query param.
