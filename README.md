# Litepub

[![codecov](https://codecov.io/gh/billvanleeuwen424/litepub/branch/main/graph/badge.svg?token=M7MXPECVAU)](https://codecov.io/gh/billvanleeuwen424/litepub)

A lightweight TCP pub/sub message broker in Go, inspired by NATS.

## Contributing

After cloning, activate the git hooks:

```bash
git config core.hooksPath .githooks
```

This enables the `prepare-commit-msg` hook which prompts for testing notes on every commit.

### Prerequisites

Install `staticcheck`:

```bash
go install honnef.co/go/tools/cmd/staticcheck@2025.1
```

Ensure `$(go env GOPATH)/bin` is on your `PATH`, then install Python dependencies for integration tests:

```bash
pip install -r tests/requirements.txt
```

### Commands

```bash
make build        # build the binary to bin/litepub
make test         # unit tests with race detector
make integration  # build and run integration tests (pytest)
make check        # vet + lint + unit tests
make coverage     # run all tests and print coverage summary
make all          # check + integration
make clean        # remove build artifacts
```