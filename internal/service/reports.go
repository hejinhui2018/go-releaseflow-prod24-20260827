package service

import (
	"encoding/json"
	"fmt"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Report(cmd Command) error {
	ids, err := s.store.ListPacketIDs()
	if err != nil {
		return err
	}
	report := model.ExportReport{}
	for _, id := range ids {
		packet, err := s.store.Packet(id)
		if err != nil {
			return err
		}
		report.Add(packet)
	}
	raw, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, string(raw))
	return nil
}

