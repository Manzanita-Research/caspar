# caspar

CLI for Ghost CMS. Use `--json` for structured output. Use `--fields` to limit response size.

## Auth

```sh
caspar auth login   # interactive setup
caspar auth status  # check connection
```

Env vars `CASPAR_URL` and `CASPAR_ADMIN_API_KEY` override saved config.

## Commands

```
caspar post list [--limit N] [--filter EXPR] [--fields F] [--order O] [--include tags,authors] [--json]
caspar post get <id-or-slug> [--fields F] [--include tags,authors] [--formats html] [--json]
caspar post create --title T [--html H | --stdin] [--status S] [--slug S] [--tag T]... [--featured] [--published-at ISO8601] [--visibility V] [--json]
caspar post update <id-or-slug> [--title T] [--html H | --stdin] [--status S] [--slug S] [--tag T]... [--published-at ISO8601] [--visibility V] [--custom-excerpt E] [--json]
caspar post delete <id> [--json]

caspar page   — same subcommands as post
caspar tag    list|get|create|update|delete
caspar member list|get|create|update
caspar newsletter list|get
caspar image  upload <file> [--json]
caspar site   [--json]
```

Aliases: `list→ls` `get→show` `create→new` `update→edit` `delete→rm`

## Newsletter safety

Creating or publishing posts via caspar **never sends newsletter emails**. Ghost's Admin API requires an explicit `newsletter` field in the request body to trigger email delivery, and caspar never sets it. Posts default to `draft` status. Even `--status published` just makes the post visible on the site — no emails go out. Bulk imports are safe.

## Gotchas

- `get` auto-detects slug vs ID. `delete` requires an ID.
- `--html` with create/update adds `?source=html` automatically.
- `--stdin` reads HTML from stdin — use for long content instead of `--html` flag.
- `--tag` is repeatable and replaces all existing tags on update.
- `--published-at` sets the publication date (ISO 8601, e.g. `2017-04-25T00:00:00.000Z`). Works on both create and update.
- `--visibility` controls who can see the post/page: `public`, `members`, `paid`, or `tiers`. Works on both create and update.
- `--custom-excerpt` sets a custom excerpt for previews and cards (update only).
- `--formats` on `get` requests specific content formats (e.g. `html`, `html+lexical`).
- Updates fetch `updated_at` automatically for Ghost's 409 conflict resolution.
- `--filter` uses Ghost NQL: `status:published`, `tag:getting-started`, `published_at:>'2024-01-01'`
- Bulk fetch by IDs: `--filter "id:[id1,id2,id3]"` returns multiple posts/pages in one request.

## Common patterns

```sh
# list recent drafts
caspar post list --filter "status:draft" --order "updated_at desc" --limit 5 --json

# create draft from stdin
echo "<p>Content here</p>" | caspar post create --title "New Post" --stdin --json

# publish a draft
caspar post update <id-or-slug> --status published --json

# set a publication date (e.g. migrating old content)
caspar post update <id-or-slug> --published-at "2017-04-25T00:00:00.000Z" --json

# make a post members-only
caspar post update <id-or-slug> --visibility members --json

# efficient listing for token savings
caspar post list --json --fields id,title,slug,status --limit 20
```
