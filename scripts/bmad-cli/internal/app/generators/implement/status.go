package implement

// File permission constant for read-write files.
const fileModeReadWrite = 0o644

// GenerationStatus represents the result of a generator that modifies files.
// Used by implement generators to report what was done.
type GenerationStatus struct {
	FilesModified  []string
	ItemsProcessed int
	Success        bool
	Message        string
}

// NewSuccessStatus creates a successful status with the given details.
func NewSuccessStatus(itemsProcessed int, filesModified []string, message string) GenerationStatus {
	return GenerationStatus{
		FilesModified:  filesModified,
		ItemsProcessed: itemsProcessed,
		Success:        true,
		Message:        message,
	}
}

// NewFailureStatus creates a failed status with the given message.
func NewFailureStatus(message string) GenerationStatus {
	return GenerationStatus{
		FilesModified:  nil,
		ItemsProcessed: 0,
		Success:        false,
		Message:        message,
	}
}
