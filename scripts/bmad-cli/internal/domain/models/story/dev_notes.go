package story

// DevNotes represents development notes with flexible structure
// Each entity must have mandatory 'source' and 'description' fields
// Additional fields are flexible and can vary by entity type
type DevNotes map[string]interface{}
