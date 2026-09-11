# ra-common-go — design notes

A Go port of [`ra-common-java`](https://github.com/resolvingarchitecture/ra-common-java),
tracking the same Phase 1 scope as the Rust, Python, TypeScript, C#, and C++
`ra-common-*` ports. First Go repo in this ecosystem.

## Decisions

- **No `ToJson()`/`FromJson()` methods on plain types.** Every other port
  writes an explicit serialization pair per type. In Go, exported fields with
  `json:"snake_case_name"` tags are marshaled/unmarshaled directly by
  `encoding/json` - writing a manual method would just re-implement what the
  stdlib does. Only the two polymorphic hierarchies need custom
  `MarshalJSON`/`UnmarshalJSON`: each concrete `Message`/`route.Route` type
  marshals itself via `type alias T; json.Marshal(struct{Type string; alias}{...})`
  to inject the wire discriminator, and a free-function dispatcher
  (`UnmarshalMessage`/`UnmarshalRoute`) peeks at that field to pick which
  concrete type to unmarshal into.
- **String-backed types stand in for enums.** `type HashAlgorithm string` with
  `const Sha256 HashAlgorithm = "Sha256"` etc. serializes as a plain JSON
  string with zero custom code, unlike an `iota`-based int enum (which would
  need a hand-written `MarshalJSON`). This is the same wire shape the other
  ports use for their enums, for free.
- **Stdlib crypto and JSON - zero third-party dependencies** (`go.mod` has no
  `require` lines at all). `crypto/sha1`, `crypto/sha256`, `crypto/sha512`,
  `crypto/hmac` cover everything `ra-common-cpp` had to implement by hand;
  PBKDF2 isn't in the stdlib (it's in the separate `golang.org/x/crypto`
  module), but its outer loop is ~15 lines on top of `crypto/hmac`, short
  enough to hand-roll rather than add a dependency for it.
- **`messaging` holds `Message`, `Envelope`, and the producer/consumer/
  channel/bus interfaces together, in one package** - not split the way the
  other ports split them into files. Every other port has the same mutual
  dependency (`Envelope` embeds `Message`; the bus/channel interfaces take an
  `Envelope` parameter) but gets away with splitting it into separate
  files/modules because their import systems tolerate circular *type-only*
  references (TS's `import type`, or just multiple types in one compiled
  unit). Go's package import graph has **no such exception** - it rejects
  cycles outright, including interface-only ones. `servicestatus` stays a
  separate leaf package (`ServiceLevel`, `ServiceStatus`, `ServiceReport`)
  specifically so `messaging.Envelope` can depend on `ServiceLevel` without
  needing the full `service` package, which depends on `messaging` - the same
  "leaf module" split the other ports already use for this, just enforced by
  the compiler here instead of being a stylistic choice.
- **Go has no virtual dispatch through struct embedding**, so `service.Service`
  isn't a Java/TS/C#-style abstract class with template methods. A `BaseService`
  struct supplies defaults that never need to call back into an override
  (`Pause`/`Unpause`/`Restart`/`HandleDocument`/`HandleEvent`/`HandleHeaders`),
  but `GracefulShutdown` and `HandleCommand` — which need to invoke whatever
  concrete `Start`/`Pause`/`Shutdown` the embedding type provides — are
  package-level functions (`DefaultGracefulShutdown(s Service)`,
  `DefaultHandleCommand(s Service, envelope)`) taking the full `Service`
  interface, which *does* dispatch dynamically. A concrete service calls them
  from its own one-line method body. See `service/service.go`.
- **`tasks.Runner` is goroutine + channel based**, not `std::thread`-based
  like `ra-common-cpp`'s `TaskRunner`. A single `stopCh` closed on `Shutdown()`
  broadcasts to the poll loop and every worker via `select`; a `sync.WaitGroup`
  tracks completion instead of joining each worker explicitly. This is a much
  smaller, lower-risk surface than the C++ version needed - see "What Go's
  race detector caught" below for the one bug it still found.
- **`route.RouteMeta.RouteID` is a plain `int64` JSON number**, not a
  string-encoded one. The string encoding in the JS-targeting ports works
  around JavaScript's float64-based JSON numbers losing precision above 2^53;
  Go's `encoding/json` round-trips `int64` exactly, so that workaround isn't
  needed here.
