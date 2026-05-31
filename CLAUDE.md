# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

**litepub** — a lightweight TCP-based pub/sub message broker in Go, inspired by NATS. Clients connect over TCP, publish messages to topics, and subscribe to receive them. The broker handles routing and fanout.

Chosen depth feature: **at-least-once delivery with client ACKs and redelivery on timeout** (connected subscribers only — no durable session state across disconnects).

## Commands

```bash
go build ./...          # build all packages
go test ./...           # run all tests
go test ./... -run TestName   # run a single test by name
go test -race ./...     # run tests with race detector
go vet ./...            # static analysis
```

Linting requires `staticcheck`:
```bash
go install honnef.co/go/tools/cmd/staticcheck@2025.1
staticcheck ./...
```

Integration tests (requires binary to be built first):
```bash
go build -o bin/litepub .
cd tests && pytest -v
```

## Planned Architecture

### Wire Protocol
Text-based, line-oriented (like NATS). Commands are terminated by `\r\n`. Expected commands:

- Client → Broker: `SUB <topic> <sid>`, `PUB <topic> <bytes>\r\n<payload>`, `UNSUB <sid>`, `ACK <msg-id>`
- Broker → Client: `MSG <topic> <sid> <msg-id> <bytes>\r\n<payload>`, `+OK`, `-ERR <reason>`

### Concurrency Model
- One goroutine per TCP connection (read loop)
- Shared subscription registry protected by a `sync.RWMutex`
- Fan-out: broker iterates subscribers and sends to each via a per-subscriber channel; slow/blocked consumers are dropped or bounded

### At-Least-Once Delivery
- Each published message gets a unique ID (monotonic counter or similar)
- Broker tracks in-flight messages per subscriber (unACKed)
- A redelivery timer loop retransmits unACKed messages after a configurable timeout
- Scope: connected subscribers only — disconnect before ACK drops the message (explicitly a non-goal to persist across reconnects)

### Package Layout (intended)
```
cmd/litepub/        # broker binary entry point
internal/broker/    # subscription registry, fanout, in-flight tracking
internal/proto/     # protocol parser and serializer
internal/server/    # TCP listener, connection handling
client/             # Go client library (pub/sub API)
```

## Git Hooks

Commit hooks live in `.githooks/`. Activate them after cloning:

```bash
git config core.hooksPath .githooks
```

## Issue Tracking

GitHub Issues is used to track upcoming work. Check open issues at the start of a session for context on what's planned — they inform direction but are not hard requirements:

```bash
gh issue list
```

## Progress

### Completed
- `hello.go` — working TCP echo server
- `internal/parser/` — full protocol parser for `SUB`, `PUB`, `UNSUB`, `ACK`; 1 MB payload cap; unit tests for all branches
- `main.go` — TCP server with `+OK`/`-ERR` responses after each command; write error handling; clean EOF disconnect
- `signals.go` — SIGTERM/SIGINT handler for coverage flush; see `TODO(#19)` for planned replacement via graceful shutdown
- `tests/` — Python integration tests covering all four commands; runs against a live binary in CI

### Where We Left Off
- Next: broker fanout — subscription registry + routing `PubCommand` to matching subscribers

### Key Decisions Made
- Text-based protocol chosen deliberately: debuggable with `nc`, loggable; tradeoff is slightly more bytes on the wire vs binary
- `sid` is client-assigned and scoped per-subscription — lets a single client route incoming `MSG` frames to the right local handler
- `signals.go` split from `main.go` deliberately to keep signal/coverage concerns out of the main connection loop; delete when issue #19 is resolved
- Client library (`client/`) deferred until broker fanout is working

## Coverage (Integration Tests)

`go build -cover` registers counter-flush as an `os.Exit` hook. SIGTERM and SIGINT both bypass that hook — the process dies without writing counter data, leaving `main.go` at 0% even though the binary was exercised. The fix is a signal handler in `main.go` that catches both signals, calls `runtime/coverage.WriteCountersDir(os.Getenv("GOCOVERDIR"))`, then calls `os.Exit(0)` (which runs the hook and writes the meta file too).

`WriteCountersDir` is a documented no-op when the binary was not built with `-cover`, so the handler is safe to leave in production builds.

The CI workflow sets `GOCOVERDIR` before running pytest; the subprocess inherits it automatically. `conftest.py` also passes `env=os.environ.copy()` explicitly for clarity, and calls `proc.wait()` after `proc.terminate()` to ensure the flush completes before `go tool covdata` runs.
