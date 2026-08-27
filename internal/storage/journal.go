package storage

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	f, err := os.Open(s.journal)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var events []model.Event
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var evt model.Event
		if err := decode([]byte(line), &evt); err != nil {
			return nil, err
		}
		if evt.PacketID == packetID {
			events = append(events, evt)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan journal: %w", err)
	}
	return events, nil
}

func (s *FileStore) ListPacketIDs() ([]string, error) {
	if ids, err := s.loadIndex(); err == nil && len(ids) > 0 {
		return ids, nil
	}
	f, err := os.Open(s.journal)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	seen := map[string]struct{}{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var evt model.Event
		if err := decode(scanner.Bytes(), &evt); err != nil {
			return nil, err
		}
		seen[evt.PacketID] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
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
