httpcc
======

Parses HTTP/1.1 Cache-Control header, and returns a struct that is convenient
for the end-user to do what they will with.

# Parsing the HTTP Request

```go
dir, err := httpcc.ParseRequest(res)
// dir.MaxAge       uint
// dir.MaxStale     uint
// dir.MinFresh     uint
// dir.NoCache      bool
// dir.NoStore      bool
// dir.NoTransform  bool
// dir.OnlyIfCached bool
// dir.Extensions   map[string]string
```

# Parsing the HTTP Response

```go
directives, err := httpcc.ParseResponse(res)
// dir.MaxAge         uint
// dir.MustRevalidate bool
// dir.NoCache        []string
// dir.NoStore        bool
// dir.NoTransform    bool
// dir.Public         bool
// dir.Private        bool
// dir.SMaxAge        uint
// dir.Extensions     map[string]string
```

