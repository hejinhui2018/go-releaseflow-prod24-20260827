package service

import (
	"encoding/json"
	"fmt"
	"time"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Metrics(cmd Command) error {
	ids, err := s.store.ListPacketIDs()
	if err != nil {
		return err
	}
	counts := model.NewStatusCounts()
	ready := 0
	blocked := 0
	for _, id := range ids {
		packet, err := s.store.Packet(id)
		if err != nil {
			return err
		}
		counts.Add(packet)
		if model.IsReleaseReady(packet) {
			ready++
		}
		if model.IsBlocked(packet) {
			blocked++
		}
	}
	report := struct {
		GeneratedAt time.Time             `json:"generated_at"`
		Counts      model.StatusCounts    `json:"counts"`
		Ready       int                   `json:"ready"`
		Blocked     int                   `json:"blocked"`
		Total       int                   `json:"total"`
	}{
		GeneratedAt: time.Now().UTC(),
		Counts:      counts,
		Ready:       ready,
		Blocked:     blocked,
		Total:       counts.Total(),
	}
	raw, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, string(raw))
	return nil
}

