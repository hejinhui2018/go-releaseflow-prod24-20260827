package service

import (
	"bytes"
	"testing"

	"example.com/releaseflow/internal/storage"
)

func TestServicePublishesApprovedPacket(t *testing.T) {
	root := t.TempDir()
	store, err := storage.NewFileStore(root)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	svc := New(store)
	buf := &bytes.Buffer{}
	svc.out = buf

	mustExec := func(args ...string) {
		if err := svc.Execute(args); err != nil {
			t.Fatalf("execute %v: %v", args, err)
		}
	}

	mustExec("init", "--id=release-100", "--service=orders", "--version=v2.1.0", "--owner=lei", "--environment=prod", "--risk=medium", "--summary=nightly cut")
	mustExec("submit", "--id=release-100", "--actor=lei", "--comment=submitted for review")
	mustExec("review", "--id=release-100", "--reviewer=qin", "--comment=looks ready")
	mustExec("approve", "--id=release-100", "--approver=ada", "--comment=approved for release")
	mustExec("publish", "--id=release-100", "--publisher=ops", "--comment=release posted")

	packet, err := store.Packet("release-100")
	if err != nil {
		t.Fatalf("packet: %v", err)
	}
	if packet.Status != "published" {
		t.Fatalf("expected published packet, got %s", packet.Status)
	}
	if buf.Len() == 0 {
		t.Fatal("expected command output for operators")
	}
}

