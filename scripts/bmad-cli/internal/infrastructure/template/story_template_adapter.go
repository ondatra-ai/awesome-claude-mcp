package template

import (
	"bmad-cli/internal/domain/models/story"
)

// FlattenedStoryData represents a flattened version of StoryDocument for template processing
// This matches the structure expected by the story.yaml.tpl template
type FlattenedStoryData struct {
	// Story fields at root level
	ID                 string                       `json:"id"`
	Title              string                       `json:"title"`
	AsA                string                       `json:"as_a"`
	IWant              string                       `json:"i_want"`
	SoThat             string                       `json:"so_that"`
	Status             string                       `json:"status"`
	AcceptanceCriteria []story.AcceptanceCriterion  `json:"acceptance_criteria"`

	// Other sections as nested structures
	Tasks          []story.Task           `json:"tasks"`
	DevNotes       story.DevNotes         `json:"dev_notes"`
	Testing        story.Testing          `json:"testing"`
	ChangeLog      []story.ChangeLogEntry `json:"change_log"`
	QAResults      *story.QAResults       `json:"qa_results,omitempty"`
	DevAgentRecord story.DevAgentRecord   `json:"dev_agent_record"`
}

// FlattenStoryDocument converts a StoryDocument to FlattenedStoryData for template processing
func FlattenStoryDocument(doc *story.StoryDocument) *FlattenedStoryData {
	return &FlattenedStoryData{
		// Flatten story fields to root level
		ID:                 doc.Story.ID,
		Title:              doc.Story.Title,
		AsA:                doc.Story.AsA,
		IWant:              doc.Story.IWant,
		SoThat:             doc.Story.SoThat,
		Status:             doc.Story.Status,
		AcceptanceCriteria: doc.Story.AcceptanceCriteria,

		// Keep other sections as nested structures
		Tasks:          doc.Tasks,
		DevNotes:       doc.DevNotes,
		Testing:        doc.Testing,
		ChangeLog:      doc.ChangeLog,
		QAResults:      doc.QAResults,
		DevAgentRecord: doc.DevAgentRecord,
	}
}
