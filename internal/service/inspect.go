package service

import (
	"encoding/json"
	"fmt"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Inspect(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	report := struct {
		Summary model.Summary `json:"summary"`
		Plan    Plan          `json:"plan"`
	}{
		Summary: model.NewSummary(packet),
		Plan:    PlanNext(packet),
	}
	raw, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, string(raw))
	return nil
}

