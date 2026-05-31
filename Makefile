BINARY     := bin/litepub
COVERDATA  := coverdata
COVOUT     := coverage.out
COVOUT_INT := coverage-integration.out

.PHONY: all build test integration check coverage clean

all: check integration

build:
	go build -o $(BINARY) .

test:
	go test -race -coverprofile=$(COVOUT) ./...

integration: build
	mkdir -p $(COVERDATA)
	go build -cover -o $(BINARY) .
	cd tests && GOCOVERDIR=$(CURDIR)/$(COVERDATA) pytest -v
	go tool covdata textfmt -i=$(COVERDATA) -o=$(COVOUT_INT)

check:
	go vet ./...
	staticcheck ./...
	go test -race ./...

coverage: test integration
	go tool cover -func=$(COVOUT)
	go tool cover -func=$(COVOUT_INT)

clean:
	rm -f $(BINARY) $(COVOUT) $(COVOUT_INT)
	rm -rf $(COVERDATA)
