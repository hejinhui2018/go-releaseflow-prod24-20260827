package service

import (
	"fmt"

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
	if err := fs.RestoreBackup(path); err != nil {
		return err
	}
	fmt.Fprintln(s.out, "restore complete")
	return nil
}
