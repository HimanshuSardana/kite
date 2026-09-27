// Package version holds the kite release version.
//
// The defaults below describe a local build. Release builds inject real
// values via -ldflags (see the Makefile):
//
//	go build -ldflags "-X github.com/HimanshuSardana/kite/internal/version.Version=v0.1.0 \
//	  -X github.com/HimanshuSardana/kite/internal/version.Commit=abc1234 \
//	  -X github.com/HimanshuSardana/kite/internal/version.Date=2026-09-27T00:00:00Z" .
package version

// Version is the release tag (e.g. "v0.1.0") or "dev" for local builds.
var Version = "dev"

// Commit is the short git hash the binary was built from.
var Commit = "none"

// Date is the UTC build timestamp.
var Date = "unknown"

// Known reports whether the binary carries a real release version
// (as opposed to a local "dev" build).
func Known() bool {
	return Version != "" && Version != "dev"
}

// String renders the full version line, e.g.
// "v0.1.0 (commit abc1234, built 2026-09-27T00:00:00Z)".
func String() string {
	s := Version
	var extra []string
	if Commit != "" && Commit != "none" {
		extra = append(extra, "commit "+Commit)
	}
	if Date != "" && Date != "unknown" {
		extra = append(extra, "built "+Date)
	}
	if len(extra) > 0 {
		s += " (" + extra[0]
		for _, e := range extra[1:] {
			s += ", " + e
		}
		s += ")"
	}
	return s
}
