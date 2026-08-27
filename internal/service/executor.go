package service

import (
	"encoding/json"
	"fmt"

	"example.com/releaseflow/internal/model"
)

type packetSpec struct {
	ID          string
	Service     string
	Version     string
	Owner       string
	Environment string
	Risk        string
	Summary     string
}

func (s *Service) createPacket(spec packetSpec) (model.Packet, error) {
	id, err := model.NormalizeID(spec.ID)
	if err != nil {
		return model.Packet{}, err
	}
	if err := s.policy.ValidatePacket(model.Packet{Service: spec.Service, Version: spec.Version, Owner: spec.Owner}); err != nil {
		return model.Packet{}, err
	}
	now := s.clock.Now()
	packet := model.Packet{
		ID:          id,
		Service:     spec.Service,
		Version:     spec.Version,
		Owner:       spec.Owner,
		Environment: spec.Environment,
		Risk:        spec.Risk,
		Summary:     spec.Summary,
		Status:      model.StateDraft,
		Revision:    1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	evt := model.Event{
		ID:       eventID(id, 1, model.EventCreated, spec.Owner, spec.Summary),
		PacketID: id,
		Kind:     model.EventCreated,
		Revision: 1,
		At:       now,
		Payload:  mustPayload(map[string]string{"actor": spec.Owner, "comment": spec.Summary}),
	}
	if err := model.ApplyEvent(&packet, evt); err != nil {
		return model.Packet{}, err
	}
	if err := s.store.AppendEvent(evt); err != nil {
		return model.Packet{}, err
	}
	if err := s.store.SaveSnapshot(packet); err != nil {
		return model.Packet{}, err
	}
	return packet, nil
}

func (s *Service) transition(packet model.Packet, kind model.EventKind, actor, comment string) (model.Packet, error) {
	var next model.State
	switch kind {
	case model.EventSubmitted:
		next = model.StateSubmitted
	case model.EventReviewed:
		next = model.StateReviewed
	case model.EventApproved:
		next = model.StateApproved
	case model.EventPublished:
		next = model.StatePublished
	case model.EventRolledBack:
		next = model.StateRolledBack
	case model.EventNoted:
		next = packet.Status
	default:
		return model.Packet{}, fmt.Errorf("unsupported transition %s", kind)
	}
	if err := model.ValidateTransition(packet.Status, next); err != nil && kind != model.EventRolledBack {
		return model.Packet{}, err
	}
	now := s.clock.Now()
	packet.Revision++
	evt := model.Event{
		ID:       eventID(packet.ID, packet.Revision, kind, actor, comment),
		PacketID: packet.ID,
		Kind:     kind,
		Revision: packet.Revision,
		At:       now,
		Payload:  mustPayload(map[string]string{"actor": actor, "comment": comment}),
	}
	if err := model.ApplyEvent(&packet, evt); err != nil {
		return model.Packet{}, err
	}
	if err := s.store.SaveSnapshot(packet); err != nil {
		return model.Packet{}, err
	}
	if err := s.store.AppendEvent(evt); err != nil {
		return model.Packet{}, err
	}
	return packet, nil
}

func mustPayload(v any) []byte {
	raw, _ := json.Marshal(v)
	return raw
}
