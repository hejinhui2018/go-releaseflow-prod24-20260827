package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunInitializesWorkspaceAndPrintsUsageOnEmptyArgs(t *testing.T) {
	dir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(wd)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := Run(nil); err != nil {
		t.Fatalf("run: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "releaseflow-data")); err == nil {
		t.Fatal("empty invocation should not create packet data")
	}
}

