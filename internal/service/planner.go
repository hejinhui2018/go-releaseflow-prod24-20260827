package service

import (
	"fmt"

	"example.com/releaseflow/internal/model"
)

type Plan struct {
	PacketID string
	Next     string
	Blocker  string
}

func PlanNext(packet model.Packet) Plan {
	plan := Plan{PacketID: packet.ID}
	switch packet.Status {
	case model.StateDraft:
		plan.Next = "submit"
	case model.StateSubmitted:
		plan.Next = "review"
	case model.StateReviewed:
		plan.Next = "approve"
	case model.StateApproved:
		plan.Next = "publish"
	case model.StatePublished:
		plan.Next = "rollback"
	default:
		plan.Blocker = "packet is terminal"
	}
	return plan
}

func (p Plan) Describe() string {
	if p.Blocker != "" {
		return fmt.Sprintf("%s: %s", p.PacketID, p.Blocker)
	}
	return fmt.Sprintf("%s: next=%s", p.PacketID, p.Next)
}

