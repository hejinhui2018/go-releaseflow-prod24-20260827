package model

import "testing"

func TestManifestDetectsDrift(t *testing.T) {
	packet := Packet{ID: "release-002", Service: "billing", Version: "v1.2.3", Owner: "mia", Status: StateDraft}
	manifest, err := NewManifest(packet)
	if err != nil {
		t.Fatalf("new manifest: %v", err)
	}
	packet.Version = "v1.2.4"
	if err := manifest.Validate(packet); err == nil {
		t.Fatal("expected manifest drift to be detected")
	}
}

