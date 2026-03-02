# caspar

The friendly Ghost CLI. Agents are the primary user — humans are welcome too.

Named after Caspar, the tiny town in Mendocino County, California — and yeah, the friendly ghost too.

## Why

Ghost has no good agent-friendly interface. The existing MCP server dumps 34 tool definitions into your context window, requires Node.js, and is brittle to set up. Caspar takes a different approach: a clean Go CLI that Claude shells out to, with a [skill file](./cmd/skill.md) that costs ~400 tokens instead of ~4000.

Built on the [linctl](https://github.com/dorkitude/linctl) pattern.

## Install

```sh
go install github.com/manzanita-research/caspar@latest
```

## Setup

```sh
caspar auth login
```

Prompts for your Ghost site URL and admin API key, validates the connection, and saves config to `~/.caspar.json`.

You can also set environment variables:

```sh
export CASPAR_URL=https://your-site.ghost.io
export CASPAR_ADMIN_API_KEY=your-id:your-secret
```

## Usage

```
caspar
├── auth login|status|logout
├── site
├── post list|get|create|update|delete
├── page list|get|create|update|delete
├── tag  list|get|create|update|delete
├── member list|get|create|update
├── newsletter list|get
└── image upload
```

### For agents

Pass `--json` for structured output. Use `--fields` to request only what you need and keep token costs low.

```sh
caspar post list --json --fields id,title,slug,status --limit 10
```

Pipe content in with `--stdin` to avoid putting HTML in flag values:

```sh
echo "<p>Hello world</p>" | caspar post create --title "New Post" --stdin --json
```

### For humans

Without `--json`, output is styled with [Charm](https://charm.sh) — tables, colors, the works.

```sh
caspar post list --limit 10
caspar post get my-post-slug
caspar tag list
```

### Common flags

| Flag | Description |
|------|-------------|
| `--json` | Structured JSON output |
| `--fields` | Comma-separated field list |
| `--limit` | Number of results |
| `--filter` | Ghost NQL filter expression |
| `--order` | Sort order |
| `--stdin` | Read content from stdin |

### Aliases

`list → ls` · `get → show` · `create → new` · `update → edit` · `delete → rm`

## Agent integration

```sh
caspar skill install
```

Copies the built-in [skill file](./cmd/skill.md) to `~/.claude/skills/caspar/SKILL.md`. Claude reads it once (~400 tokens), then knows how to use caspar for any Ghost content task.

## Status

Under active development. We're building this because we needed it.

## License

MIT
