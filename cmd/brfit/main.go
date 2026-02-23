// Package main is the entry point for the brfit CLI.
package main

// Build-time variables injected via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Pass build info to root command
	SetBuildInfo(version, commit, date)
	Execute()
}
