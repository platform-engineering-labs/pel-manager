package cli

import (
	"os"
	"path/filepath"

	filepathx "github.com/platform-engineering-labs/pelx/filepath"
)

func cleanup(path string) error {
	if filepathx.FileExists(filepath.Join(path, "formae", "bin")) {
		err := os.RemoveAll(filepath.Join(path, "formae"))
		if err != nil {
			return err
		}
	}

	return nil
}
