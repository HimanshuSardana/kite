package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	cases := []struct{ in, want string }{
		{"v0.1.0", "0.1.0"},
		{"V0.1.0", "0.1.0"},
		{"0.1.0", "0.1.0"},
		{"  v1.2.3  ", "1.2.3"},
		{"dev", ""},
		{"", ""},
		{"v", ""},
	}
	for _, c := range cases {
		if got := normalizeVersion(c.in); got != c.want {
			t.Errorf("normalizeVersion(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPickAsset(t *testing.T) {
	base := fmt.Sprintf("kite-%s-%s", runtime.GOOS, runtime.GOARCH)
	rel := &githubRelease{TagName: "v0.1.0", Assets: []releaseAsset{
		{Name: "checksums.txt"},
		{Name: base + ".tar.gz"},
		{Name: base + ".zip"},
	}}
	if got := pickAsset(rel); got == nil || got.Name != base+".tar.gz" {
		t.Fatalf("expected %s.tar.gz to win, got %+v", base, got)
	}

	// Bare binary fallback when no archive exists.
	bare := &githubRelease{Assets: []releaseAsset{{Name: base}}}
	if got := pickAsset(bare); got == nil || got.Name != base {
		t.Fatalf("expected bare %s, got %+v", base, got)
	}

	// Windows-style .exe asset.
	exeRel := &githubRelease{Assets: []releaseAsset{{Name: base + ".exe"}}}
	_ = exeRel // picked only when GOOS/GOARCH match; just ensure no panic
	pickAsset(exeRel)

	if got := pickAsset(&githubRelease{}); got != nil {
		t.Fatalf("expected nil for empty release, got %+v", got)
	}
}

func TestFindChecksumsAsset(t *testing.T) {
	for _, name := range []string{"checksums.txt", "SHA256SUMS.TXT", "sha256sums", "sha256sums.txt"} {
		rel := &githubRelease{Assets: []releaseAsset{{Name: name}}}
		if findChecksumsAsset(rel) == nil {
			t.Errorf("expected to find %q", name)
		}
	}
	if findChecksumsAsset(&githubRelease{Assets: []releaseAsset{{Name: "kite-linux-amd64.tar.gz"}}}) != nil {
		t.Error("expected nil when no checksums file present")
	}
}

func TestLookupChecksum(t *testing.T) {
	body := "abc123  kite-linux-amd64.tar.gz\ndef456 *kite-darwin-arm64.tar.gz\n999999  dist/kite-windows-amd64.exe\n"
	if got, ok := lookupChecksum(body, "kite-linux-amd64.tar.gz"); !ok || got != "abc123" {
		t.Errorf("plain entry: got %q,%v", got, ok)
	}
	if got, ok := lookupChecksum(body, "kite-darwin-arm64.tar.gz"); !ok || got != "def456" {
		t.Errorf("star entry: got %q,%v", got, ok)
	}
	if got, ok := lookupChecksum(body, "kite-windows-amd64.exe"); !ok || got != "999999" {
		t.Errorf("prefixed entry: got %q,%v", got, ok)
	}
	if _, ok := lookupChecksum(body, "kite-missing.tar.gz"); ok {
		t.Error("expected no match for missing asset")
	}
}

func TestVerifyFileChecksum(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bin")
	content := []byte("fake-binary-bytes")
	if err := os.WriteFile(p, content, 0o600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(content)
	want := hex.EncodeToString(sum[:])
	if err := verifyFileChecksum(p, want, "kite-test"); err != nil {
		t.Fatalf("expected match, got %v", err)
	}
	if err := verifyFileChecksum(p, strings.ToUpper(want), "kite-test"); err != nil {
		t.Fatalf("expected case-insensitive match, got %v", err)
	}
	if err := verifyFileChecksum(p, "00"+want[2:], "kite-test"); err == nil {
		t.Fatal("expected mismatch error")
	}
}

func TestEnsureExecutable(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "kite-bin")
	if err := os.WriteFile(p, []byte("#!/bin/sh\necho hi\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ensureExecutable(p); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm()&0o111 == 0 {
		t.Fatalf("expected exec bits, got %o", fi.Mode().Perm())
	}
}

func TestExtractBinaryBareSetsExecBit(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "kite-download")
	if err := os.WriteFile(p, []byte("bare-binary"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := extractBinary(p, "kite-linux-amd64")
	if err != nil {
		t.Fatal(err)
	}
	if got != p {
		t.Fatalf("bare path should return input, got %q", got)
	}
	fi, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm()&0o111 == 0 {
		t.Fatalf("bare binary missing exec bit: %o", fi.Mode().Perm())
	}
}

func TestVerifyChecksumMissingFailsClosed(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "asset")
	if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	rel := &githubRelease{TagName: "v0.1.0", Assets: []releaseAsset{{Name: "kite-linux-amd64.tar.gz"}}}
	if err := verifyChecksum(rel, p, "kite-linux-amd64.tar.gz", false); err == nil {
		t.Fatal("expected error when checksums file missing and insecure=false")
	}
	if err := verifyChecksum(rel, p, "kite-linux-amd64.tar.gz", true); err != nil {
		t.Fatalf("expected warn-and-continue with insecure=true, got %v", err)
	}
}

func TestUpdateHelpMentionsInsecure(t *testing.T) {
	if !strings.Contains(updateHelp, "--insecure") {
		t.Error("updateHelp should document --insecure")
	}
}
