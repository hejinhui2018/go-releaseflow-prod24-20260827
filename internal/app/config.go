package app

import "path/filepath"

type Config struct {
	DataDir string
}

func DefaultConfig() Config {
	return Config{DataDir: filepath.Join(".", "releaseflow-data")}
}

