package model

import (
	"fmt"
	"strings"
	"time"
)

type PacketReport struct {
	Packet      Packet        `json:"packet"`
	Timeline    []HistoryEntry `json:"timeline"`
	GeneratedAt time.Time     `json:"generated_at"`
}

func NewPacketReport(packet Packet) PacketReport {
	timeline := make([]HistoryEntry, len(packet.History))
	copy(timeline, packet.History)
	return PacketReport{
		Packet:      packet,
		Timeline:    timeline,
		GeneratedAt: time.Now().UTC(),
	}
}

func (r PacketReport) Title() string {
	return fmt.Sprintf("%s %s %s", r.Packet.ID, r.Packet.Service, r.Packet.Status)
}

func (r PacketReport) Lines() []string {
	lines := []string{
		r.Title(),
		fmt.Sprintf("version=%s owner=%s", r.Packet.Version, r.Packet.Owner),
		fmt.Sprintf("revision=%d notes=%d", r.Packet.Revision, len(r.Packet.Notes)),
	}
	for _, note := range r.Packet.Notes {
		lines = append(lines, fmt.Sprintf("note %s %s", note.Author, note.Message))
	}
	return lines
}

func (r PacketReport) NeedsAttention() bool {
	return r.Packet.Status == StateDraft || r.Packet.Status == StateSubmitted
}

func (r PacketReport) StatusReason() string {
	switch r.Packet.Status {
	case StateDraft:
		return "waiting for submission"
	case StateSubmitted:
		return "waiting for review"
	case StateReviewed:
		return "waiting for approval"
	case StateApproved:
		return "waiting for publication"
	case StatePublished:
		return "published"
	case StateRolledBack:
		return "rolled back"
	default:
		return "unknown"
	}
}

func (r PacketReport) Render() string {
	return strings.Join(r.Lines(), "\n")
}

func (r PacketReport) Age(now time.Time) time.Duration {
	return now.Sub(r.Packet.UpdatedAt)
}

type StatusCounts map[State]int

func NewStatusCounts() StatusCounts {
	return StatusCounts{
		StateDraft:      0,
		StateSubmitted:  0,
		StateReviewed:   0,
		StateApproved:   0,
		StatePublished:  0,
		StateRolledBack: 0,
	}
}

func (c StatusCounts) Add(packet Packet) {
	c[packet.Status]++
}

func (c StatusCounts) Total() int {
	total := 0
	for _, v := range c {
		total += v
	}
	return total
}

func (c StatusCounts) MostActive() State {
	var best State
	var count int
	for state, n := range c {
		if n > count {
			best = state
			count = n
		}
	}
	return best
}

