---
name: caspar
description: CLI for Ghost CMS — manage posts, pages, tags, members, newsletters, and images from the terminal. Use this skill whenever the user wants to interact with their Ghost blog or site, including creating or editing posts, publishing drafts, managing tags or members, uploading images, or querying site content. Also triggers on mentions of Ghost CMS, blog management, "publish this", "create a post", "list my drafts", or any Ghost content operations.
---

# caspar

CLI for Ghost CMS. Designed for agents — use `--json` for all reads and writes. Use `--fields` to keep responses small.

## Auth

```sh
caspar auth login    # interactive setup (prompts for URL + API key)
caspar auth status   # check connection
caspar auth logout   # remove saved credentials
```

Env vars `CASPAR_URL` and `CASPAR_ADMIN_API_KEY` override saved config. Run `caspar auth status` before other commands if you're unsure whether auth is configured.

## Posts and pages

`caspar page` accepts the same subcommands and flags as `caspar post`.

```
caspar post list [--limit N] [--page P] [--filter EXPR] [--fields F] [--order O] [--include tags,authors] [--json]
caspar post get <id-or-slug> [--fields F] [--formats html|html+lexical] [--include tags,authors] [--json]
caspar post create --title T [--html H | --stdin | --lexical JSON] [--status S] [--slug S] [--tag T]... [--featured] [--visibility V] [--published-at ISO8601] [--json]
caspar post update <id-or-slug> [--title T] [--html H | --stdin | --lexical JSON] [--status S] [--slug S] [--tag T]... [--featured | --no-featured] [--visibility V] [--published-at ISO8601] [--custom-excerpt E] [--json]
caspar post delete <id> [--json]
```

### Key flags

- `--status`: `draft` (default), `published`, `scheduled`. Use with `--published-at` for scheduling.
- `--visibility`: `public`, `members`, `paid`, `tiers`
- `--published-at`: ISO 8601 datetime, e.g. `2025-06-15T09:00:00.000Z`. Required when `--status scheduled`.
- `--tag`: repeatable. On update, replaces all existing tags.
- `--featured` / `--no-featured`: toggle the featured flag.
- `--custom-excerpt`: sets the excerpt shown in cards and previews.
- `--formats`: controls which content format `get` returns. Default is html. Use `html+lexical` if you need both.

## Tags

```
caspar tag list   [--limit N] [--page P] [--filter EXPR] [--fields F] [--order O] [--json]
caspar tag get    <id-or-slug> [--fields F] [--json]
caspar tag create --name N [--slug S] [--description D] [--visibility public|internal] [--json]
caspar tag update <id-or-slug> [--name N] [--slug S] [--description D] [--visibility V] [--json]
caspar tag delete <id> [--json]
```

Internal tags (visibility `internal`) are prefixed with `#` in Ghost and hidden from the site. Useful for grouping or automation.

## Members

```
caspar member list   [--limit N] [--page P] [--filter EXPR] [--fields F] [--order O] [--json]
caspar member get    <id> [--fields F] [--json]
caspar member create --email E [--name N] [--label L]... [--json]
caspar member update <id> [--email E] [--name N] [--label L]... [--json]
```

- `--label` is repeatable. On update, replaces all existing labels.
- Member `get` takes an ID, not a slug.

## Newsletters, images, site

```
caspar newsletter list [--json]
caspar newsletter get <id> [--json]
caspar image upload <file> [--json]
caspar site [--json]
```

## Aliases

`list` → `ls`, `get` → `show`, `create` → `new`, `update` → `edit`, `delete` → `rm`

## Newsletter safety

Creating or publishing posts via caspar **never sends newsletter emails**. Ghost's Admin API requires an explicit `newsletter` field to trigger email delivery, and caspar never sets it. Posts default to `draft` status. Even `--status published` just makes the post visible on the site — no emails go out. Bulk operations are safe.

## Working with content

- `--html` is for short inline content. For anything longer than a sentence, pipe through `--stdin` instead.
- `--lexical` accepts Ghost's Lexical JSON format. Prefer `--html` unless you need Lexical-specific features.
- When creating HTML content, Ghost expects block-level elements (`<p>`, `<h2>`, `<ul>`, etc.). Don't pass bare text.

## Pagination

List commands return `--limit` results per page (default 15). For large collections, use `--page` to step through:

```sh
caspar post list --json --fields id,title,status --limit 50 --page 1
caspar post list --json --fields id,title,status --limit 50 --page 2
```

JSON output includes `meta.pagination` with `total`, `pages`, `page`, and `limit`.

## Gotchas

- `get` auto-detects slug vs ID. `delete` requires an ID.
- `--html` with create/update adds `?source=html` automatically.
- `--tag` and `--label` are repeatable and replace all existing values on update.
- Updates fetch `updated_at` automatically for Ghost's 409 conflict resolution.
- `--filter` uses Ghost NQL: `status:published`, `tag:getting-started`, `published_at:>'2024-01-01'`
- `member get` takes an ID only — no slug lookup.

## Common patterns

```sh
# check auth before starting
caspar auth status

# list recent drafts
caspar post list --filter "status:draft" --order "updated_at desc" --limit 5 --json

# efficient listing for token savings
caspar post list --json --fields id,title,slug,status --limit 20

# create draft from stdin
echo "<p>Content here</p>" | caspar post create --title "New Post" --stdin --json

# create a scheduled members-only post
caspar post create --title "Coming Soon" --html "<p>Exciting news.</p>" \
  --status scheduled --published-at "2025-07-01T09:00:00.000Z" \
  --visibility members --json

# publish a draft
caspar post update <id-or-slug> --status published --json

# add tags to a post (replaces existing tags)
caspar post update <slug> --tag "news" --tag "featured" --json

# create an internal tag for automation
caspar tag create --name "#imported" --visibility internal --json

# upload an image and get its URL
caspar image upload ./photo.jpg --json
```
