package service

import (
	"fmt"
	"time"
)

func (s *Service) Recover(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	rebuilt, err := s.store.Rebuild(packet.ID)
	if err != nil {
		return err
	}
	fmt.Fprintf(s.out, "%s recovered at %s\n", rebuilt.ID, time.Now().UTC().Format(time.RFC3339))
	return nil
}

