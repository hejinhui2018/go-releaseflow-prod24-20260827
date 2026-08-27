package model

import "fmt"

type Policy struct {
	RequireReview bool `json:"require_review"`
	RequireOwner   bool `json:"require_owner"`
}

func DefaultPolicy() Policy {
	return Policy{RequireReview: true, RequireOwner: true}
}

func (p Policy) ValidatePacket(packet Packet) error {
	if packet.Service == "" || packet.Version == "" {
		return fmt.Errorf("%w: missing service or version", ErrNotReady)
	}
	if p.RequireOwner && packet.Owner == "" {
		return fmt.Errorf("%w: missing owner", ErrNotReady)
	}
	return nil
}

func (p Policy) CanPublish(packet Packet) error {
	if packet.Status != StateApproved {
		return fmt.Errorf("%w: packet not approved", ErrNotReady)
	}
	return nil
}