- **`[]byte` fields (`Multihash.Digest`, `Content.Body`) serialize as base64
  strings**, Go's built-in behavior for byte slices in `encoding/json` -
  unlike the other ports, which serialize a digest as a JSON array of byte
  values. Not wire-compatible with them (nothing in this family is), and more
  compact besides.

## Idiom mapping

| Java | Go |
|---|---|
| `JSONSerializable` | exported fields + `json:"..."` tags, `encoding/json.Marshal`/`Unmarshal` directly |
| abstract base + `Class.forName("type")` | an interface + a free-function `Unmarshal*` dispatcher on a `type`/`kind` tag |
| abstract base w/ shared fields (`BaseRoute`) | a `RouteMeta` struct held by each route |
| abstract class w/ behaviour (`BaseService`) | `BaseService` embedding for non-dispatching defaults + `Default*` free functions for anything needing dynamic dispatch (see above) |
| static utility class (`HashUtil`) | package-level functions (`crypto.GenerateHash`, ...) |
| `enum X { A("a") }` + `value()` | `type X string` + `const` values, or `.JcaName()`/`.Name()` methods where a second display form is needed |
| checked `*Exception` | `*raerror.RaError` implementing `error`, with a `Kind` field |
| `Stack<T>` / `DequeStack` | a `[]Route` slice, append/pop at the end (a Go slice's natural O(1) end is its "front" for this purpose - see `route/dynamic_routing_slip.go`) |
| `Properties` | `map[string]string` |
| `null` | a pointer (`*string`, `*int`) for an optional scalar, or a nil interface/pointer for an optional object |

## What Go's race detector caught

`go test -race ./tasks/...` found two real bugs before this port was done -
neither would have been visible without it:

1. **In the library**: `Runner.Shutdown()` closed `stopCh` and waited for
   workers, but never reset `stopCh` to `nil` afterward - so `Status()`'s
   `stopCh == nil` check (meant to detect "fully shut down") never fired, and
   `Status()` reported `Stopping` forever instead of `Shutdown`. Fixed by
   resetting `stopCh` (and `tasks`) to `nil` once `Shutdown()`'s wait returns.
2. **In the test, not the library**: the ported test counted task runs with a
   plain `int` closure variable (`runs++` from a worker goroutine, `runs < 1`
   polled from the test goroutine) - a real data race even though it "worked"
   in practice, because Go doesn't guarantee the poller ever observes the
   write without synchronization. `ra-common-cpp`'s C++ port and every other
   port's test does the same "poll a plain counter" pattern; only Go's
   toolchain has a race detector built in to actually catch it. Fixed with
   `atomic.Int32`.

`ra-common-cpp`'s `TaskRunner` needed three real fixes for dangling pointers,
unsynchronized fields, and destroying a joinable thread - see that port's
`DESIGN.md`. This port's `tasks.Runner` needed one (a status-tracking bug, not
a memory-safety one), and had it caught by a tool instead of by careful
reading. That gap is the point of using it, not a claim that Go code doesn't
need concurrency review.

## Bugs fixed during the port (from the Java original, inherited by every port)

Same fixes as the other `ra-common-*` ports:

- `Signature` (de)serialization was an empty stub in Java - implemented fully here.
- `BaseRoute.fromMap` read the key `"routedId"` instead of `"routeId"` - not
  applicable here in the same way (Go doesn't hand-write a `fromMap`), but the
  correct key (`route_id`) is what the struct tag uses.
- `Nonce` prune computed `max * (pct / 100)` -> 0 in integer division in Java -
  `maxSize * prunePercent / 100` here, with an O(1) map membership check plus
  a `container/list` for eviction order.
- `DID.getPassphraseHashAlgorithm()` could NPE - `EffectivePassphraseHashAlgorithm()`
  falls back to the stored field.

## Not here (Phase 2, deferred)

`currency/*`, `locale/*`, `Scrubber`, `RegExGen`, the full network service
layer, `Protocol` (multiaddr), `ShellCommand`, `BrowserUtil`, `FileUtil`,
`InfoVault*`, `social/*`. `DLC` is folded into `Envelope` methods (not ported
as a type).
