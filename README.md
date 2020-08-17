Kudosaurus
==========

How to fire it up
-----------------

- Spin up ngrok or a similar proxy.
- Go to `https://api.slack.com/apps` (assuming you're logged into your Slack account) and create an app.
- Give it the following scopes: `chat:write`, `chat:write.public`, `commands`, `incoming-webhook`, `users:read` and `users:read.email`. Remember to save changes in this and following config steps.
- In "Interactivity & Shortcuts", set `[YOUR ADDRESS]/interactivity` as request URL.
- In "Slash Commands", add a `/kudo` command and set request URL to `[YOUR ADDRESS]/slash`.
- In "App Home" > "Show Tabs", enable "Home Tab".
- Rename `config.yml.example` to `config.yml` and edit it to provide the tokens and secrets.
- `go run cmd/main.go`
- In "Event Subscriptions", enable events and set request URL to `[YOUR ADDRESS]/events`.

You're good to go. It's possible that app home tab will fail if you don't `/kudo` someone before using it.


To do
-----

[ ] better structure
[ ] drop nlopes/slack lib - it's not that useful when you're not using RTM and its API coverage isn't :100:
[ ] admin's "preview all kudos" modal - some work has been done, but it's super incomplete
[ ] additional functionality
