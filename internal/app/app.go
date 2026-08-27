package app

import (
	"fmt"
	"os"

	"example.com/releaseflow/internal/service"
	"example.com/releaseflow/internal/storage"
)

func Run(args []string) error {
	if len(args) == 0 {
		return printUsage()
	}

	cfg := DefaultConfig()
	store, err := storage.NewFileStore(cfg.DataDir)
	if err != nil {
		return err
	}
	svc := service.New(store)
	return svc.Execute(args)
}

func printUsage() error {
	_, err := fmt.Fprintln(os.Stdout, "releaseflow init|submit|review|approve|publish|rollback|inspect|history|export")
	return err
}
