package service

import (
	"encoding/json"
	"fmt"
)

func (s *Service) List(cmd Command) error {
	ids, err := s.store.ListPacketIDs()
	if err != nil {
		return err
	}
	type item struct {
		ID      string `json:"id"`
		Status  string `json:"status"`
		Service string `json:"service"`
		Version string `json:"version"`
	}
	var items []item
	for _, id := range ids {
		packet, err := s.store.Packet(id)
		if err != nil {
			return err
		}
		items = append(items, item{ID: packet.ID, Status: string(packet.Status), Service: packet.Service, Version: packet.Version})
	}
	raw, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, string(raw))
	return nil
}

