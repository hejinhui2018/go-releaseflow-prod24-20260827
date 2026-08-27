package model

import (
	"fmt"
	"time"
)

type Lifecycle struct {
	Status    State         `json:"status"`
	Next      string        `json:"next"`
	Terminal  bool          `json:"terminal"`
	Reason    string        `json:"reason"`
	Age       time.Duration `json:"age"`
	Attention bool          `json:"attention"`
}

func DescribeLifecycle(packet Packet, now time.Time) Lifecycle {
	plan := Lifecycle{
		Status:    packet.Status,
		Terminal:  packet.Status.IsTerminal(),
		Age:       now.Sub(packet.UpdatedAt),
		Attention: IsBlocked(packet),
	}
	switch packet.Status {
	case StateDraft:
		plan.Next = "submit"
		plan.Reason = "draft packets need a submit event"
	case StateSubmitted:
		plan.Next = "review"
		plan.Reason = "submitted packets wait for review"
	case StateReviewed:
		plan.Next = "approve"
		plan.Reason = "reviewed packets wait for approval"
	case StateApproved:
		plan.Next = "publish"
		plan.Reason = "approved packets can be published"
	case StatePublished:
		plan.Next = "rollback"
		plan.Reason = "published packets may still be rolled back"
	case StateRolledBack:
		plan.Reason = "rolled back packets are terminal"
	}
	return plan
}

func (l Lifecycle) String() string {
	if l.Reason == "" {
		return fmt.Sprintf("%s next=%s", l.Status, l.Next)
	}
	return fmt.Sprintf("%s next=%s reason=%s", l.Status, l.Next, l.Reason)
}

func (l Lifecycle) CanMove() bool {
	return !l.Terminal && l.Next != ""
}

func (l Lifecycle) NeedsAttention() bool {
	return l.Attention || l.Age > 72*time.Hour
}

func (l Lifecycle) IsStale() bool {
	return l.Age > 7*24*time.Hour
}

func (l Lifecycle) SummaryLine() string {
	return fmt.Sprintf("status=%s next=%s terminal=%t", l.Status, l.Next, l.Terminal)
}

