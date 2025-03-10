# togo

`togo` is a file-system based webserver with Markdown and templating support.

Think of [Hugo](https://gohugo.io) or [Jekyll](https://jekyllrb.com/) but
without a build step.

Check [togo-example](https://github.com/Pitasi/x/tree/main/togo-example) for a
sample website built using `togo`.

## Usage

As a library in your Go project:

```go
import "anto.pt/x/togo"
```

As a CLI:

```go
go install anto.pt/x/togo/cmd/togo
```

## Features

- Templating: .html files are rendered using Go's [html/template](https://pkg.go.dev/html/template)
- Markdown: .md files are rendered using [goldmark](https://github.com/yuin/goldmark) (CommonMark compliant)
- Layouts: _layout.html can be used an html wrapper to avoid repeating common
  boilerplate (+ can be nested)
- Dynamic paths: a path can have one or more parameters by naming your folder
  or file using a `:` prefix (e.g. `:slug`)
- Static files: the `static` folder is server as-is, setting the Content-Type
  and Cache headers

## Possible future works

- More templating features and functions (right now you can only access the
  path parameters)
- A database integration. I'd like to play a bit more around that, maybe
  supporting one single DBMS (i.e. PostgreSQL) and exposing querying
  capabilities as templating functions
- More filetypes support: why being limited at `.html`s and `.md`s? Can we
  embed `.wasm`s implementing a known interface?

