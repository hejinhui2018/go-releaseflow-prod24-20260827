package service

import (
	"fmt"

	"example.com/releaseflow/internal/model"
)

func (s *Service) Init(cmd Command) error {
	spec := packetSpec{
		ID:          cmd.String("id"),
		Service:     cmd.String("service"),
		Version:     cmd.String("version"),
		Owner:       cmd.String("owner"),
		Environment: cmd.String("environment"),
		Risk:        cmd.String("risk"),
		Summary:     cmd.String("summary"),
	}
	packet, err := s.createPacket(spec)
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, s.statusLine(packet))
	return nil
}

func (s *Service) Submit(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	actor, _ := cmd.MustString("actor")
	comment := cmd.String("comment")
	packet, err = s.transition(packet, model.EventSubmitted, actor, comment)
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, s.statusLine(packet))
	return nil
}

func (s *Service) Export(cmd Command) error {
	packet, err := s.packetFromCommand(cmd)
	if err != nil {
		return err
	}
	path, err := s.store.ExportSummary(packet)
	if err != nil {
		return err
	}
	fmt.Fprintln(s.out, path)
	return nil
}

func (s *Service) packetFromCommand(cmd Command) (model.Packet, error) {
	id := cmd.String("id")
	if id == "" {
		if len(cmd.Args) > 0 {
			id = cmd.Args[0]
		}
	}
	if id == "" {
		return model.Packet{}, fmt.Errorf("missing --id")
	}
	packet, err := s.store.Packet(id)
	if err != nil {
		return model.Packet{}, err
	}
	return packet, nil
}
