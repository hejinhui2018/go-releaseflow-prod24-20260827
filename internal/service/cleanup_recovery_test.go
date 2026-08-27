package service

import (
	"bytes"
	"testing"

	"example.com/releaseflow/internal/storage"
)

func TestCleanupRemovesDrillPacketsAfterRestart(t *testing.T) {
	root := t.TempDir()
	store, err := storage.NewFileStore(root)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	svc := New(store)
	svc.out = &bytes.Buffer{}
	exec := func(args ...string) {
		t.Helper()
		if err := svc.Execute(args); err != nil {
			t.Fatalf("execute %v: %v", args, err)
		}
	}
	exec("init", "--id=drill-packet", "--service=payments", "--version=v1.0.0", "--owner=li", "--environment=drill", "--risk=low", "--summary= rehearsal")
	exec("init", "--id=release-packet", "--service=payments", "--version=v2.0.0", "--owner=mei", "--environment=prod", "--risk=low", "--summary=production")
	exec("submit", "--id=release-packet", "--actor=mei")
	exec("review", "--id=release-packet", "--reviewer=qi")
	exec("approve", "--id=release-packet", "--approver=ops")
	exec("publish", "--id=release-packet", "--publisher=ops")
	exec("clean")

	restarted, err := storage.NewFileStore(root)
	if err != nil {
		t.Fatalf("restart store: %v", err)
	}
	ids, err := restarted.ListPacketIDs()
	if err != nil {
		t.Fatalf("list after restart: %v", err)
	}
	if len(ids) != 1 || ids[0] != "release-packet" {
		t.Fatalf("expected only published packet after cleanup, got %v", ids)
	}
	if _, err := restarted.Packet("drill-packet"); err == nil {
		t.Fatal("drill packet reappeared after cleanup and restart")
	}
}
