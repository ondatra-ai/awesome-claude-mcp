package docs

import (
	"fmt"
	"os"

	"bmad-cli/internal/infrastructure/config"
)

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
func (l *ArchitectureLoader) LoadAllArchitectureDocs() (map[string]string, error) {
	docs := make(map[string]string)

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
		docs[key] = content
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

// LoadSpecificDoc loads a specific architecture document
func (l *ArchitectureLoader) LoadSpecificDoc(docType string) (string, error) {
	docConfigKeys := map[string]string{
		"Architecture":         "documents.architecture",
		"FrontendArchitecture": "documents.frontend_architecture",
		"CodingStandards":      "documents.coding_standards",
		"SourceTree":           "documents.source_tree",
		"TechStack":            "documents.tech_stack",
	}

	configKey, exists := docConfigKeys[docType]
	if !exists {
		return "", fmt.Errorf("unknown document type: %s", docType)
	}

	filepath := l.config.GetString(configKey)
	if filepath == "" {
		return "", fmt.Errorf("document path not configured for key: %s", configKey)
	}

	return l.loadDocument(filepath)
}

// GetAvailableDocTypes returns the list of available document types
func (l *ArchitectureLoader) GetAvailableDocTypes() []string {
	return []string{
		"Architecture",
		"FrontendArchitecture",
		"CodingStandards",
		"SourceTree",
		"TechStack",
	}
}
