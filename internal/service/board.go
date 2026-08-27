package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"example.com/releaseflow/internal/model"
)

type boardRow struct {
	ID        string            `json:"id"`
	Service   string            `json:"service"`
	Status    string            `json:"status"`
	Next      string            `json:"next"`
	Reason    string            `json:"reason"`
	Terminal  bool              `json:"terminal"`
	Age       string            `json:"age"`
	Lifecycle model.Lifecycle   `json:"lifecycle"`
}

func (s *Service) Board(cmd Command) error {
	ids, err := s.store.ListPacketIDs()
	if err != nil {
		return err
	}
	rows := make([]boardRow, 0, len(ids))
	for _, id := range ids {
		packet, err := s.store.Packet(id)
		if err != nil {
			return err
		}
		lifecycle := model.DescribeLifecycle(packet, time.Now().UTC())
		rows = append(rows, boardRow{
			ID:        packet.ID,
			Service:   packet.Service,
			Status:    string(packet.Status),
			Next:      lifecycle.Next,
			Reason:    lifecycle.Reason,
			Terminal:  lifecycle.Terminal,
			Age:       lifecycle.Age.String(),
			Lifecycle: lifecycle,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Service == rows[j].Service {
			return rows[i].ID < rows[j].ID
		}
		return rows[i].Service < rows[j].Service
	})
	payload := struct {
		GeneratedAt time.Time  `json:"generated_at"`
		Rows        []boardRow `json:"rows"`
		Count       int        `json:"count"`
	}{
		GeneratedAt: time.Now().UTC(),
		Rows:        rows,
		Count:       len(rows),
	}
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, string(raw))
	return nil
}

func (s *Service) StatusLine(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, s.statusLine(packet))
	return nil
}

func (s *Service) NextAction(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, PlanNext(packet).Describe())
	return nil
}

