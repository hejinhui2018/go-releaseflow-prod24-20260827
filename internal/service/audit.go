package service

import (
	"fmt"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Audit(packet model.Packet) error {
	if err := model.ValidatePacket(packet); err != nil {
		return err
	}
	if packet.Checkpoint.IsZero() {
		return fmt.Errorf("%w: missing checkpoint", model.ErrNotReady)
	}
	digest, err := model.PacketDigest(packet)
	if err != nil {
		return err
	}
	if digest != packet.Checkpoint.Digest {
		return fmt.Errorf("%w: checkpoint mismatch", model.ErrCorruptPacket)
	}
	return nil
}

func (s *Service) statusLine(packet model.Packet) string {
	return fmt.Sprintf("%s %s %s rev=%d", packet.ID, packet.Service, packet.Status, packet.Revision)
}

