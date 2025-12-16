# hooker

This service runs a little HTTP server that can receive incoming webhook
requests from GitHub and executes a script/program named after the repository
that generated the event.

Example: I push to `Pitasi/foo`, `scripts/pitasi_foo` is executed (if exists).

## Usage

```sh
go run . -addr :9000 -secret 'foo'
```

Then, setup the webhooks from the GitHub repo settings.
