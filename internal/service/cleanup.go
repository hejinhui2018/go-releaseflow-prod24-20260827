package service

import (
	"fmt"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Clean(cmd Command) error {
	ids, err := s.store.ListPacketIDs()
	if err != nil {
		return err
	}
	var keep []string
	for _, id := range ids {
		packet, err := s.store.Packet(id)
		if err != nil {
			return err
		}
		if packet.Status == model.StatePublished || packet.Status == model.StateRolledBack {
			keep = append(keep, id)
		}
	}
	if fs, ok := s.store.(interface{ PruneSnapshots([]string) error }); ok {
		if err := fs.PruneSnapshots(keep); err != nil {
			return err
		}
	}
	fmt.Fprintln(s.out, "cleanup complete")
	return nil
}
