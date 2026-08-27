package service

import (
	"fmt"
	"io"
	"os"

	"example.com/releaseflow/internal/model"
	"example.com/releaseflow/internal/storage"
)

type Service struct {
	store  storage.Store
	clock  model.Clock
	policy model.Policy
	out    io.Writer
}

func New(store storage.Store) *Service {
	return &Service{
		store:  store,
		clock:  model.SystemClock{},
		policy: model.DefaultPolicy(),
		out:    os.Stdout,
	}
}

func (s *Service) Execute(args []string) error {
	cmd, err := Parse(args)
	if err != nil {
		return err
	}
	switch cmd.Name {
	case "init":
		return s.Init(cmd)
	case "submit":
		return s.Submit(cmd)
	case "review":
		return s.Review(cmd)
	case "approve":
		return s.Approve(cmd)
	case "publish":
		return s.Publish(cmd)
	case "rollback":
		return s.Rollback(cmd)
	case "inspect":
		return s.Inspect(cmd)
	case "history":
		return s.History(cmd)
	case "export":
		return s.Export(cmd)
	case "recover":
		return s.Recover(cmd)
	case "list":
		return s.List(cmd)
	case "metrics":
		return s.Metrics(cmd)
	case "report":
		return s.Report(cmd)
	case "backup":
		return s.Backup(cmd)
	case "restore":
		return s.Restore(cmd)
	case "clean":
		return s.Clean(cmd)
	case "summary":
		return s.Summary(cmd)
	case "quality":
		return s.Quality(cmd)
	case "board":
		return s.Board(cmd)
	case "status":
		return s.StatusLine(cmd)
	case "next":
		return s.NextAction(cmd)
	default:
		return fmt.Errorf("unknown command %q", cmd.Name)
	}
}
