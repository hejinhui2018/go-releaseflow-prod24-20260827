package storage

import (
	"os"
	"path/filepath"
	"time"

	"example.com/releaseflow/internal/model"
)

func (s *FileStore) Export(report model.ExportReport) (string, error) {
	if report.GeneratedAt.IsZero() {
		report.GeneratedAt = time.Now().UTC()
	}
	raw, err := encode(report)
	if err != nil {
		return "", err
	}
	path := filepath.Join(s.root, "reports", "release-report.json")
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func (s *FileStore) ExportSummary(packet model.Packet) (string, error) {
	report := model.ExportReport{GeneratedAt: time.Now().UTC()}
	report.Add(packet)
	return s.Export(report)
}

func (s *FileStore) ReportPath() string {
	return filepath.Join(s.root, "reports", "release-report.json")
}

func (s *FileStore) PacketReport(packetID string) (string, error) {
	packet, err := s.Packet(packetID)
	if err != nil {
		return "", err
	}
	return s.ExportSummary(packet)
}
