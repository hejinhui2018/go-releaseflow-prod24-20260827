package service

import (
	"fmt"

	"example.com/releaseflow/internal/model"
	"example.com/releaseflow/internal/storage"
)

func (s *Service) Backup(cmd Command) error {
	fs, ok := s.store.(interface{ CreateBackup() (storage.Backup, error) })
	if !ok {
		return fmt.Errorf("store does not support backups")
	}
	backup, err := fs.CreateBackup()
	if err != nil {
		return err
	}
	fmt.Fprintf(s.out, "%v\n", backup)
	return nil
}

func (s *Service) Restore(cmd Command) error {
	path := cmd.String("path")
	if path == "" {
		return fmt.Errorf("missing --path")
	}
	fs, ok := s.store.(interface{ RestoreBackup(string) error })
	if !ok {
		return fmt.Errorf("store does not support restore")
	}
	ids, err := s.store.ListPacketIDs()
	if err != nil {
		return err
	}
	packets := make([]model.Packet, 0, len(ids))
	for _, id := range ids {
		packet, err := s.store.Packet(id)
		if err != nil {
			return err
		}
		packets = append(packets, packet)
	}
	if err := fs.RestoreBackup(path); err != nil {
		return err
	}
	for _, packet := range packets {
		if err := s.store.SaveSnapshot(packet); err != nil {
			return err
		}
	}
	fmt.Fprintln(s.out, "restore complete")
	return nil
}
