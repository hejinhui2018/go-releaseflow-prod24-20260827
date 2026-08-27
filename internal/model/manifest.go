package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Manifest struct {
	PacketID  string    `json:"packet_id"`
	Digest    string    `json:"digest"`
	Revision  int       `json:"revision"`
	Status    State     `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewManifest(packet Packet) (Manifest, error) {
	digest, err := PacketDigest(packet)
	if err != nil {
		return Manifest{}, err
	}
	return Manifest{
		PacketID:  packet.ID,
		Digest:    digest,
		Revision:  packet.Revision,
		Status:    packet.Status,
		UpdatedAt: packet.UpdatedAt,
	}, nil
}

func (m Manifest) Validate(packet Packet) error {
	actual, err := PacketDigest(packet)
	if err != nil {
		return err
	}
	if actual != m.Digest {
		return fmt.Errorf("%w: manifest digest mismatch", ErrCorruptPacket)
	}
	if packet.ID != m.PacketID {
		return fmt.Errorf("%w: manifest packet mismatch", ErrCorruptPacket)
	}
	return nil
}

func (m Manifest) Marshal() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

