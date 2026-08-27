package model

import (
	"testing"
	"time"
)

func TestApplyEventBuildsCheckpointAndHistory(t *testing.T) {
	now := time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC)
	packet := Packet{ID: "release-001", Status: StateDraft}
	evt := Event{
		ID:       "evt-1",
		PacketID: "release-001",
		Kind:     EventSubmitted,
		Revision: 2,
		At:       now,
		Payload:  []byte(`{"actor":"alice","comment":"ready"}`),
	}
	if err := ApplyEvent(&packet, evt); err != nil {
		t.Fatalf("apply event: %v", err)
	}
	if packet.Status != StateSubmitted {
		t.Fatalf("expected submitted status, got %s", packet.Status)
	}
	if packet.Checkpoint.Revision != 2 || packet.Checkpoint.Digest == "" {
		t.Fatalf("expected checkpoint to be recorded, got %+v", packet.Checkpoint)
	}
	if len(packet.History) != 1 {
		t.Fatalf("expected one history item, got %d", len(packet.History))
	}
}

