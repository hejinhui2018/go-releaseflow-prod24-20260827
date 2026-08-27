package service

import (
	"fmt"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Publish(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	actor, _ := cmd.MustString("publisher")
	comment := cmd.String("comment")
	if err := s.policy.CanPublish(packet); err != nil {
		return err
	}
	packet, err = s.transition(packet, model.EventPublished, actor, comment)
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, packet.Status)
	return nil
}
