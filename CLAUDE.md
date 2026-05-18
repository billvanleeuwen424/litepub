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

Linting requires `golangci-lint`:
```bash
golangci-lint run       # lint all packages
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

## Progress

### Completed
- `hello.go` — working TCP echo server; covers net.Listen, Accept loop, goroutine-per-connection, io.Copy, conn.Close. Tested with `nc localhost 8080`.

### Where We Left Off
- Next step: create `internal/proto/` and write the protocol parser
- First command to parse: `SUB <topic> <sid>\r\n`
- Test strategy: use raw `nc` to send commands manually — no client library needed yet
- Client library (`client/`) is deferred until the broker is working

### Key Decisions Made
- Text-based protocol chosen deliberately: debuggable with `nc`, loggable, tradeof is slightly more bytes on the wire vs binary
- `sid` is client-assigned and scoped per-subscription (not a global client ID) — lets a single client route incoming `MSG` frames to the right local handler
- Client library would own the sid→handler mapping and the read loop; application code just registers callbacks
