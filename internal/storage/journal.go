package storage

import (
	"os"
	"path/filepath"
	"sort"
	"sync"

	"example.com/releaseflow/internal/model"
)

type FileStore struct {
	root    string
	mu      sync.Mutex
	journal string
}

func NewFileStore(root string) (*FileStore, error) {
	store := &FileStore{root: root, journal: filepath.Join(root, "journal", "events.jsonl")}
	if err := store.Ensure(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *FileStore) Root() string { return s.root }

func (s *FileStore) Ensure() error {
	for _, dir := range []string{
		s.root,
		filepath.Dir(s.journal),
		filepath.Join(s.root, "snapshots"),
		filepath.Join(s.root, "reports"),
		filepath.Join(s.root, "index"),
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	if _, err := os.Stat(s.journal); os.IsNotExist(err) {
		f, err := os.OpenFile(s.journal, os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		return f.Close()
	}
	return nil
}

func (s *FileStore) AppendEvent(evt model.Event) error {
	if err := evt.Validate(); err != nil {
		return err
	}
	raw, err := encode(evt)
	if err != nil {
		return err
	}
	if err := s.addIndex(evt.PacketID); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.OpenFile(s.journal, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(raw); err != nil {
		return err
	}
	return nil
}

func (s *FileStore) Events(packetID string) ([]model.Event, error) {
	all, err := s.AllEvents()
	if err != nil {
		return nil, err
	}
	var events []model.Event
	for _, evt := range all {
		if evt.PacketID == packetID {
			events = append(events, evt)
		}
	}
	return events, nil
}

func (s *FileStore) ListPacketIDs() ([]string, error) {
	if ids, err := s.loadIndex(); err == nil && len(ids) > 0 {
		return ids, nil
	}
	all, err := s.AllEvents()
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	for _, evt := range all {
		seen[evt.PacketID] = struct{}{}
	}
	var ids []string
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	if err := s.saveIndex(ids); err != nil {
		return nil, err
	}
	return ids, nil
}
