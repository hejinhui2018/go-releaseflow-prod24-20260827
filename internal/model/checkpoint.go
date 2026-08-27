package model

import "time"

type Checkpoint struct {
	Revision   int       `json:"revision"`
	Digest     string    `json:"digest"`
	RecordedAt time.Time `json:"recorded_at"`
}

func (c *Checkpoint) Set(revision int, digest string, at time.Time) {
	c.Revision = revision
	c.Digest = digest
	c.RecordedAt = at
}

func (c Checkpoint) IsZero() bool {
	return c.Revision == 0 && c.Digest == ""
}

