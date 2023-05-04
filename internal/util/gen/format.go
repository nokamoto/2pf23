package gen

import (
	"go/format"
	"os"
)

// WriteFormattedGo writes a go file with format.
func WriteFormattedGo(file string, b []byte) error {
	formatted, err := format.Source(b)
	if err != nil {
		return err
	}
	return os.WriteFile(file, formatted, 0o644)
}
