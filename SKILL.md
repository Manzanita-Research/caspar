# ghostctl

CLI for Ghost CMS. Use `--json` for structured output. Use `--fields` to limit response size.

## Auth

```sh
ghostctl auth login   # interactive setup
ghostctl auth status  # check connection
```

Env vars `GHOSTCTL_URL` and `GHOSTCTL_ADMIN_API_KEY` override saved config.

## Commands

```
ghostctl post list [--limit N] [--filter EXPR] [--fields F] [--order O] [--include tags,authors] [--json]
ghostctl post get <id-or-slug> [--fields F] [--include tags,authors] [--json]
ghostctl post create --title T [--html H | --stdin] [--status S] [--slug S] [--tag T]... [--featured] [--json]
ghostctl post update <id-or-slug> [--title T] [--html H | --stdin] [--status S] [--slug S] [--tag T]... [--json]
ghostctl post delete <id> [--json]

ghostctl page   — same subcommands as post
ghostctl tag    list|get|create|update|delete
ghostctl member list|get|create|update
ghostctl newsletter list|get
ghostctl image  upload <file> [--json]
ghostctl site   [--json]
```

Aliases: `list→ls` `get→show` `create→new` `update→edit` `delete→rm`

## Gotchas

- `get` auto-detects slug vs ID. `delete` requires an ID.
- `--html` with create/update adds `?source=html` automatically.
- `--stdin` reads HTML from stdin — use for long content instead of `--html` flag.
- `--tag` is repeatable and replaces all existing tags on update.
- Updates fetch `updated_at` automatically for Ghost's 409 conflict resolution.
- `--filter` uses Ghost NQL: `status:published`, `tag:getting-started`, `published_at:>'2024-01-01'`

## Common patterns

```sh
# list recent drafts
ghostctl post list --filter "status:draft" --order "updated_at desc" --limit 5 --json

# create draft from stdin
echo "<p>Content here</p>" | ghostctl post create --title "New Post" --stdin --json

# publish a draft
ghostctl post update <id-or-slug> --status published --json

# efficient listing for token savings
ghostctl post list --json --fields id,title,slug,status --limit 20
```
