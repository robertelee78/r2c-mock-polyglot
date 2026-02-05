// Package native bridges the gateway to its sibling first-party helper units:
// the C `exiftrim` metadata stripper and the Rust `worker` image processor.
// Both are spawned as child processes.
package native

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// helperDir is where the compiled sibling-unit binaries are installed.
// Overridable for local runs where the units build into their own trees.
func helperDir() string {
	if d := os.Getenv("SNAPSTORE_HELPERS"); d != "" {
		return d
	}
	return "/usr/local/libexec/snapstore"
}

// RunExiftrim invokes the C exiftrim unit on an artifact and returns the
// number of bytes it reclaimed.
func RunExiftrim(artifactPath string) (string, error) {
	bin := filepath.Join(helperDir(), "exiftrim")
	out, err := exec.Command(bin, artifactPath).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// JobSpec is the argument bundle passed to the Rust worker unit.
type JobSpec struct {
	Input   string
	Output  string
	Filter  string
	Quality string
}

// RunWorker forwards a job to the Rust worker unit and returns its stdout.
func RunWorker(spec JobSpec) (string, error) {
	bin := filepath.Join(helperDir(), "worker")
	out, err := exec.Command(bin,
		"--input", spec.Input,
		"--output", spec.Output,
		"--filter", spec.Filter,
		"--quality", spec.Quality,
	).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
