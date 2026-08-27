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
	if fs, ok := s.store.(interface {
		PruneSnapshots([]string) error
		PruneJournal([]string) error
	}); ok {
		// PruneSnapshots removes snapshot files for packets that are no longer
		// kept, but the journal and index still reference those packets. If they
		// are left in place, Packet() rebuilds them from the journal and they
		// reappear in every listing, which is exactly what the handover team saw.
		// PruneJournal rewrites the journal so only kept packets survive, and the
		// index is rebuilt afterwards so it matches the surviving set.
		if err := fs.PruneJournal(keep); err != nil {
			return err
		}
		if err := fs.PruneSnapshots(keep); err != nil {
			return err
		}
	}
	if rb, ok := s.store.(interface{ RebuildIndex() error }); ok {
		if err := rb.RebuildIndex(); err != nil {
			return err
		}
	}
	fmt.Fprintln(s.out, "cleanup complete")
	return nil
}

