package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Backup struct {
	CreatedAt time.Time `json:"created_at"`
	Path      string    `json:"path"`
}

func (s *FileStore) CreateBackup() (Backup, error) {
	now := time.Now().UTC()
	dir := filepath.Join(s.layout().BackupDir(), now.Format("20060102T150405Z"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Backup{}, err
	}
	for _, src := range []string{s.journal, s.indexPath()} {
		if err := copyFile(src, filepath.Join(dir, filepath.Base(src))); err != nil {
			return Backup{}, err
		}
	}
	if err := copySnapshotDir(s.layout().SnapshotDir(), filepath.Join(dir, "snapshots")); err != nil {
		return Backup{}, err
	}
	return Backup{CreatedAt: now, Path: dir}, nil
}

func (s *FileStore) RestoreBackup(path string) error {
	currentIDs, err := s.loadIndex()
	if err != nil {
		return err
	}
	if err := copyFile(filepath.Join(path, filepath.Base(s.journal)), s.journal); err != nil {
		return err
	}
	if err := copyFile(filepath.Join(path, filepath.Base(s.indexPath())), s.indexPath()); err != nil {
		return err
	}
	restoredIDs, err := s.loadIndex()
	if err != nil {
		return err
	}
	seen := make(map[string]struct{}, len(restoredIDs))
	for _, id := range restoredIDs {
		seen[id] = struct{}{}
	}
	for _, id := range currentIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		restoredIDs = append(restoredIDs, id)
	}
	if err := s.saveIndex(restoredIDs); err != nil {
		return err
	}
	return copySnapshotDir(filepath.Join(path, "snapshots"), s.layout().SnapshotDir())
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func copySnapshotDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if err := copyFile(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())); err != nil {
			return fmt.Errorf("copy %s: %w", entry.Name(), err)
		}
	}
	return nil
}
