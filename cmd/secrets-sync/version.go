package main

import (
	"fmt"
	"runtime"
)

var (
	Version   = "0.1.0"
	GitCommit = "dev"
	BuildDate = "unknown"
)

func printVersion() {
	version := Version
	if GitCommit != "" && GitCommit != "dev" {
		version = fmt.Sprintf("%s-%s", Version, GitCommit[:7])
	} else if GitCommit == "dev" {
		version = "dev-" + getGitHash()
	}

	fmt.Printf("secrets-sync version %s\n", version)
	fmt.Printf("  Build date: %s\n", BuildDate)
	fmt.Printf("  Go version: %s\n", runtime.Version())
	fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func getGitHash() string {
	// This will be replaced by build script
	return "unknown"
}
