package service

import (
	"fmt"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Review(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	actor, _ := cmd.MustString("reviewer")
	comment := cmd.String("comment")
	packet, err = s.transition(packet, model.EventReviewed, actor, comment)
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, PlanNext(packet).Describe())
	return nil
}

func (s *Service) Approve(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	actor, _ := cmd.MustString("approver")
	comment := cmd.String("comment")
	packet, err = s.transition(packet, model.EventApproved, actor, comment)
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, PlanNext(packet).Describe())
	return nil
}
