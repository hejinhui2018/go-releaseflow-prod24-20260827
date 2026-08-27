package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type EventKind string

const (
	EventCreated   EventKind = "created"
	EventSubmitted EventKind = "submitted"
	EventReviewed  EventKind = "reviewed"
	EventApproved  EventKind = "approved"
	EventPublished EventKind = "published"
	EventRolledBack EventKind = "rolled_back"
	EventNoted     EventKind = "noted"
)

type Event struct {
	ID       string          `json:"id"`
	PacketID string          `json:"packet_id"`
	Kind     EventKind       `json:"kind"`
	Revision int             `json:"revision"`
	At       time.Time       `json:"at"`
	Payload  json.RawMessage `json:"payload,omitempty"`
}

func (e Event) Validate() error {
	if e.ID == "" || e.PacketID == "" || e.Kind == "" {
		return fmt.Errorf("%w: missing event identity", ErrNotReady)
	}
	return nil
}

func (e Event) Actor() string {
	var wrapper struct {
		Actor string `json:"actor"`
	}
	_ = json.Unmarshal(e.Payload, &wrapper)
	return wrapper.Actor
}

func ApplyEvent(packet *Packet, evt Event) error {
	if packet == nil {
		return fmt.Errorf("%w: nil packet", ErrNotReady)
	}
	if packet.ID != evt.PacketID {
		return fmt.Errorf("%w: event packet mismatch", ErrCorruptPacket)
	}
	switch evt.Kind {
	case EventCreated:
		return applyCreated(packet, evt)
	case EventSubmitted:
		return applySubmitted(packet, evt)
	case EventReviewed:
		return applyReviewed(packet, evt)
	case EventApproved:
		return applyApproved(packet, evt)
	case EventPublished:
		return applyPublished(packet, evt)
	case EventRolledBack:
		return applyRolledBack(packet, evt)
	case EventNoted:
		return applyNoted(packet, evt)
	default:
		return fmt.Errorf("%w: unknown event kind %s", ErrNotReady, evt.Kind)
	}
}

func applyCreated(packet *Packet, evt Event) error {
	return applyCommon(packet, evt, StateDraft)
}

func applySubmitted(packet *Packet, evt Event) error {
	if err := ValidateTransition(packet.Status, StateSubmitted); err != nil {
		return err
	}
	return applyCommon(packet, evt, StateSubmitted)
}

func applyReviewed(packet *Packet, evt Event) error {
	if err := ValidateTransition(packet.Status, StateReviewed); err != nil {
		return err
	}
	return applyCommon(packet, evt, StateReviewed)
}

func applyApproved(packet *Packet, evt Event) error {
	if err := ValidateTransition(packet.Status, StateApproved); err != nil {
		return err
	}
	return applyCommon(packet, evt, StateApproved)
}

func applyPublished(packet *Packet, evt Event) error {
	if err := ValidateTransition(packet.Status, StatePublished); err != nil {
		return err
	}
	return applyCommon(packet, evt, StatePublished)
}

func applyRolledBack(packet *Packet, evt Event) error {
	if err := ValidateTransition(packet.Status, StateRolledBack); err != nil && packet.Status != StatePublished {
		return err
	}
	return applyCommon(packet, evt, StateRolledBack)
}

func applyNoted(packet *Packet, evt Event) error {
	return applyCommon(packet, evt, packet.Status)
}

func applyCommon(packet *Packet, evt Event, next State) error {
	var payload struct {
		Actor   string `json:"actor"`
		Comment string `json:"comment"`
	}
	_ = json.Unmarshal(evt.Payload, &payload)
	packet.Status = next
	packet.Revision = evt.Revision
	packet.UpdatedAt = evt.At
	packet.appendHistory(NewHistoryEntry(string(evt.Kind), payload.Actor, payload.Comment, evt.Revision, evt.At))
	switch evt.Kind {
	case EventCreated:
		packet.CreatedAt = evt.At
	case EventSubmitted:
		packet.SubmittedAt = &evt.At
	case EventReviewed:
		packet.ReviewedAt = &evt.At
		packet.Reviewer = payload.Actor
		packet.ReviewNote = payload.Comment
	case EventApproved:
		packet.ApprovedAt = &evt.At
		packet.Approver = payload.Actor
		packet.ApprovalNote = payload.Comment
	case EventPublished:
		packet.PublishedAt = &evt.At
		packet.Publisher = payload.Actor
		packet.PublishNote = payload.Comment
	case EventRolledBack:
		packet.RolledBackAt = &evt.At
		packet.Rollbacker = payload.Actor
		packet.RollbackNote = payload.Comment
	case EventNoted:
		packet.Notes = append(packet.Notes, Note{Author: payload.Actor, Message: payload.Comment, At: evt.At})
	}
	digest, err := PacketDigest(*packet)
	if err != nil {
		return err
	}
	packet.Checkpoint.Set(packet.Revision, digest, evt.At)
	return nil
}

