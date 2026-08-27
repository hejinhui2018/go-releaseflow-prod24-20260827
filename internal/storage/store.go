package storage

import (
	"example.com/releaseflow/internal/model"
)

type Store interface {
	Ensure() error
	AppendEvent(model.Event) error
	Events(string) ([]model.Event, error)
	Packet(string) (model.Packet, error)
	SaveSnapshot(model.Packet) error
	LoadSnapshot(string) (model.Packet, error)
	Rebuild(string) (model.Packet, error)
	RebuildAll() ([]model.Packet, error)
	ListPacketIDs() ([]string, error)
	Export(model.ExportReport) (string, error)
	ExportSummary(model.Packet) (string, error)
	Root() string
}
