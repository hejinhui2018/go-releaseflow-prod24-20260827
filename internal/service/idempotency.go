package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"example.com/releaseflow/internal/model"
)

func eventID(packetID string, rev int, kind model.EventKind, actor, comment string) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%s|%s|%s", packetID, rev, kind, actor, comment)))
	return hex.EncodeToString(sum[:])
}

func recentKind(packet model.Packet) model.EventKind {
	if last := packet.LastHistory(); last != nil {
		return model.EventKind(last.Kind)
	}
	return ""
}

