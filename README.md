# Loggo

An opinionated, batteries-included structured logging library, built on top of [go-kit's logger](https://github.com/go-kit/kit)

## Install

`go get -u github.com/nubunto/loggo`

## Getting started

```go
package main

import (
  "github.com/nubunto/loggo"
)

func main() {
  l := loggo.New(
    loggo.JSON(os.Stdout),
  )
  l.Log("hello", "world") // logs '{"hello": "world"}' to STDOUT
}
```

The bulk of this library lives inside the `adapters` package, in which an adapter for ElasticSearch is builtin:

```go
package main

import (
  "github.com/nubunto/loggo"
  es "github.com/nubunto/loggo/adapters"
)

func main() {
  esHandler, err := adapters.NewElasticHandler(
    es.Type("logger"),
    es.Index("logging"),
    es.Client(), // default ES client
  )
  if err != nil {
    // ...
  }
  esAdapter, err := adapters.NewAdapter(esHandler)
  if err != nil {
    // ...
  }
  l := loggo.New(
    loggo.JSON(esAdapter),
  )
  
  // log directly to elastic search
  l.Log("this-key", "goes to elastic search under index 'logging' and type 'logger'")
}
```

Custom or new adapters can be coded by satisfying the `adapters.Handler` interface:

```go
type Handler interface {
  HandlePayload([]byte) error
}
```

And you use it by supplying it to the `adapters.NewAdapter` function.

## Composability

Since this plays with go-kit's `log.Logger`, you can pretty much default to it for any missing feature, such as levels or passing through the standard logging package.

## Contributing

Tests expect a ElasticSearch listening on 9200, the default port.
