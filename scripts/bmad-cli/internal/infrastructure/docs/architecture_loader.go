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

// ArchitectureDocs represents all loaded architecture documents
type ArchitectureDocs struct {
	Architecture         ArchitectureDoc
	FrontendArchitecture ArchitectureDoc
	CodingStandards      ArchitectureDoc
	SourceTree           ArchitectureDoc
	TechStack            ArchitectureDoc
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

// LoadAllArchitectureDocsStruct loads all architecture documents and returns them as a struct
func (l *ArchitectureLoader) LoadAllArchitectureDocsStruct() (*ArchitectureDocs, error) {
	// Load architecture document
	archContent, err := l.loadDocumentWithPath("documents.architecture")
	if err != nil {
		return nil, fmt.Errorf("failed to load architecture document: %w", err)
	}

	// Load frontend architecture document
	frontendContent, err := l.loadDocumentWithPath("documents.frontend_architecture")
	if err != nil {
		return nil, fmt.Errorf("failed to load frontend architecture document: %w", err)
	}

	// Load coding standards document
	codingContent, err := l.loadDocumentWithPath("documents.coding_standards")
	if err != nil {
		return nil, fmt.Errorf("failed to load coding standards document: %w", err)
	}

	// Load source tree document
	sourceContent, err := l.loadDocumentWithPath("documents.source_tree")
	if err != nil {
		return nil, fmt.Errorf("failed to load source tree document: %w", err)
	}

	// Load tech stack document
	techContent, err := l.loadDocumentWithPath("documents.tech_stack")
	if err != nil {
		return nil, fmt.Errorf("failed to load tech stack document: %w", err)
	}

	return &ArchitectureDocs{
		Architecture:         archContent,
		FrontendArchitecture: frontendContent,
		CodingStandards:      codingContent,
		SourceTree:           sourceContent,
		TechStack:            techContent,
	}, nil
}

// loadDocumentWithPath loads a document given its config key path
func (l *ArchitectureLoader) loadDocumentWithPath(configKey string) (ArchitectureDoc, error) {
	filepath := l.config.GetString(configKey)
	if filepath == "" {
		return ArchitectureDoc{}, fmt.Errorf("document path not configured for key: %s", configKey)
	}

	content, err := l.loadDocument(filepath)
	if err != nil {
		return ArchitectureDoc{}, fmt.Errorf("failed to load document %s (from %s): %w", configKey, filepath, err)
	}

	return ArchitectureDoc{
		Content:  content,
		FilePath: filepath,
	}, nil
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
