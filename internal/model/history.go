package model

import "time"

type HistoryEntry struct {
	Kind      string    `json:"kind"`
	Actor     string    `json:"actor"`
	Comment   string    `json:"comment"`
	Revision  int       `json:"revision"`
	RecordedAt time.Time `json:"recorded_at"`
}

func NewHistoryEntry(kind, actor, comment string, revision int, at time.Time) HistoryEntry {
	return HistoryEntry{Kind: kind, Actor: actor, Comment: comment, Revision: revision, RecordedAt: at}
}

func (p *Packet) appendHistory(entry HistoryEntry) {
	p.History = append(p.History, entry)
}

func (p *Packet) LastHistory() *HistoryEntry {
	if len(p.History) == 0 {
		return nil
	}
	return &p.History[len(p.History)-1]
}

