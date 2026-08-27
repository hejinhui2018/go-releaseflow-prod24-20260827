package model

import "fmt"

type Transition struct {
	From  State `json:"from"`
	To    State `json:"to"`
	Label string `json:"label"`
}

var transitionTable = []Transition{
	{From: StateDraft, To: StateSubmitted, Label: "submit"},
	{From: StateSubmitted, To: StateReviewed, Label: "review"},
	{From: StateReviewed, To: StateApproved, Label: "approve"},
	{From: StateApproved, To: StatePublished, Label: "publish"},
	{From: StateApproved, To: StateRolledBack, Label: "rollback"},
	{From: StatePublished, To: StateRolledBack, Label: "rollback"},
}

func AllowedTransitions(from State) []State {
	var out []State
	for _, tr := range transitionTable {
		if tr.From == from {
			out = append(out, tr.To)
		}
	}
	return out
}

func TransitionLabel(from, to State) string {
	for _, tr := range transitionTable {
		if tr.From == from && tr.To == to {
			return tr.Label
		}
	}
	return ""
}

func CanReach(from, to State) bool {
	return TransitionLabel(from, to) != ""
}

func ExplainTransition(from, to State) string {
	if label := TransitionLabel(from, to); label != "" {
		return fmt.Sprintf("%s -> %s via %s", from, to, label)
	}
	return fmt.Sprintf("%s cannot reach %s", from, to)
}

func IsBlocked(packet Packet) bool {
	return packet.Status == StateDraft || packet.Status == StateSubmitted
}

func IsReleaseReady(packet Packet) bool {
	return packet.Status == StateApproved
}

