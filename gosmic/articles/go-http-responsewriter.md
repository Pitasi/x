---
title: "Build your own ResponseWriter: safer HTTP in Go"
date: "2025-04-23"
description: "TODO"
categories:
  - "go"
published: false
---

Go's `http.ResponseWriter` writes directly to the socket, which can lead to
subtle bugs like forgetting to set a status code or accidentally modifying
headers too late.

This article shows how it's possible to wrap the ResponseWriter to enforce
custom rules like requiring `WriteHeader()` and blocking writes after an error,
making your handlers safer and easier to reason about.

I've written hundreds of HTTP handlers in Go and I kept making the same subtle
mistake without realizing it.
It wasn't until [empijei](https://empijei.science/)'s [workshop on secure
code](https://golab.io/talks/a-simple-approach-to-secure-code) that it finally
clicked: `http.ResponseWriter` is unsafe by default, but it's meant to be used
as a base for your own custom logic.

Some takeaway from the workshop were:

- `http.ResponseWriter` is an interface
- you can implement your own `ResponseWriter` by wrapping another
`ResponseWriter`, in order to enforce certain rules

Every handler I've written more or less always starts like this:

```go
http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("helo world"))
})

http.ListenAndServe(":8080", nil)
```

One notable thing is that the `w.Write` call will also call the
`w.WriteHeader` for me, since I didn't do it explicitly myself:

```go
http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(200)
  w.Write([]byte("helo world"))
})
```

`w.WriteHeader` writes the status code of the response, and all the header
entries that were stored in the `w.Header()` map.

What if I want to enforce all my handlers to explicitly set the status code
even if it's `200 OK`, so I can be sure I didn't forget to set it?

Additionally, setting new headers or changing existing ones won't have any
effect and you won't receive any sort of errors:

```go
http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(200)
  w.Header().Set("content-type", "application/json") // silently has no effects
  w.Write([]byte("helo world"))
})
```

Wouldn't it be great to have at least a warning that we're doing something by mistake?

Let's see another example:

```go
http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
  response, err := database.LoadResponse()
  if err != nil {
    log.Println("error loading response:", err)
    w.WriteHeader(500)
    w.Write([]byte("error :("))
  }

  w.Write([]byte("response: "))
  w.Write(response)
})
```

There's a bug in the above code that I've made more than once in my life.

Can you spot it?

Solution in...

3...

2...

1...

It's missing an early return from the `if` condition. The rest of the handler
printing the response will continue to execute even if the database returned an
error!

Great. What can we do about it?

One of the problems is that the `http.ResponseWriter` that we are using, is an
actual writer, that writes in the underlying TCP socket without *preparing* a
response entirely before writing it.

`http.ResponseWriter` is an interface. Let's implement it ourselves enforcing
our custom rules:

```go
type HttpWriter struct {
  w http.ResponseWriter // wrap an existing writer
}

func NewHttpWriter(w http.ResponseWriter) http.ResponseWriter {
  return &HttpWriter{
    w: w,
  }
}
```

We need to implement only a handful of methods:

```go
// implement http.ResponseWriter

func (w *HttpWriter) Header() http.Header {
  return w.w.Header()
}

func (w *HttpWriter) Write(data []byte) (int, error) {
  return w.w.Write(data)
}

func (w *HttpWriter) WriteHeader(statusCode int) {
  w.w.WriteHeader(statusCode)
}
```

```go
// it's actually a good idea to implement a Flusher version for our writer as
// well

type HttpWriterFlusher struct {
  *HttpWriter   // wrap our writer
  http.Flusher  // keep a ref to the wrapped Flusher
}

func (w *HttpWriterFlusher) Flush() {
  w.Flusher.Flush()
}

// modify the constructor to either return HttpWriter or HttpWriterFlusher
// depending on the writer being wrapped

func NewHttpWriter(w http.ResponseWriter) http.ResponseWriter {
  httpWriter := &HttpWriter{
    w: w,
  }

  if flusher, ok := w.(http.Flusher); ok {
    return &HttpWriterFlusher{
      HttpWriter: httpWriter,
      Flusher:    flusher,
    }
  }

  return httpWriter
}
```

Let's wire everything up.

You can easily start using `HttpWriter` instead of whichever default
`http.ResponseWriter` was being used, by writing a middleware:

```go
func middleware(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    writer := NewHttpWriter(w)
    h.ServeHTTP(writer, r)
  })
}
```

And now, it's time to have some fun by customizing the implementation of `HttpWriter`.

Want a warning log every time you invoke `Write()` without `WriteHeader()`? You can!

```go
type HttpWriter struct {
  w http.ResponseWriter // wrap an existing writer

  headerWritten bool
}

func (w *HttpWriter) Write(data []byte) (int, error) {
  if !w.headerWritten {
    log.Println("warn: invoked Write() without WriteHeader(statusCode)")
  }
  return w.w.Write(data)
}

func (w *HttpWriter) WriteHeader(statusCode int) {
  w.w.WriteHeader(statusCode)
  w.headerWritten = true
}
```

Want to avoid writing anything at all if the status code has been set to 500?

```go
type HttpWriter struct {
  w http.ResponseWriter // wrap an existing writer

  statusCode int
}

func (w *HttpWriter) Write(data []byte) (int, error) {
  if w.statusCode >= 500 {
    log.Println("warn: ignoring Write(), status code is 500")
    return 0, nil
  }
  return w.w.Write(data)
}

func (w *HttpWriter) WriteHeader(statusCode int) {
  w.w.WriteHeader(statusCode)
  w.statusCode = statusCode
}
```

It's really up to you! You can tweak your rules however you wish.
`http.ResponseWriter` is yours to hack.

---

Thanks to [empijei](https://empijei.science/) and
[loresuso](https://github.com/loresuso) for reading an early draft of this
article ❤️
