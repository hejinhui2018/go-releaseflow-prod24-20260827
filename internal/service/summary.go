package service

import (
	"encoding/json"
	"fmt"
	"time"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Summary(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	report := model.NewPacketReport(packet)
	view := struct {
		Report   model.PacketReport `json:"report"`
		Lifecycle model.Lifecycle    `json:"lifecycle"`
		CheckedAt time.Time          `json:"checked_at"`
	}{
		Report:   report,
		Lifecycle: model.DescribeLifecycle(packet, time.Now().UTC()),
		CheckedAt: time.Now().UTC(),
	}
	raw, err := json.MarshalIndent(view, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, string(raw))
	return nil
}

