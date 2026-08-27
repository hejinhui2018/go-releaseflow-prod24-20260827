package service

import (
	"testing"

	"example.com/releaseflow/internal/storage"
)

func TestRestoreDropsPacketsCreatedAfterBackup(t *testing.T) {
	store, err := storage.NewFileStore(t.TempDir())
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	svc := New(store)

	exec := func(args ...string) {
		t.Helper()
		if err := svc.Execute(args); err != nil {
			t.Fatalf("execute %v: %v", args, err)
		}
	}

	exec("init", "--id=release-before-backup", "--service=delivery", "--version=v1.2.0", "--owner=mei", "--environment=prod", "--risk=low", "--summary=baseline")
	backup, err := store.CreateBackup()
	if err != nil {
		t.Fatalf("create backup: %v", err)
	}

	exec("init", "--id=release-after-backup", "--service=delivery", "--version=v1.3.0", "--owner=mei", "--environment=prod", "--risk=medium", "--summary=late packet")
	exec("submit", "--id=release-after-backup", "--actor=mei", "--comment=ready")
	exec("review", "--id=release-after-backup", "--reviewer=qin", "--comment=reviewed")
	exec("approve", "--id=release-after-backup", "--approver=ada", "--comment=approved")
	exec("publish", "--id=release-after-backup", "--publisher=ops", "--comment=published")

	exec("restore", "--path="+backup.Path)
	if packet, err := store.Packet("release-after-backup"); err == nil {
		t.Fatalf("packet created after backup survived restore: status=%s revision=%d", packet.Status, packet.Revision)
	}
	ids, err := store.ListPacketIDs()
	if err != nil {
		t.Fatalf("list restored packet ids: %v", err)
	}
	if len(ids) != 1 || ids[0] != "release-before-backup" {
		t.Fatalf("restored ids = %v, want only release-before-backup", ids)
	}
}
