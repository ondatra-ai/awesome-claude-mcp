package docs

import (
	"fmt"
	"os"

	"bmad-cli/internal/infrastructure/config"
)

// ArchitectureDoc represents an architecture document with content and file path
type ArchitectureDoc struct {
	Content  string
	FilePath string
}

// ArchitectureLoader loads architecture documents using configured paths
type ArchitectureLoader struct {
	config *config.ViperConfig
}

// NewArchitectureLoader creates a new ArchitectureLoader instance
func NewArchitectureLoader(config *config.ViperConfig) *ArchitectureLoader {
	return &ArchitectureLoader{
		config: config,
	}
}

// LoadAllArchitectureDocs loads all architecture documents and returns them as a map
func (l *ArchitectureLoader) LoadAllArchitectureDocs() (map[string]ArchitectureDoc, error) {
	docs := make(map[string]ArchitectureDoc)

	// Define the documents to load from config
	docConfigKeys := map[string]string{
		"Architecture":         "documents.architecture",
		"FrontendArchitecture": "documents.frontend_architecture",
		"CodingStandards":      "documents.coding_standards",
		"SourceTree":           "documents.source_tree",
		"TechStack":            "documents.tech_stack",
	}

	// Load each document - fail immediately if any are missing
	for key, configKey := range docConfigKeys {
		filepath := l.config.GetString(configKey)
		if filepath == "" {
			return nil, fmt.Errorf("document path not configured for key: %s", configKey)
		}

		content, err := l.loadDocument(filepath)
		if err != nil {
			return nil, fmt.Errorf("failed to load required architecture document %s (from %s): %w", configKey, filepath, err)
		}
		docs[key] = ArchitectureDoc{
			Content:  content,
			FilePath: filepath,
		}
	}

	return docs, nil
}

// loadDocument loads a single document from the specified path
func (l *ArchitectureLoader) loadDocument(filepath string) (string, error) {

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", fmt.Errorf("document not found: %s", filepath)
	}

	// Read the file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read document %s: %w", filepath, err)
	}

	return string(content), nil
}
