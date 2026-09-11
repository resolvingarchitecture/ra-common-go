# TODO

## P0 — done

- [x] Core packages ported: raerror, encoding (base32/58; base64 is stdlib),
      util (bytes/strings/version/random/UniqueID/Nonce), crypto (Hash/
      HashUtil/Multihash/HashCash/EncryptionAlgorithm, on stdlib SHA + HMAC),
      identity (Signature/PublicKey/Did), network, servicestatus, route
      (Route/SimpleRoute/SimpleExternalRoute/RelayedExternalRoute/
      DynamicRoutingSlip), lifecycle, multipart, messaging (Message hierarchy
      + Envelope + channel/bus interfaces, one package - see DESIGN.md),
      service (ServiceCore/Service), tasks (goroutine-based Runner), content, config.
- [x] Test suite ported from `ra-common-ts`'s `core.test.ts` + `envelope.test.ts`,
      one `_test.go` per package, all passing under `go test ./...` **and**
      `go test -race ./...`.
- [x] One real bug found and fixed via the race detector (`Runner.Status()`
      never reporting `Shutdown`) - see DESIGN.md.

## P1 — hardening

- [ ] `Runner.Shutdown()` has the same narrow non-linearizable-against-a-
      concurrent-`AddTask` gap noted in `ra-common-cpp`'s TODO - revisit if a
      real service host ever exercises that race in practice.
- [ ] Fuzz-test `route.UnmarshalRoute`/`messaging.UnmarshalMessage` against
      malformed/truncated JSON (right now a bad `type`/`kind` tag returns a
      clean error, but deeply malformed nested JSON hasn't been fuzzed).
- [ ] Decide whether `Content.Body`/`Multihash.Digest`'s base64-string wire
      shape (Go's default `[]byte` JSON encoding, vs. the byte-array shape
      the other ports use) matters for any future cross-port tooling - it
      doesn't today since none of these ports are wire-compatible with each
      other anyway.

## Phase 2 (deferred, mirrors the other ra-common-* ports)

- [ ] `currency/*` (~72 Java files) → enum + table
- [ ] `locale/*` + `LocaleUtil` / `LanguageUtil` / `Resources`
- [ ] `Scrubber`, `RegExGen` / `IntegerRangeRegex`
- [ ] full `network` service layer
- [ ] `Protocol` (multiaddr)
- [ ] `ShellCommand`, `BrowserUtil`, `FileUtil`
- [ ] `InfoVault*`, `SimpleByteCache`, `OrderedProperties`
- [ ] `social/*`

## Packaging

- [ ] Tag a `v0.1.0` release once the API settles, so consumers can pin via
      `go get github.com/resolvingarchitecture/ra-common-go@v0.1.0` instead
      of a commit hash.
