package service

import (
	"fmt"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Rollback(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	actor, _ := cmd.MustString("rollbacker")
	comment := cmd.String("reason")
	if packet.Status != model.StatePublished && packet.Status != model.StateApproved {
		return fmt.Errorf("rollback requires a published or approved packet")
	}
	packet, err = s.transition(packet, model.EventRolledBack, actor, comment)
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, packet.Status)
	return nil
}
