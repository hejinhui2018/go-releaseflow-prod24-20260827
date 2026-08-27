package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"example.com/releaseflow/internal/model"
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

// PruneJournal rewrites the event journal so that only events for packets in
// keepIDs remain. Packets that are no longer kept have no surviving snapshot,
// so their journal events must be removed as well; otherwise Packet() would
// rebuild them from the journal and they would reappear in every listing.
func (s *FileStore) PruneJournal(keepIDs []string) error {
	keep := map[string]struct{}{}
	for _, id := range keepIDs {
		keep[id] = struct{}{}
	}
	events, err := s.AllEvents()
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.OpenFile(s.journal, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, evt := range events {
		if _, ok := keep[evt.PacketID]; !ok {
			continue
		}
		raw, err := encode(evt)
		if err != nil {
			return err
		}
		if _, err := f.Write(raw); err != nil {
			return err
		}
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}

func (s *FileStore) ResetJournal() error {
	return os.WriteFile(s.journal, nil, 0o644)
}

// AllEvents reads every event from the journal. The journal stores events as
// pretty-printed JSON objects that can span multiple lines, so the stream is
// decoded object-by-object rather than line-by-line.
func (s *FileStore) AllEvents() ([]model.Event, error) {
	f, err := os.Open(s.journal)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var events []model.Event
	dec := json.NewDecoder(f)
	for {
		var evt model.Event
		if err := dec.Decode(&evt); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("decode journal: %w", err)
		}
		events = append(events, evt)
	}
	return events, nil
}

