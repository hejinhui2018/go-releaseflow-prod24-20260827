package model

import "time"

type Packet struct {
	ID           string        `json:"id"`
	Service      string        `json:"service"`
	Version      string        `json:"version"`
	Owner        string        `json:"owner"`
	Environment  string        `json:"environment"`
	Risk         string        `json:"risk"`
	Summary      string        `json:"summary"`
	Status       State         `json:"status"`
	Revision     int           `json:"revision"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	SubmittedAt  *time.Time    `json:"submitted_at,omitempty"`
	ReviewedAt   *time.Time    `json:"reviewed_at,omitempty"`
	ApprovedAt   *time.Time    `json:"approved_at,omitempty"`
	PublishedAt  *time.Time    `json:"published_at,omitempty"`
	RolledBackAt *time.Time    `json:"rolled_back_at,omitempty"`
	Reviewer     string        `json:"reviewer,omitempty"`
	Approver     string        `json:"approver,omitempty"`
	Publisher    string        `json:"publisher,omitempty"`
	Rollbacker   string        `json:"rollbacker,omitempty"`
	ReviewNote   string        `json:"review_note,omitempty"`
	ApprovalNote string        `json:"approval_note,omitempty"`
	PublishNote  string        `json:"publish_note,omitempty"`
	RollbackNote string        `json:"rollback_note,omitempty"`
	Notes        []Note        `json:"notes,omitempty"`
	History      []HistoryEntry `json:"history,omitempty"`
	Checkpoint   Checkpoint    `json:"checkpoint"`
}

type Note struct {
	Author  string    `json:"author"`
	Message string    `json:"message"`
	At      time.Time `json:"at"`
}

type Review struct {
	Reviewer string    `json:"reviewer"`
	Comment  string    `json:"comment"`
	At       time.Time `json:"at"`
}

type Approval struct {
	Approver string    `json:"approver"`
	Comment  string    `json:"comment"`
	At       time.Time `json:"at"`
}

type Publication struct {
	Publisher string    `json:"publisher"`
	Comment   string    `json:"comment"`
	At        time.Time `json:"at"`
}

type Rollback struct {
	Rollbacker string    `json:"rollbacker"`
	Reason     string    `json:"reason"`
	At         time.Time `json:"at"`
}

