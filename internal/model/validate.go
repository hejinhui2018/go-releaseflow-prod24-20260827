package model

import "fmt"

func ValidatePacket(packet Packet) error {
	if packet.ID == "" {
		return fmt.Errorf("%w: missing packet id", ErrNotReady)
	}
	if packet.Service == "" {
		return fmt.Errorf("%w: missing service", ErrNotReady)
	}
	if packet.Version == "" {
		return fmt.Errorf("%w: missing version", ErrNotReady)
	}
	if packet.Status == "" {
		return fmt.Errorf("%w: missing status", ErrNotReady)
	}
	return nil
}

func ValidateTransition(from, to State) error {
	if !from.CanAdvanceTo(to) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, from, to)
	}
	return nil
}

func ValidateNote(note Note) error {
	if note.Author == "" || note.Message == "" {
		return fmt.Errorf("%w: incomplete note", ErrNotReady)
	}
	return nil
}

