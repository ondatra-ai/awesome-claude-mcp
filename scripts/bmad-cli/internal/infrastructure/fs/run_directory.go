package fs

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RunDirectory manages timestamped run directories for organizing tmp files
type RunDirectory struct {
	basePath string
	runPath  string
}

// NewRunDirectory creates a new timestamped run directory
// Format: basePath/YYYY-MM-DD-XXXX where XXXX is a 4-character hex salt
func NewRunDirectory(basePath string) (*RunDirectory, error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Format: YYYY-MM-DD-XXXX
	timestamp := time.Now().Format("2006-01-02")
	dirName := fmt.Sprintf("%s-%s", timestamp, salt)
	runPath := filepath.Join(basePath, dirName)

	if err := os.MkdirAll(runPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create run directory: %w", err)
	}

	return &RunDirectory{
		basePath: basePath,
		runPath:  runPath,
	}, nil
}

// GetPath returns the full path to the run directory
func (rd *RunDirectory) GetPath() string {
	return rd.runPath
}

// GetFilePath returns the full path to a file within the run directory
func (rd *RunDirectory) GetFilePath(filename string) string {
	return filepath.Join(rd.runPath, filename)
}

// generateSalt generates a 4-character hex string from 2 random bytes
func generateSalt() (string, error) {
	bytes := make([]byte, 2) // 2 bytes = 4 hex chars
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
