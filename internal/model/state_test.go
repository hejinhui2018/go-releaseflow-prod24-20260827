package model

import "testing"

func TestStateAdvancesInReleaseOrder(t *testing.T) {
	if !StateDraft.CanAdvanceTo(StateSubmitted) {
		t.Fatal("draft packets should move into submission")
	}
	if !StateSubmitted.CanAdvanceTo(StateReviewed) {
		t.Fatal("submitted packets should move into review")
	}
	if !StateReviewed.CanAdvanceTo(StateApproved) {
		t.Fatal("reviewed packets should move into approval")
	}
	if !StateApproved.CanAdvanceTo(StatePublished) {
		t.Fatal("approved packets should publish")
	}
	if !StatePublished.IsTerminal() || !StateRolledBack.IsTerminal() {
		t.Fatal("published and rolled back packets should be terminal")
	}
}

