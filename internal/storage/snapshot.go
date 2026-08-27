package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/releaseflow/internal/model"
)

func (s *FileStore) snapshotPath(packetID string) string {
	return filepath.Join(s.root, "snapshots", packetID+".json")
}

func (s *FileStore) SaveSnapshot(packet model.Packet) error {
	raw, err := encode(packet)
	if err != nil {
		return err
	}
	return os.WriteFile(s.snapshotPath(packet.ID), raw, 0o644)
}

func (s *FileStore) LoadSnapshot(packetID string) (model.Packet, error) {
	raw, err := os.ReadFile(s.snapshotPath(packetID))
	if err != nil {
		return model.Packet{}, err
	}
	var packet model.Packet
	if err := decode(raw, &packet); err != nil {
		return model.Packet{}, err
	}
	if packet.ID != packetID {
		return model.Packet{}, fmt.Errorf("snapshot packet mismatch")
	}
	return packet, nil
}

