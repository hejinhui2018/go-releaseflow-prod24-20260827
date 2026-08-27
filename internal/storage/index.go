package storage

import (
	"os"
	"sort"
)

type packetIndex struct {
	IDs []string `json:"ids"`
}

func (s *FileStore) layout() Layout {
	return Layout{Root: s.root}
}

func (s *FileStore) indexPath() string {
	return s.layout().IndexFile()
}

func (s *FileStore) loadIndex() ([]string, error) {
	raw, err := os.ReadFile(s.indexPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var idx packetIndex
	if err := decode(raw, &idx); err != nil {
		return nil, err
	}
	return idx.IDs, nil
}

func (s *FileStore) saveIndex(ids []string) error {
	sorted := append([]string(nil), ids...)
	sort.Strings(sorted)
	raw, err := encode(packetIndex{IDs: sorted})
	if err != nil {
		return err
	}
	return os.WriteFile(s.indexPath(), raw, 0o644)
}

func (s *FileStore) addIndex(packetID string) error {
	ids, err := s.loadIndex()
	if err != nil {
		return err
	}
	for _, id := range ids {
		if id == packetID {
			return nil
		}
	}
	ids = append(ids, packetID)
	return s.saveIndex(ids)
}

func (s *FileStore) rebuildIndexFromJournal() error {
	ids, err := s.ListPacketIDs()
	if err != nil {
		return err
	}
	return s.saveIndex(ids)
}

func (s *FileStore) removeIndex(packetID string) error {
	ids, err := s.loadIndex()
	if err != nil {
		return err
	}
	filtered := ids[:0]
	for _, id := range ids {
		if id != packetID {
			filtered = append(filtered, id)
		}
	}
	if filtered == nil {
		filtered = []string{}
	}
	return s.saveIndex(filtered)
}

func (s *FileStore) IndexContains(packetID string) bool {
	ids, err := s.loadIndex()
	if err != nil {
		return false
	}
	for _, id := range ids {
		if id == packetID {
			return true
		}
	}
	return false
}

func (s *FileStore) IndexedIDs() ([]string, error) {
	return s.loadIndex()
}
