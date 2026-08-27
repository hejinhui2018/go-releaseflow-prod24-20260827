package service

import (
	"encoding/json"
	"fmt"
)

func (s *Service) History(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	raw, err := json.MarshalIndent(packet.History, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, string(raw))
	return nil
}

