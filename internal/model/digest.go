package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func PacketDigest(packet Packet) (string, error) {
	sum, err := digestValue(map[string]any{
		"id":          packet.ID,
		"service":     packet.Service,
		"version":     packet.Version,
		"owner":       packet.Owner,
		"environment": packet.Environment,
		"risk":        packet.Risk,
		"status":      packet.Status,
		"revision":    packet.Revision,
		"notes":       packet.Notes,
		"history":     packet.History,
	})
	if err != nil {
		return "", err
	}
	return sum, nil
}

func EventDigest(evt Event) (string, error) {
	return digestValue(evt)
}

func digestValue(v any) (string, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("digest json: %w", err)
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:]), nil
}

