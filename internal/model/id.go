package model

import (
	"errors"
	"fmt"
	"strings"
)

const MinPacketIDLen = 6

func NormalizeID(raw string) (string, error) {
	id := strings.TrimSpace(strings.ToLower(raw))
	if len(id) < MinPacketIDLen {
		return "", fmt.Errorf("packet id too short")
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return "", errors.New("packet id contains invalid characters")
		}
	}
	return id, nil
}

func MustNormalizeID(raw string) string {
	id, err := NormalizeID(raw)
	if err != nil {
		panic(err)
	}
	return id
}

