package docs

import (
	"os"

	"bmad-cli/internal/infrastructure/config"
	"bmad-cli/internal/pkg/errors"
)

// ArchitectureDoc represents an architecture document with content and file path.
type ArchitectureDoc struct {
	Content  string
	FilePath string
}

// ArchitectureDocs represents all loaded architecture documents.
type ArchitectureDocs struct {
	Architecture         ArchitectureDoc
	FrontendArchitecture ArchitectureDoc
	CodingStandards      ArchitectureDoc
	SourceTree           ArchitectureDoc
	TechStack            ArchitectureDoc
}

// ArchitectureLoader loads architecture documents using configured paths.
type ArchitectureLoader struct {
	config *config.ViperConfig
}

// NewArchitectureLoader creates a new ArchitectureLoader instance.
func NewArchitectureLoader(config *config.ViperConfig) *ArchitectureLoader {
	return &ArchitectureLoader{
		config: config,
	}
}

// LoadAllArchitectureDocs loads all architecture documents and returns them as a map.
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
			return nil, errors.ErrDocumentPathNotConfigured(configKey)
		}

		content, err := l.loadDocument(filepath)
		if err != nil {
			return nil, errors.ErrLoadDocumentFailed(configKey, filepath, err)
		}

		docs[key] = ArchitectureDoc{
			Content:  content,
			FilePath: filepath,
		}
	}

	return docs, nil
}

// LoadAllArchitectureDocsStruct loads all architecture documents and returns them as a struct.
func (l *ArchitectureLoader) LoadAllArchitectureDocsStruct() (*ArchitectureDocs, error) {
	// Load architecture document
	archContent, err := l.loadDocumentWithPath("documents.architecture")
	if err != nil {
		return nil, errors.ErrLoadArchitectureFailed(err)
	}

	// Load frontend architecture document
	frontendContent, err := l.loadDocumentWithPath("documents.frontend_architecture")
	if err != nil {
		return nil, errors.ErrLoadDocumentFailed("frontend_architecture", "documents.frontend_architecture", err)
	}

	// Load coding standards document
	codingContent, err := l.loadDocumentWithPath("documents.coding_standards")
	if err != nil {
		return nil, errors.ErrLoadDocumentFailed("coding_standards", "documents.coding_standards", err)
	}

	// Load source tree document
	sourceContent, err := l.loadDocumentWithPath("documents.source_tree")
	if err != nil {
		return nil, errors.ErrLoadDocumentFailed("source_tree", "documents.source_tree", err)
	}

	// Load tech stack document
	techContent, err := l.loadDocumentWithPath("documents.tech_stack")
	if err != nil {
		return nil, errors.ErrLoadDocumentFailed("tech_stack", "documents.tech_stack", err)
	}

	return &ArchitectureDocs{
		Architecture:         archContent,
		FrontendArchitecture: frontendContent,
		CodingStandards:      codingContent,
		SourceTree:           sourceContent,
		TechStack:            techContent,
	}, nil
}

// loadDocumentWithPath loads a document given its config key path.
func (l *ArchitectureLoader) loadDocumentWithPath(configKey string) (ArchitectureDoc, error) {
	filepath := l.config.GetString(configKey)
	if filepath == "" {
		return ArchitectureDoc{}, errors.ErrDocumentPathNotConfigured(configKey)
	}

	content, err := l.loadDocument(filepath)
	if err != nil {
		return ArchitectureDoc{}, errors.ErrLoadDocumentFailed(configKey, filepath, err)
	}

	return ArchitectureDoc{
		Content:  content,
		FilePath: filepath,
	}, nil
}

// loadDocument loads a single document from the specified path.
func (l *ArchitectureLoader) loadDocument(filepath string) (string, error) {
	// Check if file exists
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return "", errors.ErrDocumentNotFound(filepath)
	}

	// Read the file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", errors.ErrReadDocumentFailed(filepath, err)
	}

	return string(content), nil
}
