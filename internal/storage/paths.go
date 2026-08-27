package storage

import "path/filepath"

type Layout struct {
	Root string
}

func (l Layout) Journal() string {
	return filepath.Join(l.Root, "journal", "events.jsonl")
}

func (l Layout) SnapshotDir() string {
	return filepath.Join(l.Root, "snapshots")
}

func (l Layout) ReportDir() string {
	return filepath.Join(l.Root, "reports")
}

func (l Layout) IndexFile() string {
	return filepath.Join(l.Root, "index", "packets.json")
}

func (l Layout) BackupDir() string {
	return filepath.Join(l.Root, "backups")
}

func (l Layout) Snapshot(packetID string) string {
	return filepath.Join(l.SnapshotDir(), packetID+".json")
}

