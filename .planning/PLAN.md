# caspar — Plan

## Context

Ghost CMS has no good agent-friendly interface. The existing [Ghost MCP Server](https://github.com/jgardner04/Ghost-MCP-Server) dumps 34 tool definitions into your context window, requires Node.js, and is brittle to set up. We want the [linctl](https://github.com/dorkitude/linctl) pattern instead: a clean Go CLI that Claude shells out to, with a SKILL.md that costs ~400 tokens instead of ~4000.

This is a Claude-first CLI — agents are the primary user, humans are welcome too. v2 adds a Bubbletea TUI for the human experience.

---

## Milestone 1: MVP CLI + SKILL.md

### Phase 1: Skeleton + Auth

**Files:**
- `main.go` — entry point
- `go.mod` — `github.com/manzanita-research/caspar`
- `cmd/root.go` — Cobra root, global flags (`--json`, `--version`)
- `cmd/auth.go` — `auth login`, `auth status`, `auth logout`
- `cmd/site.go` — `site` (GET /site/, doubles as auth check)
- `pkg/ghost/client.go` — HTTP client, JWT header injection, error parsing
- `pkg/ghost/auth.go` — JWT generation from `id:secret` key (HS256, 5min expiry, `/admin/` audience)
- `pkg/config/config.go` — read/write `~/.caspar.json` (url + admin_api_key, 0600 perms)
- `pkg/output/output.go` — JSON vs plaintext rendering

**Dependencies:** `spf13/cobra`, `golang-jwt/jwt/v5`, `charmbracelet/bubbletea`, `charmbracelet/lipgloss`, `charmbracelet/huh` (forms/prompts), `charmbracelet/bubbles` (reusable components)

Charm stack is in from day one — Huh for interactive prompts (auth login), Lipgloss for styled human output, Bubbles for tables/spinners. The full interactive TUI dashboard (`caspar tui`) comes in Milestone 2.

**Auth flow:** interactive prompts for URL + API key → validate via GET /site/ → save to `~/.caspar.json`. Env vars `GHOSTCTL_URL` and `GHOSTCTL_ADMIN_API_KEY` override config file.

### Phase 2: Posts + Pages

**Files:**
- `pkg/ghost/types.go` — Post, Page, Tag, Member, Pagination, Error types
- `pkg/ghost/posts.go` — List, Get, Create, Update, Delete
- `pkg/ghost/pages.go` — same interface, different endpoint
- `cmd/post.go` — `post list|get|create|update|delete`
- `cmd/page.go` — `page list|get|create|update|delete`

**Shared list flags:** `--limit`, `--page`, `--filter` (Ghost NQL), `--order`, `--fields`, `--include`
**Post/page flags:** `--title`, `--html`, `--lexical`, `--status`, `--slug`, `--tag` (repeatable), `--featured`, `--stdin`
**Aliases:** `list→ls`, `get→show`, `create→new`, `update→edit`, `delete→rm`

**Key detail:** `get` auto-detects slug vs UUID. Updates require sending current `updated_at` (Ghost 409 conflict resolution).

### Phase 3: Tags + Members + Images + Newsletters

**Files:**
- `pkg/ghost/tags.go`, `cmd/tag.go` — full CRUD
- `pkg/ghost/members.go`, `cmd/member.go` — list, get, create, update
- `pkg/ghost/images.go`, `cmd/image.go` — multipart upload
- `pkg/ghost/newsletters.go`, `cmd/newsletter.go` — read-only (list, get)

### Phase 4: SKILL.md + Polish

**Files:**
- `SKILL.md` — agent interface contract (~400 tokens). Structure: quick rules, gotchas, auth, core workflow, command map, common patterns, troubleshooting
- `README.md` — human-facing docs
- `Makefile` — build, test, install targets
- `CLAUDE.md` — project-specific dev instructions

**Tests:** JWT generation, client request building, output formatting, slug-vs-ID detection

---

## Milestone 2: Bubbletea TUI (fast-follow)

- `caspar tui` command launches interactive TUI
- Post list with filtering, preview pane, status indicators
- Inline editing (open in `$EDITOR`)
- Uses `/tui-design` skill + `brand-tokens` for Manzanita visual identity
- Charm stack: `bubbletea`, `lipgloss`, `bubbles`

## Milestone 3: Distribution + Extended (as needed)

- Homebrew formula, GitHub Actions release builds
- Tier/offer management, webhooks, themes
- Bulk operations, content migration helpers

---

## Command Tree (MVP)

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

---

## Key Design Decisions

1. **Write Ghost client from scratch** — existing Go libs (libecto, go-ghost) are undermaintained. The API surface is small (~500 lines of client code). Keeps us dependency-light.
2. **No viper** — overkill for two fields. Plain `encoding/json` + `os.Getenv` to `~/.caspar.json`.
3. **Charm stack from day one** — Huh for interactive prompts, Lipgloss for styled output, Bubbles for tables/spinners. Full TUI dashboard (`caspar tui`) is Milestone 2.
4. **`--json` is the agent interface** — structured output, no parsing ambiguity. Human mode gets Lipgloss-styled output with Bubbles tables.
5. **`--fields` for token efficiency** — agents can request only `id,title,slug,status` instead of full post bodies.
6. **`--stdin` for content** — avoids putting HTML in flag values, keeps command tokens low.
7. **SKILL.md as the product** — this is what differentiates caspar from an MCP server. ~400 tokens vs ~4000. Agent reads it once, then shells out.

---

## Linear Issues to Create

### Milestone 1: MVP
1. **caspar: scaffold project + auth flow** — go mod, cobra root, auth login/status/logout, site info, config, JWT generation
2. **caspar: post + page CRUD** — types, client methods, cobra commands, flags, slug detection, stdin support
3. **caspar: tags, members, images, newsletters** — remaining resource commands
4. **caspar: SKILL.md + README + tests** — agent skill file, docs, Makefile, test coverage

### Milestone 2: TUI
5. **caspar: bubbletea TUI for interactive use** — `caspar tui`, post browser, preview, editing

### Milestone 3: Distribution
6. **caspar: homebrew + CI/CD releases** — formula, GitHub Actions, cross-platform builds

---

## Verification

1. `caspar auth login` → prompts for URL + key → saves config → prints site title
2. `caspar auth status` → confirms connection, prints site info
3. `caspar post list --json --limit 5` → returns structured JSON array
4. `caspar post create --title "Test" --html "<p>Hello</p>" --json` → creates draft, returns post object
5. `caspar post get <slug> --json` → fetches by slug
6. `caspar post update <id> --status published --json` → publishes, returns updated post
7. `caspar post delete <id>` → deletes, confirms
8. `echo "<p>Piped content</p>" | caspar post create --title "Stdin Test" --stdin --json` → creates from stdin
9. `caspar image upload ./test.jpg --json` → returns image URL
10. Install SKILL.md into `.claude/skills/` → Claude agent can discover and use caspar with ~400 tokens of context
