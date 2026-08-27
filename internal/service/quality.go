package service

import (
	"encoding/json"
	"fmt"
	"time"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Quality(cmd Command) error {
	ids, err := s.store.ListPacketIDs()
	if err != nil {
		return err
	}
	type quality struct {
		ID      string            `json:"id"`
		Issue   string            `json:"issue,omitempty"`
		Status  string            `json:"status"`
		Lifecycle model.Lifecycle `json:"lifecycle"`
	}
	result := struct {
		GeneratedAt time.Time `json:"generated_at"`
		Items       []quality `json:"items"`
	}{GeneratedAt: time.Now().UTC()}
	for _, id := range ids {
		packet, err := s.store.Packet(id)
		if err != nil {
			return err
		}
		check := model.CheckPacket(packet)
		item := quality{
			ID:        id,
			Status:    string(packet.Status),
			Lifecycle: model.DescribeLifecycle(packet, time.Now().UTC()),
		}
		if !check.OK {
			item.Issue = check.FirstIssue()
		}
		result.Items = append(result.Items, item)
	}
	raw, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, string(raw))
	return nil
}

