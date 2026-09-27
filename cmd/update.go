package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/HimanshuSardana/kite/internal/version"
)

const (
	updateOwner = "HimanshuSardana"
	updateRepo  = "kite"
)

const updateHelp = `Update kite to the latest GitHub release.

USAGE:
  kite update [--check] [--force]

FLAGS:
  --check       Only check for updates, don't install anything
  --force, -f   Reinstall even if already on the latest version
  -h, --help    Show this help message

DESCRIPTION:
  Downloads the prebuilt binary for your platform from the latest
  GitHub release and swaps it over the current executable.
  Set GITHUB_TOKEN to raise the API rate limit.
`

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type githubRelease struct {
	TagName string         `json:"tag_name"`
	Assets  []releaseAsset `json:"assets"`
}

var updateClient = &http.Client{Timeout: 30 * time.Second}

func runUpdate(args []string) {
	checkOnly := false
	force := false
	for _, a := range args[2:] {
		switch a {
		case "--check":
			checkOnly = true
		case "--force", "-f":
			force = true
		case "-h", "--help":
			fmt.Print(updateHelp)
			return
		default:
			fmt.Fprintf(os.Stderr, "Unknown flag %q\n\n%s", a, updateHelp)
			os.Exit(1)
		}
	}

	fmt.Println("Checking for updates...")
	rel, err := latestRelease()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not check for updates: %v\n", err)
		os.Exit(1)
	}

	cur, lat := normalizeVersion(version.Version), normalizeVersion(rel.TagName)
	fmt.Printf("Current: %s   Latest: %s\n", displayVersion(version.Version), displayVersion(rel.TagName))

	if cur != "" && cur == lat && !force {
		fmt.Println("Already up to date.")
		return
	}
	if checkOnly {
		fmt.Printf("Update available: run `kite update` to install %s.\n", rel.TagName)
		return
	}

	asset := pickAsset(rel)
	if asset == nil {
		fmt.Fprintf(os.Stderr, "Error: release %s has no binary for %s/%s\n", rel.TagName, runtime.GOOS, runtime.GOARCH)
		os.Exit(1)
	}
	fmt.Printf("Downloading %s...\n", asset.Name)

	if err := installAsset(rel, asset); err != nil {
		fmt.Fprintf(os.Stderr, "Error: update failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Updated to %s. Run `kite version` to verify.\n", rel.TagName)
}

// normalizeVersion strips a leading "v" so "v0.1.0" and "0.1.0" compare
// equal. "dev"/empty (a local build) normalizes to "" — always updatable.
func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimPrefix(v, "V")
	if v == "" || v == "dev" {
		return ""
	}
	return v
}

func displayVersion(v string) string {
	if normalizeVersion(v) == "" {
		return "dev (local build)"
	}
	return v
}

func latestRelease() (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", updateOwner, updateRepo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "kite-updater")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := updateClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no releases published yet (releases are cut from v* tags)")
	}
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("GitHub API rate limit exceeded (set GITHUB_TOKEN to raise the limit)")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("GitHub API returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("decoding release metadata: %w", err)
	}
	if rel.TagName == "" {
		return nil, fmt.Errorf("release has no tag name")
	}
	return &rel, nil
}

// pickAsset finds the binary for the current platform. The release workflow
// publishes kite-<os>-<arch>.tar.gz (unix) and kite-<os>-<arch>.exe (windows).
func pickAsset(rel *githubRelease) *releaseAsset {
	base := fmt.Sprintf("kite-%s-%s", runtime.GOOS, runtime.GOARCH)
	for _, ext := range []string{".tar.gz", ".zip", ".exe", ""} {
		for i := range rel.Assets {
			if rel.Assets[i].Name == base+ext {
				return &rel.Assets[i]
			}
		}
	}
	return nil
}

func downloadFile(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "kite-updater")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := updateClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned %s", resp.Status)
	}

	tmp, err := os.CreateTemp("", "kite-update-*")
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			os.Remove(tmp.Name())
		}
	}()

	if _, err = io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return "", err
	}
	if err = tmp.Close(); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

