package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/coverage"
	"syscall"
)

// handleSignals catches SIGTERM/SIGINT and flushes coverage data before
// exiting. go build -cover registers the flush as an os.Exit hook, but signal
// delivery bypasses that hook, so we call WriteCountersDir explicitly here.
// WriteCountersDir is a no-op in non-instrumented builds.
//
// TODO(#19): delete this once main() blocks on a shutdown channel instead of
// looping forever — at that point main returns naturally and coverage flushes
// via the os.Exit hook without any explicit call here.
func handleSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	<-sigs
	if dir := os.Getenv("GOCOVERDIR"); dir != "" {
		if err := coverage.WriteCountersDir(dir); err != nil {
			log.Printf("coverage flush error: %v", err)
		}
	}
	os.Exit(0)
}
