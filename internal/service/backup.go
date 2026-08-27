package service

import (
	"fmt"

	"example.com/releaseflow/internal/storage"
)

func (s *Service) Backup(cmd Command) error {
	fs, ok := s.store.(interface {
		CreateBackup() (storage.Backup, error)
	})
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
		return fmt.Errorf("store does not restore")
	}
	if err := fs.RestoreBackup(path); err != nil {
		return err
	}
	// Rebuild snapshots from the restored journal so in-memory state matches the
	// backup point. Rebuilding must happen after RestoreBackup reverts the journal,
	// index, and snapshots; saving pre-restore packets here would overwrite the
	// restored state and resurrect packets created after the backup.
	ids, err := s.store.ListPacketIDs()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if _, err := s.store.Rebuild(id); err != nil {
			return err
		}
	}
	fmt.Fprintln(s.out, "restore complete")
	return nil
}
