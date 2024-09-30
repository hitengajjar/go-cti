package filesys

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

// Copy directory from src to dst
// remove dst repository if it exists
func ReplaceWithCopy(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		if err = os.RemoveAll(dst); err != nil {
			return fmt.Errorf("remove existing bundle: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat %s: %w", dst, err)
	}

	if err := os.MkdirAll(dst, os.ModePerm); err != nil {
		return fmt.Errorf("create parent directory: %w", err)
	}

	if err := copy.Copy(src, dst); err != nil {
		return fmt.Errorf("copy %s -> %s: %w", src, dst, err)
	}
	return nil
}

func ReplaceWithMove(src, dst string) error {
	if err := os.RemoveAll(dst); err != nil {
		return fmt.Errorf("remove existing dir: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("move %s -> %s: %w", src, dst, err)
	}
	return nil
}
