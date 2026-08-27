package model

import (
	"fmt"
	"sort"
	"time"
)

type Issue struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type CheckResult struct {
	OK      bool     `json:"ok"`
	Issues  []Issue  `json:"issues,omitempty"`
	Packet  string   `json:"packet"`
	Checked time.Time `json:"checked"`
}

func CheckPacket(packet Packet) CheckResult {
	result := CheckResult{OK: true, Packet: packet.ID, Checked: time.Now().UTC()}
	if packet.ID == "" {
		result.add("id", "missing packet id")
	}
	if packet.Service == "" {
		result.add("service", "missing service")
	}
	if packet.Version == "" {
		result.add("version", "missing version")
	}
	if packet.Status == "" {
		result.add("status", "missing status")
	}
	if packet.Checkpoint.IsZero() {
		result.add("checkpoint", "missing checkpoint")
	}
	if err := validateHistory(packet); err != nil {
		result.add("history", err.Error())
	}
	return result
}

func validateHistory(packet Packet) error {
	if len(packet.History) == 0 {
		return nil
	}
	previous := packet.History[0].Revision
	for i, entry := range packet.History {
		if entry.Revision < previous {
			return fmt.Errorf("history revision moved backward at %d", i)
		}
		previous = entry.Revision
	}
	return nil
}

func (r *CheckResult) add(field, message string) {
	r.OK = false
	r.Issues = append(r.Issues, Issue{Field: field, Message: message})
}

func (r CheckResult) SortedIssues() []Issue {
	out := append([]Issue(nil), r.Issues...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Field == out[j].Field {
			return out[i].Message < out[j].Message
		}
		return out[i].Field < out[j].Field
	})
	return out
}

func (r CheckResult) FirstIssue() string {
	if len(r.Issues) == 0 {
		return ""
	}
	return r.Issues[0].Message
}

func (r CheckResult) Summary() string {
	if r.OK {
		return "ok"
	}
	return fmt.Sprintf("%d issue(s)", len(r.Issues))
}

