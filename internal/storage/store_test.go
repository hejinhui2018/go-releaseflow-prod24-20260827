package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/releaseflow/internal/model"
)

func TestFileStoreRebuildsPacketFromJournal(t *testing.T) {
	root := t.TempDir()
	store, err := NewFileStore(root)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	now := time.Date(2026, time.February, 2, 11, 0, 0, 0, time.UTC)
	packet := model.Packet{ID: "release-010", Service: "ingest", Version: "v1.0.0", Owner: "li", Status: model.StateDraft, Revision: 1, CreatedAt: now, UpdatedAt: now}
	evt := model.Event{ID: "evt-10", PacketID: packet.ID, Kind: model.EventCreated, Revision: 1, At: now, Payload: []byte(`{"actor":"li","comment":"start"}`)}
	if err := model.ApplyEvent(&packet, evt); err != nil {
		t.Fatalf("apply event: %v", err)
	}
	if err := store.AppendEvent(evt); err != nil {
		t.Fatalf("append event: %v", err)
	}
	if err := store.SaveSnapshot(packet); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	recovered, err := store.Packet(packet.ID)
	if err != nil {
		t.Fatalf("recover packet: %v", err)
	}
	if recovered.Status != model.StateDraft {
		t.Fatalf("expected draft packet, got %s", recovered.Status)
	}
	reportPath, err := store.ExportSummary(recovered)
	if err != nil {
		t.Fatalf("export summary: %v", err)
	}
	if filepath.Base(reportPath) != "release-report.json" {
		t.Fatalf("unexpected report file %s", reportPath)
	}
	if _, err := os.Stat(reportPath); err != nil {
		t.Fatalf("expected report file to exist: %v", err)
	}
}

