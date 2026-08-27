package model

type State string

const (
	StateDraft      State = "draft"
	StateSubmitted  State = "submitted"
	StateReviewed   State = "reviewed"
	StateApproved   State = "approved"
	StatePublished  State = "published"
	StateRolledBack State = "rolled_back"
)

func (s State) String() string { return string(s) }

func (s State) IsTerminal() bool {
	return s == StatePublished || s == StateRolledBack
}

func (s State) CanAdvanceTo(next State) bool {
	switch s {
	case StateDraft:
		return next == StateSubmitted
	case StateSubmitted:
		return next == StateReviewed
	case StateReviewed:
		return next == StateApproved
	case StateApproved:
		return next == StatePublished || next == StateRolledBack
	case StatePublished:
		return next == StateRolledBack
	default:
		return false
	}
}