// verifyChecksum checks the downloaded archive against checksums.txt when the
// release publishes one; it warns (but proceeds) when it doesn't.
func verifyChecksum(rel *githubRelease, archivePath, assetName string) error {
	var sums *releaseAsset
	for i := range rel.Assets {
		n := strings.ToLower(rel.Assets[i].Name)
		if n == "checksums.txt" || n == "sha256sums.txt" || n == "sha256sums" {
			sums = &rel.Assets[i]
			break
		}
	}
	if sums == nil {
		fmt.Println("Warning: release has no checksums file, skipping verification.")
		return nil
	}

	sumsPath, err := downloadFile(sums.BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("downloading checksums: %w", err)
	}
	defer os.Remove(sumsPath)

	raw, err := os.ReadFile(sumsPath)
	if err != nil {
		return fmt.Errorf("reading checksums: %w", err)
	}
	want := ""
	for _, line := range strings.Split(string(raw), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := strings.TrimPrefix(fields[1], "*")
		if name == assetName || strings.HasSuffix(name, "/"+assetName) {
			want = fields[0]
			break
		}
	}
	if want == "" {
		return fmt.Errorf("no checksum entry for %s", assetName)
	}

	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	if got := hex.EncodeToString(h.Sum(nil)); !strings.EqualFold(got, want) {
		return fmt.Errorf("checksum mismatch for %s", assetName)
	}
	fmt.Println("Checksum verified.")
	return nil
}

// extractBinary returns a temp file holding the kite executable from the
// downloaded asset (an archive or, for bare .exe assets, the file itself).
func extractBinary(archivePath, assetName string) (string, error) {
	binName := "kite"
	if runtime.GOOS == "windows" {
		binName = "kite.exe"
	}

	switch {
	case strings.HasSuffix(assetName, ".tar.gz"):
		return extractTarGz(archivePath, binName)
	case strings.HasSuffix(assetName, ".zip"):
		return extractZip(archivePath, binName)
	default:
		return archivePath, nil // bare binary asset
	}
}

func writeTempExecutable(r io.Reader) (string, error) {
	tmp, err := os.CreateTemp("", "kite-bin-*")
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			os.Remove(tmp.Name())
		}
	}()
	if _, err = io.Copy(tmp, r); err != nil {
		tmp.Close()
		return "", err
	}
	if err = tmp.Close(); err != nil {
		return "", err
	}
	if err = os.Chmod(tmp.Name(), 0o755); err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

func extractTarGz(archivePath, binName string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", fmt.Errorf("not a gzip archive: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		name := filepath.Base(hdr.Name)
		if hdr.Typeflag == tar.TypeReg && (name == binName || name == "kite" || name == "kite.exe") {
			return writeTempExecutable(tr)
		}
	}
	return "", fmt.Errorf("archive contains no %s binary", binName)
}

func extractZip(archivePath, binName string) (string, error) {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer zr.Close()

	for _, f := range zr.File {
		name := filepath.Base(f.Name)
		if name != binName && name != "kite" && name != "kite.exe" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		path, werr := writeTempExecutable(rc)
		rc.Close()
		if werr != nil {
			return "", werr
		}
		return path, nil
	}
	return "", fmt.Errorf("archive contains no %s binary", binName)
}

func installAsset(rel *githubRelease, asset *releaseAsset) error {
	archivePath, err := downloadFile(asset.BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("downloading %s: %w", asset.Name, err)
	}
	defer os.Remove(archivePath)

	if err := verifyChecksum(rel, archivePath, asset.Name); err != nil {
		return err
	}

	binPath, err := extractBinary(archivePath, asset.Name)
	if err != nil {
		return err
	}
	if binPath != archivePath {
		defer os.Remove(binPath)
	}

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locating current executable: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}

	// Rename-over is atomic on unix and works while the old binary runs.
	if err := os.Rename(binPath, exe); err != nil {
		if runtime.GOOS == "windows" {
			nextTo := exe + ".new"
			if err := copyExecutable(binPath, nextTo); err != nil {
				return err
			}
			return fmt.Errorf("Windows cannot replace a running executable: new binary saved to %s — close kite and rename it over %s", nextTo, exe)
		}
		return fmt.Errorf("replacing %s: %w (hint: you may need write permission on that directory)", exe, err)
	}
	return nil
}

func copyExecutable(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
