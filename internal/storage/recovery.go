package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/releaseflow/internal/model"
)

func (s *FileStore) Packet(packetID string) (model.Packet, error) {
	if packet, err := s.LoadSnapshot(packetID); err == nil {
		if packet.Status != "" {
			return packet, nil
		}
	}
	events, err := s.Events(packetID)
	if err != nil {
		return model.Packet{}, err
	}
	if len(events) == 0 {
		return model.Packet{}, model.ErrNotFound
	}
	var packet model.Packet
	for _, evt := range events {
		if err := model.ApplyEvent(&packet, evt); err != nil {
			return model.Packet{}, err
		}
	}
	return packet, nil
}

func (s *FileStore) Rebuild(packetID string) (model.Packet, error) {
	packet, err := s.Packet(packetID)
	if err != nil {
		return model.Packet{}, err
	}
	if err := s.SaveSnapshot(packet); err != nil {
		return model.Packet{}, err
	}
	return packet, nil
}

func (s *FileStore) RebuildAll() ([]model.Packet, error) {
	ids, err := s.ListPacketIDs()
	if err != nil {
		return nil, err
	}
	var packets []model.Packet
	for _, id := range ids {
		packet, err := s.Rebuild(id)
		if err != nil {
			return nil, fmt.Errorf("rebuild %s: %w", id, err)
		}
		packets = append(packets, packet)
	}
	return packets, nil
}

func (s *FileStore) SnapshotExists(packetID string) bool {
	_, err := os.Stat(filepath.Join(s.root, "snapshots", packetID+".json"))
	return err == nil
}

