package model

import "time"

type Summary struct {
	PacketID     string     `json:"packet_id"`
	Service      string     `json:"service"`
	Version      string     `json:"version"`
	Status       State      `json:"status"`
	Revision     int        `json:"revision"`
	HistoryCount int        `json:"history_count"`
	NoteCount    int        `json:"note_count"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Checkpoint   Checkpoint `json:"checkpoint"`
}

func NewSummary(packet Packet) Summary {
	return Summary{
		PacketID:     packet.ID,
		Service:      packet.Service,
		Version:      packet.Version,
		Status:       packet.Status,
		Revision:     packet.Revision,
		HistoryCount: len(packet.History),
		NoteCount:    len(packet.Notes),
		UpdatedAt:    packet.UpdatedAt,
		Checkpoint:   packet.Checkpoint,
	}
}

type ExportReport struct {
	GeneratedAt time.Time `json:"generated_at"`
	Items       []Summary `json:"items"`
}

func (r *ExportReport) Add(packet Packet) {
	r.Items = append(r.Items, NewSummary(packet))
}

