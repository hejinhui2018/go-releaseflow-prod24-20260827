package storage

import (
	"encoding/json"
	"fmt"
)

func encode(v any) ([]byte, error) {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}
	return append(raw, '\n'), nil
}

func encodeJournal(v any) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("encode journal: %w", err)
	}
	return append(raw, '\n'), nil
}

func decode(raw []byte, v any) error {
	if err := json.Unmarshal(raw, v); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	return nil
}
