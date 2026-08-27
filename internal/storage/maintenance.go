package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

func (s *FileStore) PruneSnapshots(keepIDs []string) error {
	keep := map[string]struct{}{}
	for _, id := range keepIDs {
		keep[id] = struct{}{}
	}
	entries, err := os.ReadDir(filepath.Join(s.root, "snapshots"))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		id := entry.Name()
		if filepath.Ext(id) != ".json" {
			continue
		}
		base := id[:len(id)-len(filepath.Ext(id))]
		if _, ok := keep[base]; ok {
			continue
		}
		if err := os.Remove(filepath.Join(s.root, "snapshots", id)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove snapshot %s: %w", id, err)
		}
	}
	return nil
}

func (s *FileStore) ResetJournal() error {
	return os.WriteFile(s.journal, nil, 0o644)
}

