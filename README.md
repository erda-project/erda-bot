# Erda Bot

Erda bot for GitHub and others.

## Instrument

### Cherry-Pick

Add comment in pr like: `/cherry-pick release/1.2 release/1.1 release/1.0`

### Approve

Add comment in pr: `/approve`.

Bot will auto add `approved` label and try to merge continuously until pr can be merged(all checks passed).

### Assign

Add comment in pr like: `/approve @sfwn sfwn effet` (with or without `@` is both supported)

Bot will auto add specified users as pr reviewers.

## GitHub Star

Auto send a DingTalk message when star count added.

You should set two environment variables:

- `DINGTALK_ACCESS_TOKEN`
- `DINGTALK_SECRET` (optional)
