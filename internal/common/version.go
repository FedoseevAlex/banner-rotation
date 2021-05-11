package common

import (
	"encoding/json"
	"io"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func PrintVersion(destination io.Writer) error {
	if err := json.NewEncoder(destination).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		return err
	}
	return nil
}
