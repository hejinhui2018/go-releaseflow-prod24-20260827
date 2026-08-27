package model

import "errors"

var (
	ErrNotFound          = errors.New("packet not found")
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrDuplicateEvent    = errors.New("duplicate event")
	ErrNotReady          = errors.New("packet is not ready")
	ErrCorruptPacket     = errors.New("packet is corrupt")
)

