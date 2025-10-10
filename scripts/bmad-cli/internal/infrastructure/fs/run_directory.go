package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RunDirectory manages timestamped run directories for organizing tmp files
type RunDirectory struct {
	runPath string
}

// NewRunDirectory creates a new timestamped run directory
// Format: basePath/YYYY-MM-DD-HH-MM where HH-MM is hours and minutes
func NewRunDirectory(basePath string) (*RunDirectory, error) {
	// Format: YYYY-MM-DD-HH-MM
	timestamp := time.Now().Format("2006-01-02-15-04")
	dirName := timestamp
	runPath := filepath.Join(basePath, dirName)

	if err := os.MkdirAll(runPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create run directory: %w", err)
	}

	return &RunDirectory{
		runPath: runPath,
	}, nil
}

// GetTmpOutPath returns the full path to the run directory
func (rd *RunDirectory) GetTmpOutPath() string {
	return rd.runPath
}
