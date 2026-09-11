# Resolving Architecture Common Library (Go)

A Go port of
[`ra-common-java`](https://github.com/resolvingarchitecture/ra-common-java) —
the foundational types for the Resolving Architecture / 1M5 ecosystem,
alongside the TypeScript, Python, C#, and C++ `ra-common-*` ports.

**Zero dependencies** — `go.mod` lists nothing beyond the Go toolchain.
Go's standard library already covers JSON, SHA-1/256/512, HMAC, base64, and
cross-platform home/config/cache directory resolution, so this port needed
none of `ra-common-cpp`'s vendored headers or hand-rolled crypto.

## Packages

- **`messaging`** — `Envelope`, the `Message` hierarchy (`DocumentMessage` /
  `CommandMessage` / `EventMessage` / `TextMessage`), and the producer/
  consumer/channel/bus interfaces. These live together in one package
  deliberately — see "No `ToJson()`/`FromJson()` methods" and "One package,
  not several" below.
- **`route`** — routing slips (`DynamicRoutingSlip`) and external/relayed routes.
- **`service`** — the `Service` interface and `ServiceCore` shared state.
- **`identity`** — `Did`, `PublicKey`, `Signature`.
- **`crypto`** — `Hash`, `Multihash`, `HashCash`, password hashing (stdlib
  `crypto/sha1`, `crypto/sha256`, `crypto/sha512`, `crypto/hmac`).
- **`content`** — typed content (text / html / json / image / audio / video / binary).
- **`tasks`** — `Task` + a goroutine-based `Runner`.
- **`config`** — `.properties` loading + cross-platform directory resolution
  (`os.UserHomeDir`/`os.UserConfigDir`/`os.UserCacheDir`).
- **`lifecycle`**, **`network`**, **`servicestatus`**, **`multipart`**,
  **`encoding`** (base32/58 — base64 is `encoding/base64`), **`util`**
  (bytes/strings/version/random/`UniqueID`/`Nonce`), **`raerror`**.

Serialization is JSON-based (`encoding/json`) and **not** wire-compatible
with the Java version.

## No `ToJson()`/`FromJson()` methods

Every other port in this family writes an explicit `ToJson()`/`FromJson()`
pair per type. Go doesn't need that: exported struct fields with `json:"..."`
tags are marshaled/unmarshaled directly by `encoding/json.Marshal`/
`Unmarshal` — that *is* Go's idiomatic equivalent, and writing a manual
`ToJson()` method here would just be re-implementing what the standard
library already does correctly. The only place custom `MarshalJSON`/
`UnmarshalJSON` methods exist is on the two polymorphic hierarchies
(`Message`, `route.Route`), which need a `type`/`kind` discriminator field
injected on the way out and a dispatcher on the way in — see `DESIGN.md`.

## Use

```go
import "github.com/resolvingarchitecture/ra-common-go/messaging"

e := messaging.DocumentEnvelope()
e.AddRoute("ra.http.HttpService", "SEND")
e.AddContent(map[string]any{"hello": "world"})
e.Ratchet()

fmt.Println(*e.GetRoute().Meta().Service) // "ra.http.HttpService"

text, _ := e.ToJSONString()
back, _ := messaging.EnvelopeFromJSONString(text)
fmt.Println(back.Equals(e)) // true
```

## Develop

```sh
go build ./...
go vet ./...
go test ./...
go test -race ./...   # especially the tasks package
```

## Status

Phase 1 (core), matching the other ports' scope. Deferred: currency, locale/
i18n, the full network service layer, `Protocol`, shell/file/browser
utilities, `InfoVault`. See `TODO.md`.

## Reference

- [`DESIGN.md`](DESIGN.md)
- [`TODO.md`](TODO.md)

## License

MIT — see `LICENSE`.
