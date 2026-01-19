package errors

import (
	"errors"
	"fmt"
)

// Category represents the type of error.
type Category string

const (
	CategoryAI             Category = "ai"
	CategoryGitHub         Category = "github"
	CategoryParsing        Category = "parsing"
	CategoryInfrastructure Category = "infrastructure"
	CategoryParser         Category = "parser"
)

// AppError represents a structured application error.
type AppError struct {
	Category Category
	Code     string
	Message  string
	Cause    error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.Category, e.Code, e.Message, e.Cause)
	}

	return fmt.Sprintf("[%s:%s] %s", e.Category, e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// AI Errors.
var (
	ErrEmptyOutput               = errors.New("client returned empty output")
	ErrCreateTmpDirectory        = errors.New("failed to create tmp directory")
	ErrLoadData                  = errors.New("failed to load data")
	ErrLoadPrompts               = errors.New("failed to load prompts")
	ErrGenerateContent           = errors.New("failed to generate content")
	ErrWriteResponseFile         = errors.New("failed to write response file")
	ErrParseResponse             = errors.New("failed to parse response")
	ErrValidation                = errors.New("validation failed")
	ErrFileNotFound              = errors.New("file not found")
	ErrParseYAML                 = errors.New("failed to parse YAML")
	ErrKeyNotFoundInYAML         = errors.New("key not found in YAML")
	ErrReadTemplateFile          = errors.New("failed to read template file")
	ErrReadChecklistFile         = errors.New("failed to read checklist file")
	ErrParseTemplate             = errors.New("failed to parse template")
	ErrSendQuery                 = errors.New("failed to send query")
	ErrClaudeReturnedError       = errors.New("claude returned error")
	ErrResponseTooLarge          = errors.New("claude response too large for buffer")
	ErrBuildHeuristicPrompt      = errors.New("failed to build heuristic prompt")
	ErrAIClientExecution         = errors.New("AI client execution failed")
	ErrParseAIOutput             = errors.New("failed to parse AI output")
	ErrBuildImplementationPrompt = errors.New("failed to build implementation prompt")
	ErrAIClientImplementation    = errors.New("AI client implementation failed")
)

func ErrEmptyClientOutput(clientName string) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "EMPTY_OUTPUT",
		Message:  clientName + " returned empty output",
		Cause:    ErrEmptyOutput,
	}
}

func ErrCreateTmpDirectoryFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "CREATE_TMP_DIRECTORY_FAILED",
		Message:  "failed to create tmp directory",
		Cause:    errors.Join(ErrCreateTmpDirectory, cause),
	}
}

func ErrLoadDataFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "LOAD_DATA_FAILED",
		Message:  "failed to load data",
		Cause:    errors.Join(ErrLoadData, cause),
	}
}

func ErrLoadPromptsFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "LOAD_PROMPTS_FAILED",
		Message:  "failed to load prompts",
		Cause:    errors.Join(ErrLoadPrompts, cause),
	}
}

func ErrGenerateContentFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "GENERATE_CONTENT_FAILED",
		Message:  "failed to generate content",
		Cause:    errors.Join(ErrGenerateContent, cause),
	}
}

func ErrGenerateContentWithSystemPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "GENERATE_CONTENT_WITH_SYSTEM_PROMPT_FAILED",
		Message:  "failed to generate content with system prompt",
		Cause:    errors.Join(ErrGenerateContent, cause),
	}
}

func ErrWriteResponseFileFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "WRITE_RESPONSE_FILE_FAILED",
		Message:  "failed to write response file",
		Cause:    errors.Join(ErrWriteResponseFile, cause),
	}
}

func ErrParseResponseFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "PARSE_RESPONSE_FAILED",
		Message:  "failed to parse response",
		Cause:    errors.Join(ErrParseResponse, cause),
	}
}

func ErrValidationFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "VALIDATION_FAILED",
		Message:  "validation failed",
		Cause:    errors.Join(ErrValidation, cause),
	}
}

func ErrYAMLFileNotFound(filePrefix, filePath string) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "YAML_FILE_NOT_FOUND",
		Message:  filePrefix + " file not found: " + filePath,
		Cause:    ErrFileNotFound,
	}
}

func ErrParseYAMLFailed(filePrefix string, cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "PARSE_YAML_FAILED",
		Message:  "failed to parse " + filePrefix + " YAML",
		Cause:    errors.Join(ErrParseYAML, cause),
	}
}

func ErrYAMLKeyNotFound(yamlKey string) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "YAML_KEY_NOT_FOUND",
		Message:  yamlKey + " key not found in YAML",
		Cause:    ErrKeyNotFoundInYAML,
	}
}

func ErrReadTemplateFileFailed(filePath string, cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "READ_TEMPLATE_FILE_FAILED",
		Message:  "failed to read template file " + filePath,
		Cause:    errors.Join(ErrReadTemplateFile, cause),
	}
}

func ErrReadChecklistFileFailed(filePath string, cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "READ_CHECKLIST_FILE_FAILED",
		Message:  "failed to read checklist file " + filePath,
		Cause:    errors.Join(ErrReadChecklistFile, cause),
	}
}

func ErrParseTemplateFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "PARSE_TEMPLATE_FAILED",
		Message:  "failed to parse template",
		Cause:    errors.Join(ErrParseTemplate, cause),
	}
}

func ErrSendQueryFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "SEND_QUERY_FAILED",
		Message:  "failed to send query",
		Cause:    errors.Join(ErrSendQuery, cause),
	}
}

func ErrClaudeError(result string) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "CLAUDE_RETURNED_ERROR",
		Message:  "Claude returned error: " + result,
		Cause:    ErrClaudeReturnedError,
	}
}

func ErrResponseTooLargeForBuffer(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "RESPONSE_TOO_LARGE",
		Message:  "Claude response too large for buffer (using streaming approach)",
		Cause:    errors.Join(ErrResponseTooLarge, cause),
	}
}

func ErrClaudeExecutionFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "CLAUDE_EXECUTION_FAILED",
		Message:  "claude execution failed",
		Cause:    errors.Join(ErrGenerateContent, cause),
	}
}

func ErrBuildHeuristicPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "BUILD_HEURISTIC_PROMPT_FAILED",
		Message:  "failed to build heuristic prompt",
		Cause:    errors.Join(ErrBuildHeuristicPrompt, cause),
	}
}

func ErrAIClientExecutionFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "AI_CLIENT_EXECUTION_FAILED",
		Message:  "AI client execution failed",
		Cause:    errors.Join(ErrAIClientExecution, cause),
	}
}

func ErrParseAIOutputFailed(clientName string, cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "PARSE_AI_OUTPUT_FAILED",
		Message:  "failed to parse " + clientName + " output",
		Cause:    errors.Join(ErrParseAIOutput, cause),
	}
}

func ErrBuildImplementationPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "BUILD_IMPLEMENTATION_PROMPT_FAILED",
		Message:  "failed to build implementation prompt",
		Cause:    errors.Join(ErrBuildImplementationPrompt, cause),
	}
}

func ErrAIClientImplementationFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "AI_CLIENT_IMPLEMENTATION_FAILED",
		Message:  "AI client implementation failed",
		Cause:    errors.Join(ErrAIClientImplementation, cause),
	}
}

// GitHub Errors.
var (
	ErrNoPRFound               = errors.New("no PR found")
	ErrUnexpectedRepoOutput    = errors.New("unexpected repo output")
	ErrNoEligibleReviewThreads = errors.New("no eligible review threads found")
	ErrUnexpectedOutput        = errors.New("unexpected output format")
	ErrNoComments              = errors.New("no comments found in thread")
	ErrGitBranch               = errors.New("git branch command failed")
	ErrGHPRList                = errors.New("gh pr list command failed")
	ErrParsePRList             = errors.New("failed to parse PR list")
	ErrRepoView                = errors.New("repo view command failed")
	ErrGetCurrentBranch        = errors.New("failed to get current branch")
	ErrListPRs                 = errors.New("failed to list PRs for branch")
	ErrCheckWorkingTreeStatus  = errors.New("failed to check working tree status")
	ErrFetchOriginMain         = errors.New("failed to fetch origin/main")
	ErrResolveThread           = errors.New("resolve thread failed")
	ErrReplyThread             = errors.New("reply thread failed")
	ErrGetRepoOwnerAndName     = errors.New("failed to get repository owner and name")
	ErrGraphQLListThreads      = errors.New("graphql list threads failed")
	ErrParseThreadsPage        = errors.New("parse threads page failed")
	ErrCompareMainWithOrigin   = errors.New("failed to compare main with origin/main")
	ErrSwitchBranch            = errors.New("failed to switch to branch")
	ErrCheckoutRemoteBranch    = errors.New("failed to checkout remote branch")
	ErrCreateBranch            = errors.New("failed to create branch")
	ErrPushBranch              = errors.New("failed to push branch")
	ErrSwitchToMain            = errors.New("failed to switch to main")
	ErrDeleteLocalBranch       = errors.New("failed to delete local branch")
	ErrDeleteRemoteBranch      = errors.New("failed to delete remote branch")
	ErrCheckHEADState          = errors.New("failed to check HEAD state")
	ErrCheckGitRepository      = errors.New("failed to check git repository")
	ErrCheckMainBehindOrigin   = errors.New("failed to check if main is behind origin")
	ErrUnrelatedBranch         = errors.New("currently on unrelated branch")
)

// Infrastructure Errors.
var (
	ErrDocumentNotConfigured          = errors.New("document path not configured")
	ErrLoadDocument                   = errors.New("failed to load document")
	ErrLoadTemplateFile               = errors.New("failed to load template file")
	ErrLoadDevNotesPrompt             = errors.New("failed to load devnotes system prompt")
	ErrLoadDevNotesUserPrompt         = errors.New("failed to load devnotes user prompt")
	ErrMandatoryEntityMissing         = errors.New("mandatory entity is missing")
	ErrEntityInvalidType              = errors.New("entity must be a map")
	ErrEntityMissingField             = errors.New("entity is missing mandatory field")
	ErrLoadQASystemPrompt             = errors.New("failed to load QA system prompt")
	ErrLoadQAUserPrompt               = errors.New("failed to load QA user prompt")
	ErrGenerateQAResults              = errors.New("failed to generate QA results")
	ErrLoadQAPrompt                   = errors.New("failed to load QA prompt")
	ErrInvalidRiskLevel               = errors.New("invalid risk level")
	ErrInvalidGateStatus              = errors.New("invalid gate status")
	ErrLoadScenariosSystemPrompt      = errors.New("failed to load scenarios system prompt")
	ErrLoadScenariosUserPrompt        = errors.New("failed to load scenarios user prompt")
	ErrLoadScenariosPrompt            = errors.New("failed to load scenarios prompt")
	ErrLoadTestingSystemPrompt        = errors.New("failed to load testing system prompt")
	ErrLoadTestingUserPrompt          = errors.New("failed to load testing user prompt")
	ErrGenerateTesting                = errors.New("failed to generate testing requirements")
	ErrLoadTestingPrompt              = errors.New("failed to load testing prompt")
	ErrLoadTasksSystemPrompt          = errors.New("failed to load tasks system prompt")
	ErrLoadTasksUserPrompt            = errors.New("failed to load tasks user prompt")
	ErrInitializeConfig               = errors.New("failed to initialize config")
	ErrCreateAIClient                 = errors.New("failed to create AI client")
	ErrLoadStoryFromEpic              = errors.New("failed to load story from epic file")
	ErrLoadArchitectureDocs           = errors.New("failed to load architecture documents")
	ErrGenerateTasks                  = errors.New("failed to generate tasks")
	ErrGenerateDevNotes               = errors.New("failed to generate dev_notes")
	ErrGenerateTestingReqs            = errors.New("failed to generate testing requirements")
	ErrGenerateTestScenarios          = errors.New("failed to generate test scenarios")
	ErrGenerateQA                     = errors.New("failed to generate QA results")
	ErrFindStoryFile                  = errors.New("failed to find story file")
	ErrStoryFileNotFound              = errors.New("story file not found")
	ErrMultipleStoryFiles             = errors.New("multiple story files found")
	ErrStoryFileNotAccessible         = errors.New("story file not accessible")
	ErrReadStoryFile                  = errors.New("failed to read story file")
	ErrParseStoryYAML                 = errors.New("failed to parse story YAML")
	ErrInvalidStorySlug               = errors.New("invalid story slug")
	ErrParseStoryNumber               = errors.New("failed to parse story number")
	ErrLoadEpic                       = errors.New("failed to load epic")
	ErrStoryIndexOutOfBounds          = errors.New("story index out of bounds")
	ErrInvalidStoryNumberFormat       = errors.New("invalid story number format")
	ErrSearchEpicFiles                = errors.New("failed to search for epic files")
	ErrNoEpicFile                     = errors.New("no epic file found")
	ErrMultipleEpicFiles              = errors.New("multiple epic files found")
	ErrReadEpicFile                   = errors.New("failed to read epic file")
	ErrParseEpicYAML                  = errors.New("failed to parse epic YAML")
	ErrCreateStory                    = errors.New("failed to create story")
	ErrProcessTemplate                = errors.New("failed to process template")
	ErrSaveStoryFile                  = errors.New("failed to save story file")
	ErrRegexError                     = errors.New("regex error")
	ErrApplyThreadAction              = errors.New("apply failed for thread")
	ErrTriageError                    = errors.New("triage error")
	ErrGetStorySlug                   = errors.New("failed to get story slug")
	ErrCloneRequirementsFile          = errors.New("failed to clone requirements file")
	ErrLoadStory                      = errors.New("failed to load story")
	ErrMergeScenarios                 = errors.New("failed to merge scenarios")
	ErrReplaceRequirements            = errors.New("failed to replace requirements")
	ErrGenerateTests                  = errors.New("failed to generate tests")
	ErrImplementFeatures              = errors.New("failed to implement features")
	ErrLoadUserPromptFailed           = errors.New("failed to load user prompt")
	ErrLoadSystemPromptFailed         = errors.New("failed to load system prompt")
	ErrCreateOutputDirectory          = errors.New("failed to create directory")
	ErrReadRequirementsFile           = errors.New("failed to read requirements.yaml")
	ErrWriteOutputFile                = errors.New("failed to write output file")
	ErrReadOriginalFile               = errors.New("failed to read original")
	ErrCreateBackupFile               = errors.New("failed to create backup")
	ErrReadMergedFile                 = errors.New("failed to read merged")
	ErrReplaceFile                    = errors.New("failed to replace")
	ErrLoadUserPrompt                 = errors.New("failed to load user prompt")
	ErrLoadSystemPrompt               = errors.New("failed to load system prompt")
	ErrMergeScenario                  = errors.New("failed to merge scenario")
	ErrParsePendingScenarios          = errors.New("failed to parse pending scenarios")
	ErrReadRequirements               = errors.New("failed to read requirements file")
	ErrUnmarshalRequirements          = errors.New("failed to unmarshal requirements YAML")
	ErrAssessmentSummaryEmpty         = errors.New("assessment summary cannot be empty")
	ErrAtLeastOneStrength             = errors.New("at least one strength must be identified")
	ErrRiskLevelMustBeSpecified       = errors.New("risk level must be specified")
	ErrTestabilityScoreRange          = errors.New("testability score must be between 1 and 10")
	ErrImplementationReadinessRange   = errors.New("implementation readiness must be between 1 and 10")
	ErrAtLeastOneTestScenario         = errors.New("at least one test scenario must be specified")
	ErrAIGeneratedNoTasks             = errors.New("AI generated no tasks")
	ErrTestLocationEmpty              = errors.New("test location cannot be empty")
	ErrAtLeastOneFramework            = errors.New("at least one testing framework must be specified")
	ErrAtLeastOneTestingReq           = errors.New("at least one testing requirement must be specified")
	ErrCoverageTargetsMustBeSpecified = errors.New("coverage targets must be specified")
	ErrModifierMustHaveOneKey         = errors.New("modifier must have exactly one key")
	ErrInvalidStepStatementFormat     = errors.New("invalid step statement format")
	ErrHEADDetached                   = errors.New("HEAD is detached - please checkout a branch first")
	ErrWorkingTreeDirty               = errors.New(
		"working tree has uncommitted changes - please commit or stash them first",
	)
	ErrNotGitRepository = errors.New("current directory is not a git repository")
	ErrMainBehindOrigin = errors.New(
		"main branch is behind origin/main - please pull the latest changes first",
	)
	ErrInvalidStoryFilename     = errors.New("invalid story filename format")
	ErrInvalidStoryFilenameSlug = errors.New("invalid story filename: slug cannot be empty")
	ErrSchemaAbsPathFailed      = errors.New("failed to get absolute schema path")
	ErrSchemaFileNotExist       = errors.New("schema file does not exist")
	ErrCreateTempFileFailed     = errors.New("failed to create temporary file")
	ErrWriteYAMLContentFailed   = errors.New("failed to write YAML content")
	ErrCloseTempFileFailed      = errors.New("failed to close temporary file")
	ErrYamaleValidationFailed   = errors.New("yamale validation failed")
)

func ErrDocumentPathNotConfigured(key string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "DOCUMENT_NOT_CONFIGURED",
		Message:  "document path not configured for key: " + key,
		Cause:    ErrDocumentNotConfigured,
	}
}

func ErrLoadDocumentFailed(configKey, filepath string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_DOCUMENT_FAILED",
		Message:  "failed to load required architecture document " + configKey + " (from " + filepath + ")",
		Cause:    errors.Join(ErrLoadDocument, cause),
	}
}

func ErrLoadArchitectureFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_ARCHITECTURE_FAILED",
		Message:  "failed to load architecture document",
		Cause:    errors.Join(ErrLoadDocument, cause),
	}
}

func ErrLoadTemplateFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_TEMPLATE_FILE_FAILED",
		Message:  "failed to load template file",
		Cause:    errors.Join(ErrLoadTemplateFile, cause),
	}
}

func ErrLoadDevNotesPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_DEVNOTES_PROMPT_FAILED",
		Message:  "failed to load devnotes system prompt",
		Cause:    errors.Join(ErrLoadDevNotesPrompt, cause),
	}
}

func ErrLoadDevNotesUserPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_DEVNOTES_USER_PROMPT_FAILED",
		Message:  "failed to load devnotes user prompt",
		Cause:    errors.Join(ErrLoadDevNotesUserPrompt, cause),
	}
}

func ErrMandatoryEntityMissingError(entityName string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "MANDATORY_ENTITY_MISSING",
		Message:  fmt.Sprintf("mandatory entity '%s' is missing", entityName),
		Cause:    ErrMandatoryEntityMissing,
	}
}

func ErrEntityInvalidTypeError(entityName string, gotType interface{}) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "ENTITY_INVALID_TYPE",
		Message:  fmt.Sprintf("entity '%s' must be a map, got %T", entityName, gotType),
		Cause:    ErrEntityInvalidType,
	}
}

func ErrEntityMissingFieldError(entityName string, fieldName string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "ENTITY_MISSING_FIELD",
		Message:  fmt.Sprintf("entity '%s' is missing mandatory '%s' field", entityName, fieldName),
		Cause:    ErrEntityMissingField,
	}
}

func ErrLoadQASystemPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_QA_SYSTEM_PROMPT_FAILED",
		Message:  "failed to load QA system prompt",
		Cause:    errors.Join(ErrLoadQASystemPrompt, cause),
	}
}

func ErrLoadQAUserPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_QA_USER_PROMPT_FAILED",
		Message:  "failed to load QA user prompt",
		Cause:    errors.Join(ErrLoadQAUserPrompt, cause),
	}
}

func ErrGenerateQAResultsFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "GENERATE_QA_RESULTS_FAILED",
		Message:  "failed to generate QA results",
		Cause:    errors.Join(ErrGenerateQAResults, cause),
	}
}

func ErrLoadQAPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_QA_PROMPT_FAILED",
		Message:  "failed to load QA prompt",
		Cause:    errors.Join(ErrLoadQAPrompt, cause),
	}
}

func ErrInvalidRiskLevelError(level string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "INVALID_RISK_LEVEL",
		Message:  fmt.Sprintf("invalid risk level: %s (must be Low, Medium, or High)", level),
		Cause:    ErrInvalidRiskLevel,
	}
}

func ErrInvalidGateStatusError(status string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "INVALID_GATE_STATUS",
		Message:  "invalid gate status: " + status,
		Cause:    ErrInvalidGateStatus,
	}
}

func ErrDocumentNotFound(filepath string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "DOCUMENT_NOT_FOUND",
		Message:  "document not found: " + filepath,
		Cause:    ErrLoadDocument,
	}
}

func ErrReadDocumentFailed(filepath string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "READ_DOCUMENT_FAILED",
		Message:  "failed to read document " + filepath,
		Cause:    errors.Join(ErrLoadDocument, cause),
	}
}

func ErrNoPRFoundForBranch(branch string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "NO_PR_FOUND",
		Message:  "no PR found for branch: " + branch,
		Cause:    ErrNoPRFound,
	}
}

func ErrUnexpectedRepoOutputWithDetails(output string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "UNEXPECTED_REPO_OUTPUT",
		Message:  "unexpected repo view output: " + output,
		Cause:    ErrUnexpectedRepoOutput,
	}
}

func ErrNoEligibleThreads(prNumber int) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "NO_ELIGIBLE_THREADS",
		Message:  fmt.Sprintf("no eligible review threads found for PR %d", prNumber),
		Cause:    ErrNoEligibleReviewThreads,
	}
}

func ErrNoCommentsInThread(threadID string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "NO_COMMENTS",
		Message:  "no comments found in thread " + threadID,
		Cause:    ErrNoComments,
	}
}

func ErrGitBranchFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "GIT_BRANCH_FAILED",
		Message:  "git branch command failed",
		Cause:    errors.Join(ErrGitBranch, cause),
	}
}

func ErrGHPRListFailed(cause error, output string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "GH_PR_LIST_FAILED",
		Message:  "gh pr list failed, out=" + output,
		Cause:    errors.Join(ErrGHPRList, cause),
	}
}

func ErrParsePRListFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "PARSE_PR_LIST_FAILED",
		Message:  "failed to parse PR list",
		Cause:    errors.Join(ErrParsePRList, cause),
	}
}

func ErrRepoViewFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "REPO_VIEW_FAILED",
		Message:  "repo view command failed",
		Cause:    errors.Join(ErrRepoView, cause),
	}
}

func ErrGetCurrentBranchFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "GET_CURRENT_BRANCH_FAILED",
		Message:  "failed to get current branch",
		Cause:    errors.Join(ErrGetCurrentBranch, cause),
	}
}

func ErrListPRsFailed(branch string, cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "LIST_PRS_FAILED",
		Message:  "failed to list PRs for branch " + branch,
		Cause:    errors.Join(ErrListPRs, cause),
	}
}

func ErrCheckWorkingTreeStatusFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "CHECK_WORKING_TREE_STATUS_FAILED",
		Message:  "failed to check working tree status",
		Cause:    errors.Join(ErrCheckWorkingTreeStatus, cause),
	}
}

func ErrFetchOriginMainFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "FETCH_ORIGIN_MAIN_FAILED",
		Message:  "failed to fetch origin/main",
		Cause:    errors.Join(ErrFetchOriginMain, cause),
	}
}

func ErrResolveThreadFailed(cause error, output string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "RESOLVE_THREAD_FAILED",
		Message:  "resolve thread failed, out=" + output,
		Cause:    errors.Join(ErrResolveThread, cause),
	}
}

func ErrReplyThreadFailed(cause error, output string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "REPLY_THREAD_FAILED",
		Message:  "reply thread failed, out=" + output,
		Cause:    errors.Join(ErrReplyThread, cause),
	}
}

func ErrGetRepoOwnerAndNameFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "GET_REPO_OWNER_AND_NAME_FAILED",
		Message:  "failed to get repository owner and name",
		Cause:    errors.Join(ErrGetRepoOwnerAndName, cause),
	}
}

func ErrGraphQLListThreadsFailed(cause error, output string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "GRAPHQL_LIST_THREADS_FAILED",
		Message:  "graphql list threads failed, out=" + output,
		Cause:    errors.Join(ErrGraphQLListThreads, cause),
	}
}

func ErrParseThreadsPageFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "PARSE_THREADS_PAGE_FAILED",
		Message:  "parse threads page failed",
		Cause:    errors.Join(ErrParseThreadsPage, cause),
	}
}

// Parsing Errors.
var (
	ErrRiskScoreNotFound       = errors.New("risk_score not found")
	ErrPreferredOptionNotFound = errors.New("preferred_option not found in YAML")
	ErrSummaryNotFound         = errors.New("summary not found in YAML")
	ErrItemsBlockNotFound      = errors.New("items block not found or empty")
	ErrInvalidBooleanValue     = errors.New("invalid boolean value")
	ErrMissingRequiredItem     = errors.New("missing required item")
	ErrInvalidRiskScoreValue   = errors.New("invalid risk score value")
	ErrExecuteTemplate         = errors.New("failed to execute template")
	ErrParseSummaryYAML        = errors.New("missing summary in YAML")
	ErrParseRiskScore          = errors.New("failed to parse risk score")
)

func ErrItemsMustBeBoolean(key, val string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_BOOLEAN_VALUE",
		Message:  fmt.Sprintf("items[%s] must be boolean, got %q", key, val),
		Cause:    ErrInvalidBooleanValue,
	}
}

func ErrMissingItems(key string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "MISSING_REQUIRED_ITEM",
		Message:  "missing items." + key,
		Cause:    ErrMissingRequiredItem,
	}
}

func ErrInvalidRiskScore(score int) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_RISK_SCORE",
		Message:  fmt.Sprintf("invalid risk_score: %d", score),
		Cause:    ErrInvalidRiskScoreValue,
	}
}

func ErrExecuteTemplateFailed(cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "EXECUTE_TEMPLATE_FAILED",
		Message:  "failed to execute template",
		Cause:    errors.Join(ErrExecuteTemplate, cause),
	}
}

func ErrParseSummaryYAMLFailed(cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "PARSE_SUMMARY_YAML_FAILED",
		Message:  "missing summary in YAML",
		Cause:    errors.Join(ErrParseSummaryYAML, cause),
	}
}

func ErrParseRiskScoreFailed(cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "PARSE_RISK_SCORE_FAILED",
		Message:  "failed to parse risk score",
		Cause:    errors.Join(ErrParseRiskScore, cause),
	}
}

// Validation Errors.
var (
	ErrInvalidOptions       = errors.New("invalid options")
	ErrScenarioValidation   = errors.New("scenario validation failed")
	ErrEmptyScenarioID      = errors.New("scenario ID cannot be empty")
	ErrNoCriteria           = errors.New("scenario must reference at least one acceptance criterion")
	ErrNoSteps              = errors.New("scenario must have at least one step")
	ErrNoKeywordSet         = errors.New("step must have at least one keyword set")
	ErrMultipleKeywords     = errors.New("step must have exactly one keyword set")
	ErrNoGivenStep          = errors.New("scenario must have at least one 'Given' step")
	ErrNoWhenStep           = errors.New("scenario must have at least one 'When' step")
	ErrNoThenStep           = errors.New("scenario must have at least one 'Then' step")
	ErrNoExamples           = errors.New("scenario outline must have at least one example")
	ErrInvalidLevel         = errors.New("level must be integration or e2e")
	ErrInvalidPriority      = errors.New("priority must be P0, P1, P2, or P3")
	ErrUncoveredCriterion   = errors.New("acceptance criterion is not covered by any test scenario")
	ErrNoStatements         = errors.New("step must have at least one statement")
	ErrEmptyStatement       = errors.New("statement cannot be empty")
	ErrInvalidFirstStmt     = errors.New("first statement must be main")
	ErrInvalidFollowingStmt = errors.New("following statement must be 'and' or 'but'")
	ErrInvalidModifier      = errors.New("invalid modifier type")
	ErrEmptyCoverage        = errors.New("coverage value cannot be empty")
	ErrInvalidCoverage      = errors.New("coverage value should be a percentage")
)

// Filesystem Errors.
var (
	ErrCreateDirectory       = errors.New("failed to create directory")
	ErrCheckWorkingDirectory = errors.New("failed to check working directory")
	ErrGetCLIVersion         = errors.New("failed to get CLI version")
	ErrReadConfig            = errors.New("failed to read config")
)

func ErrCreateRunDirectoryFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "CREATE_RUN_DIRECTORY_FAILED",
		Message:  "failed to create run directory",
		Cause:    errors.Join(ErrCreateDirectory, cause),
	}
}

func ErrCheckWorkingDirectoryFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "CHECK_WORKING_DIRECTORY_FAILED",
		Message:  "failed to check working directory",
		Cause:    errors.Join(ErrCheckWorkingDirectory, cause),
	}
}

func ErrGetCLIVersionFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "GET_CLI_VERSION_FAILED",
		Message:  "failed to get CLI version",
		Cause:    errors.Join(ErrGetCLIVersion, cause),
	}
}

func ErrInvalidVersionFormat(version string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "INVALID_VERSION_FORMAT",
		Message:  "invalid version format: " + version,
		Cause:    ErrGetCLIVersion,
	}
}

func ErrReadConfigFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "READ_CONFIG_FAILED",
		Message:  "failed to read config file",
		Cause:    errors.Join(ErrReadConfig, cause),
	}
}

// YAML Errors.
var (
	ErrMarshalYAML   = errors.New("failed to marshal to YAML")
	ErrUnmarshalYAML = errors.New("failed to unmarshal from YAML")
)

func ErrMarshalToYAMLFailed(cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "MARSHAL_YAML_FAILED",
		Message:  "failed to marshal to YAML",
		Cause:    errors.Join(ErrMarshalYAML, cause),
	}
}

func ErrUnmarshalFromYAMLFailed(cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "UNMARSHAL_YAML_FAILED",
		Message:  "failed to unmarshal from YAML",
		Cause:    errors.Join(ErrUnmarshalYAML, cause),
	}
}

func ErrNegativeMaxThinkingTokens(value int) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_MAX_THINKING_TOKENS",
		Message:  fmt.Sprintf("MaxThinkingTokens must be non-negative, got %d", value),
		Cause:    ErrInvalidOptions,
	}
}

func ErrNegativeMaxTurns(value int) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_MAX_TURNS",
		Message:  fmt.Sprintf("MaxTurns must be non-negative, got %d", value),
		Cause:    ErrInvalidOptions,
	}
}

func ErrToolInBothLists(tool string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "TOOL_CONFLICT",
		Message:  fmt.Sprintf("tool '%s' cannot be in both AllowedTools and DisallowedTools", tool),
		Cause:    ErrInvalidOptions,
	}
}

// Helper functions for error checking.
func IsCategory(err error, category Category) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Category == category
	}

	return false
}

func IsAIError(err error) bool {
	return IsCategory(err, CategoryAI)
}

func IsGitHubError(err error) bool {
	return IsCategory(err, CategoryGitHub)
}

func IsParsingError(err error) bool {
	return IsCategory(err, CategoryParsing)
}

// Parser Errors (for claudecode/internal/parser).
var (
	ErrBufferOverflow    = errors.New("buffer overflow")
	ErrParseContentBlock = errors.New("failed to parse content block")
	ErrParseLine         = errors.New("failed to parse line")
)

func ErrBufferSizeExceeded(bufferSize, maxSize int) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "BUFFER_SIZE_EXCEEDED",
		Message:  fmt.Sprintf("buffer size %d exceeds limit %d", bufferSize, maxSize),
		Cause:    ErrBufferOverflow,
	}
}

func ErrParseContentBlockFailed(index int, cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "PARSE_CONTENT_BLOCK_FAILED",
		Message:  fmt.Sprintf("failed to parse content block %d", index),
		Cause:    errors.Join(ErrParseContentBlock, cause),
	}
}

func ErrParseLineFailed(lineIndex int, cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "PARSE_LINE_FAILED",
		Message:  fmt.Sprintf("error parsing line %d", lineIndex),
		Cause:    errors.Join(ErrParseLine, cause),
	}
}

// Transport Errors (for claudecode/internal/subprocess).
var (
	ErrCreateStdinPipe  = errors.New("failed to create stdin pipe")
	ErrCreateStdoutPipe = errors.New("failed to create stdout pipe")
	ErrCreateStderrFile = errors.New("failed to create stderr file")
	ErrMarshalMessage   = errors.New("failed to marshal message")
	ErrWriteMessage     = errors.New("failed to write message")
	ErrStdoutScanner    = errors.New("stdout scanner error")
)

func ErrCreateStdinPipeFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "CREATE_STDIN_PIPE_FAILED",
		Message:  "failed to create stdin pipe",
		Cause:    errors.Join(ErrCreateStdinPipe, cause),
	}
}

func ErrCreateStdoutPipeFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "CREATE_STDOUT_PIPE_FAILED",
		Message:  "failed to create stdout pipe",
		Cause:    errors.Join(ErrCreateStdoutPipe, cause),
	}
}

func ErrCreateStderrFileFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "CREATE_STDERR_FILE_FAILED",
		Message:  "failed to create stderr file",
		Cause:    errors.Join(ErrCreateStderrFile, cause),
	}
}

func ErrMarshalMessageFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "MARSHAL_MESSAGE_FAILED",
		Message:  "failed to marshal message",
		Cause:    errors.Join(ErrMarshalMessage, cause),
	}
}

func ErrWriteMessageFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "WRITE_MESSAGE_FAILED",
		Message:  "failed to write message",
		Cause:    errors.Join(ErrWriteMessage, cause),
	}
}

func ErrStdoutScannerFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "STDOUT_SCANNER_FAILED",
		Message:  "stdout scanner error",
		Cause:    errors.Join(ErrStdoutScanner, cause),
	}
}

// Client Errors (for claudecode).
var (
	ErrWorkingDirectoryNotExist  = errors.New("working directory does not exist")
	ErrInvalidMaxTurns           = errors.New("invalid max_turns")
	ErrInvalidPermissionMode     = errors.New("invalid permission mode")
	ErrConnectClient             = errors.New("failed to connect client")
	ErrCLINotFound               = errors.New("claude CLI not found")
	ErrConnectTransport          = errors.New("failed to connect transport")
	ErrInvalidConfiguration      = errors.New("invalid configuration")
	ErrCloseTransport            = errors.New("failed to close transport")
	ErrClientNotConnected        = errors.New("client not connected")
	ErrTransportAlreadyConnected = errors.New("transport already connected")
	ErrTransportNotConnected     = errors.New("transport not connected or stdin closed")
	ErrProcessNotRunning         = errors.New("process not running")
	ErrInterruptNotSupported     = errors.New("interrupt not supported by windows")
	ErrTransportRequired         = errors.New("transport is required")
	ErrClaudeStreamNilMessage    = errors.New("claude stream returned nil message")
)

func ErrWorkingDirectoryDoesNotExist(path string) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "WORKING_DIRECTORY_NOT_EXIST",
		Message:  "working directory does not exist: " + path,
		Cause:    ErrWorkingDirectoryNotExist,
	}
}

func ErrMaxTurnsMustBeNonNegative(value int) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "INVALID_MAX_TURNS",
		Message:  fmt.Sprintf("max_turns must be non-negative, got: %d", value),
		Cause:    ErrInvalidMaxTurns,
	}
}

func ErrInvalidPermissionModeValue(mode string) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "INVALID_PERMISSION_MODE",
		Message:  "invalid permission mode: " + mode,
		Cause:    ErrInvalidPermissionMode,
	}
}

func ErrConnectClientFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "CONNECT_CLIENT_FAILED",
		Message:  "failed to connect client",
		Cause:    errors.Join(ErrConnectClient, cause),
	}
}

func ErrCLINotFoundError(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "CLI_NOT_FOUND",
		Message:  "claude CLI not found",
		Cause:    errors.Join(ErrCLINotFound, cause),
	}
}

func ErrConnectTransportFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "CONNECT_TRANSPORT_FAILED",
		Message:  "failed to connect transport",
		Cause:    errors.Join(ErrConnectTransport, cause),
	}
}

func ErrInvalidConfigurationError(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "INVALID_CONFIGURATION",
		Message:  "invalid configuration",
		Cause:    errors.Join(ErrInvalidConfiguration, cause),
	}
}

func ErrCloseTransportFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "CLOSE_TRANSPORT_FAILED",
		Message:  "failed to close transport",
		Cause:    errors.Join(ErrCloseTransport, cause),
	}
}

// Query Errors (for claudecode).
var (
	ErrCreateQueryTransport = errors.New("failed to create query transport")
	ErrSendMessage          = errors.New("failed to send message")
)

func ErrCreateQueryTransportFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "CREATE_QUERY_TRANSPORT_FAILED",
		Message:  "failed to create query transport",
		Cause:    errors.Join(ErrCreateQueryTransport, cause),
	}
}

func ErrSendMessageFailed(cause error) error {
	return &AppError{
		Category: CategoryParser,
		Code:     "SEND_MESSAGE_FAILED",
		Message:  "failed to send message",
		Cause:    errors.Join(ErrSendMessage, cause),
	}
}

// Infrastructure constructor functions for new errors.
func ErrLoadScenariosSystemPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_SCENARIOS_SYSTEM_PROMPT_FAILED",
		Message:  "failed to load scenarios system prompt",
		Cause:    errors.Join(ErrLoadScenariosSystemPrompt, cause),
	}
}

func ErrLoadScenariosUserPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_SCENARIOS_USER_PROMPT_FAILED",
		Message:  "failed to load scenarios user prompt",
		Cause:    errors.Join(ErrLoadScenariosUserPrompt, cause),
	}
}

func ErrGenerateTestScenariosFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "GENERATE_TEST_SCENARIOS_FAILED",
		Message:  "failed to generate test scenarios",
		Cause:    errors.Join(ErrGenerateTestScenarios, cause),
	}
}

func ErrLoadScenariosPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_SCENARIOS_PROMPT_FAILED",
		Message:  "failed to load scenarios prompt",
		Cause:    errors.Join(ErrLoadScenariosPrompt, cause),
	}
}

func ErrLoadTestingSystemPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_TESTING_SYSTEM_PROMPT_FAILED",
		Message:  "failed to load testing system prompt",
		Cause:    errors.Join(ErrLoadTestingSystemPrompt, cause),
	}
}

func ErrLoadTestingUserPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_TESTING_USER_PROMPT_FAILED",
		Message:  "failed to load testing user prompt",
		Cause:    errors.Join(ErrLoadTestingUserPrompt, cause),
	}
}

func ErrGenerateTestingFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "GENERATE_TESTING_FAILED",
		Message:  "failed to generate testing requirements",
		Cause:    errors.Join(ErrGenerateTesting, cause),
	}
}

func ErrLoadTestingPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_TESTING_PROMPT_FAILED",
		Message:  "failed to load testing prompt",
		Cause:    errors.Join(ErrLoadTestingPrompt, cause),
	}
}

func ErrLoadTasksSystemPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_TASKS_SYSTEM_PROMPT_FAILED",
		Message:  "failed to load tasks system prompt",
		Cause:    errors.Join(ErrLoadTasksSystemPrompt, cause),
	}
}

func ErrLoadTasksUserPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_TASKS_USER_PROMPT_FAILED",
		Message:  "failed to load tasks user prompt",
		Cause:    errors.Join(ErrLoadTasksUserPrompt, cause),
	}
}

func ErrInitializeConfigFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "INITIALIZE_CONFIG_FAILED",
		Message:  "failed to initialize config",
		Cause:    errors.Join(ErrInitializeConfig, cause),
	}
}

func ErrCreateAIClientFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "CREATE_AI_CLIENT_FAILED",
		Message:  "failed to create AI client",
		Cause:    errors.Join(ErrCreateAIClient, cause),
	}
}

func ErrLoadStoryFromEpicFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_STORY_FROM_EPIC_FAILED",
		Message:  "failed to load story from epic file",
		Cause:    errors.Join(ErrLoadStoryFromEpic, cause),
	}
}

func ErrLoadArchitectureDocsFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_ARCHITECTURE_DOCS_FAILED",
		Message:  "failed to load architecture documents",
		Cause:    errors.Join(ErrLoadArchitectureDocs, cause),
	}
}

func ErrGenerateTasksFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "GENERATE_TASKS_FAILED",
		Message:  "failed to generate tasks",
		Cause:    errors.Join(ErrGenerateTasks, cause),
	}
}

func ErrGenerateDevNotesFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "GENERATE_DEVNOTES_FAILED",
		Message:  "failed to generate dev_notes",
		Cause:    errors.Join(ErrGenerateDevNotes, cause),
	}
}

func ErrGenerateTestingReqsFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "GENERATE_TESTING_REQS_FAILED",
		Message:  "failed to generate testing requirements",
		Cause:    errors.Join(ErrGenerateTestingReqs, cause),
	}
}

func ErrGenerateQAFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "GENERATE_QA_FAILED",
		Message:  "failed to generate QA results",
		Cause:    errors.Join(ErrGenerateQA, cause),
	}
}

func ErrFindStoryFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "FIND_STORY_FILE_FAILED",
		Message:  "failed to find story file",
		Cause:    errors.Join(ErrFindStoryFile, cause),
	}
}

func ErrStoryFileNotFoundError(storyNumber, storiesDir, format string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "STORY_FILE_NOT_FOUND",
		Message: fmt.Sprintf(
			"no story file found for story %s in %s (expected format: %s-<slug>.yaml)",
			storyNumber,
			storiesDir,
			format,
		),
		Cause: ErrStoryFileNotFound,
	}
}

func ErrMultipleStoryFilesError(storyNumber string, matches []string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "MULTIPLE_STORY_FILES",
		Message:  fmt.Sprintf("multiple story files found for story %s: %v", storyNumber, matches),
		Cause:    ErrMultipleStoryFiles,
	}
}

func ErrStoryFileNotAccessibleError(storyFile string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "STORY_FILE_NOT_ACCESSIBLE",
		Message:  "story file not accessible: " + storyFile,
		Cause:    errors.Join(ErrStoryFileNotAccessible, cause),
	}
}

func ErrReadStoryFileFailed(storyFile string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "READ_STORY_FILE_FAILED",
		Message:  "failed to read story file " + storyFile,
		Cause:    errors.Join(ErrReadStoryFile, cause),
	}
}

func ErrParseStoryYAMLFailed(storyFile string, cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "PARSE_STORY_YAML_FAILED",
		Message:  "failed to parse story YAML from " + storyFile,
		Cause:    errors.Join(ErrParseStoryYAML, cause),
	}
}

func ErrInvalidStorySlugError(slug string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "INVALID_STORY_SLUG",
		Message:  fmt.Sprintf("invalid story slug '%s': cannot contain spaces (use dashes instead)", slug),
		Cause:    ErrInvalidStorySlug,
	}
}

func ErrParseStoryNumberFailed(storyNumber string, cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "PARSE_STORY_NUMBER_FAILED",
		Message:  "failed to parse story number " + storyNumber,
		Cause:    errors.Join(ErrParseStoryNumber, cause),
	}
}

func ErrLoadEpicFailed(epicNum int, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_EPIC_FAILED",
		Message:  fmt.Sprintf("failed to load epic %d", epicNum),
		Cause:    errors.Join(ErrLoadEpic, cause),
	}
}

func ErrStoryIndexOutOfBoundsError(storyIndex, epicNum, totalStories int) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "STORY_INDEX_OUT_OF_BOUNDS",
		Message:  fmt.Sprintf("story index %d out of bounds for epic %d (has %d stories)", storyIndex, epicNum, totalStories),
		Cause:    ErrStoryIndexOutOfBounds,
	}
}

func ErrInvalidStoryNumberFormatError(storyNumber string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_STORY_NUMBER_FORMAT",
		Message:  "invalid story number format, expected X.Y but got " + storyNumber,
		Cause:    ErrInvalidStoryNumberFormat,
	}
}

func ErrSearchEpicFilesFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "SEARCH_EPIC_FILES_FAILED",
		Message:  "failed to search for epic files",
		Cause:    errors.Join(ErrSearchEpicFiles, cause),
	}
}

func ErrNoEpicFileError(pattern string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "NO_EPIC_FILE",
		Message:  "no epic file found matching pattern " + pattern,
		Cause:    ErrNoEpicFile,
	}
}

func ErrMultipleEpicFilesError(pattern string, matches []string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "MULTIPLE_EPIC_FILES",
		Message:  fmt.Sprintf("multiple epic files found matching pattern %s: %v", pattern, matches),
		Cause:    ErrMultipleEpicFiles,
	}
}

func ErrReadEpicFileFailed(epicFilePath string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "READ_EPIC_FILE_FAILED",
		Message:  "failed to read epic file " + epicFilePath,
		Cause:    errors.Join(ErrReadEpicFile, cause),
	}
}

func ErrParseEpicYAMLFailed(epicFilePath string, cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "PARSE_EPIC_YAML_FAILED",
		Message:  "failed to parse epic YAML from " + epicFilePath,
		Cause:    errors.Join(ErrParseEpicYAML, cause),
	}
}

func ErrCreateStoryFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "CREATE_STORY_FAILED",
		Message:  "failed to create story",
		Cause:    errors.Join(ErrCreateStory, cause),
	}
}

func ErrProcessTemplateFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "PROCESS_TEMPLATE_FAILED",
		Message:  "failed to process template",
		Cause:    errors.Join(ErrProcessTemplate, cause),
	}
}

func ErrSaveStoryFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "SAVE_STORY_FILE_FAILED",
		Message:  "failed to save story file",
		Cause:    errors.Join(ErrSaveStoryFile, cause),
	}
}

func ErrRegexFailed(cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "REGEX_ERROR",
		Message:  "regex error",
		Cause:    errors.Join(ErrRegexError, cause),
	}
}

func ErrApplyThreadActionFailed(threadID string, cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "APPLY_THREAD_ACTION_FAILED",
		Message:  "apply failed for thread " + threadID,
		Cause:    errors.Join(ErrApplyThreadAction, cause),
	}
}

func ErrTriageFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "TRIAGE_FAILED",
		Message:  "triage",
		Cause:    errors.Join(ErrTriageError, cause),
	}
}

func ErrGetStorySlugFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "GET_STORY_SLUG_FAILED",
		Message:  "failed to get story slug",
		Cause:    errors.Join(ErrGetStorySlug, cause),
	}
}

func ErrCloneRequirementsFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "CLONE_REQUIREMENTS_FILE_FAILED",
		Message:  "failed to clone requirements file",
		Cause:    errors.Join(ErrCloneRequirementsFile, cause),
	}
}

func ErrLoadStoryFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_STORY_FAILED",
		Message:  "failed to load story",
		Cause:    errors.Join(ErrLoadStory, cause),
	}
}

func ErrMergeScenariosFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "MERGE_SCENARIOS_FAILED",
		Message:  "failed to merge scenarios",
		Cause:    errors.Join(ErrMergeScenarios, cause),
	}
}

func ErrReplaceRequirementsFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "REPLACE_REQUIREMENTS_FAILED",
		Message:  "failed to replace requirements",
		Cause:    errors.Join(ErrReplaceRequirements, cause),
	}
}

func ErrGenerateTestsFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "GENERATE_TESTS_FAILED",
		Message:  "failed to generate tests",
		Cause:    errors.Join(ErrGenerateTests, cause),
	}
}

func ErrImplementFeaturesFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "IMPLEMENT_FEATURES_FAILED",
		Message:  "failed to implement features",
		Cause:    errors.Join(ErrImplementFeatures, cause),
	}
}

func ErrImplementFeaturesMaxAttemptsExceeded(maxAttempts int) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "IMPLEMENT_FEATURES_MAX_ATTEMPTS_EXCEEDED",
		Message:  fmt.Sprintf("failed to make tests pass after %d attempts", maxAttempts),
		Cause:    ErrImplementFeatures,
	}
}

func ErrCreateOutputDirectoryFailed(outputDir string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "CREATE_OUTPUT_DIRECTORY_FAILED",
		Message:  "failed to create directory " + outputDir,
		Cause:    errors.Join(ErrCreateOutputDirectory, cause),
	}
}

func ErrReadRequirementsFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "READ_REQUIREMENTS_FILE_FAILED",
		Message:  "failed to read requirements.yaml",
		Cause:    errors.Join(ErrReadRequirementsFile, cause),
	}
}

func ErrWriteOutputFileFailed(outputFile string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "WRITE_OUTPUT_FILE_FAILED",
		Message:  "failed to write " + outputFile,
		Cause:    errors.Join(ErrWriteOutputFile, cause),
	}
}

func ErrReadOriginalFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "READ_ORIGINAL_FILE_FAILED",
		Message:  "failed to read original",
		Cause:    errors.Join(ErrReadOriginalFile, cause),
	}
}

func ErrCreateBackupFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "CREATE_BACKUP_FILE_FAILED",
		Message:  "failed to create backup",
		Cause:    errors.Join(ErrCreateBackupFile, cause),
	}
}

func ErrReadMergedFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "READ_MERGED_FILE_FAILED",
		Message:  "failed to read merged",
		Cause:    errors.Join(ErrReadMergedFile, cause),
	}
}

func ErrReplaceFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "REPLACE_FILE_FAILED",
		Message:  "failed to replace",
		Cause:    errors.Join(ErrReplaceFile, cause),
	}
}

func ErrLoadUserPromptForScenarioFailed(scenarioID string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_USER_PROMPT_FAILED",
		Message:  "failed to load user prompt for scenario " + scenarioID,
		Cause:    errors.Join(ErrLoadUserPrompt, cause),
	}
}

func ErrLoadSystemPromptForScenarioFailed(scenarioID string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_SYSTEM_PROMPT_FAILED",
		Message:  "failed to load system prompt for scenario " + scenarioID,
		Cause:    errors.Join(ErrLoadSystemPrompt, cause),
	}
}

func ErrMergeScenarioFailed(scenarioID string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "MERGE_SCENARIO_FAILED",
		Message:  "failed to merge scenario " + scenarioID,
		Cause:    errors.Join(ErrMergeScenario, cause),
	}
}

func ErrParsePendingScenariosFailed(cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "PARSE_PENDING_SCENARIOS_FAILED",
		Message:  "failed to parse pending scenarios",
		Cause:    errors.Join(ErrParsePendingScenarios, cause),
	}
}

func ErrReadRequirementsFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "READ_REQUIREMENTS_FAILED",
		Message:  "failed to read requirements file",
		Cause:    errors.Join(ErrReadRequirements, cause),
	}
}

func ErrUnmarshalRequirementsFailed(cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "UNMARSHAL_REQUIREMENTS_FAILED",
		Message:  "failed to unmarshal requirements YAML",
		Cause:    errors.Join(ErrUnmarshalRequirements, cause),
	}
}

// Git-related constructor functions.
func ErrCompareMainWithOriginFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "COMPARE_MAIN_WITH_ORIGIN_FAILED",
		Message:  "failed to compare main with origin/main",
		Cause:    errors.Join(ErrCompareMainWithOrigin, cause),
	}
}

func ErrSwitchBranchFailed(branch string, cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "SWITCH_BRANCH_FAILED",
		Message:  "failed to switch to branch " + branch,
		Cause:    errors.Join(ErrSwitchBranch, cause),
	}
}

func ErrCheckoutRemoteBranchFailed(branch string, cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "CHECKOUT_REMOTE_BRANCH_FAILED",
		Message:  "failed to checkout remote branch " + branch,
		Cause:    errors.Join(ErrCheckoutRemoteBranch, cause),
	}
}

func ErrCreateBranchFailed(branch string, cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "CREATE_BRANCH_FAILED",
		Message:  "failed to create branch " + branch,
		Cause:    errors.Join(ErrCreateBranch, cause),
	}
}

func ErrPushBranchFailed(branch string, cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "PUSH_BRANCH_FAILED",
		Message:  "failed to push branch " + branch,
		Cause:    errors.Join(ErrPushBranch, cause),
	}
}

func ErrSwitchToMainFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "SWITCH_TO_MAIN_FAILED",
		Message:  "failed to switch to main",
		Cause:    errors.Join(ErrSwitchToMain, cause),
	}
}

func ErrDeleteLocalBranchFailed(branch string, cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "DELETE_LOCAL_BRANCH_FAILED",
		Message:  "failed to delete local branch " + branch,
		Cause:    errors.Join(ErrDeleteLocalBranch, cause),
	}
}

func ErrDeleteRemoteBranchFailed(branch string, cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "DELETE_REMOTE_BRANCH_FAILED",
		Message:  "failed to delete remote branch " + branch,
		Cause:    errors.Join(ErrDeleteRemoteBranch, cause),
	}
}

func ErrCheckHEADStateFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "CHECK_HEAD_STATE_FAILED",
		Message:  "failed to check HEAD state",
		Cause:    errors.Join(ErrCheckHEADState, cause),
	}
}

func ErrCheckGitRepositoryFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "CHECK_GIT_REPOSITORY_FAILED",
		Message:  "failed to check git repository",
		Cause:    errors.Join(ErrCheckGitRepository, cause),
	}
}

func ErrCheckMainBehindOriginFailed(cause error) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "CHECK_MAIN_BEHIND_ORIGIN_FAILED",
		Message:  "failed to check if main is behind origin",
		Cause:    errors.Join(ErrCheckMainBehindOrigin, cause),
	}
}

func ErrUnrelatedBranchError(currentBranch, storyNumber string) error {
	return &AppError{
		Category: CategoryGitHub,
		Code:     "UNRELATED_BRANCH",
		Message: fmt.Sprintf(
			"currently on branch '%s' which is not related to story %s - please switch to main first",
			currentBranch,
			storyNumber,
		),
		Cause: ErrUnrelatedBranch,
	}
}

// Validation constructor functions.
func ErrEmptyScenarioIDError(index int) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "EMPTY_SCENARIO_ID",
		Message:  fmt.Sprintf("scenario %d: ID cannot be empty", index),
		Cause:    ErrEmptyScenarioID,
	}
}

func ErrNoCriteriaError(scenarioID string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "NO_CRITERIA",
		Message:  fmt.Sprintf("scenario %s: must reference at least one acceptance criterion", scenarioID),
		Cause:    ErrNoCriteria,
	}
}

func ErrNoStepsError(scenarioID string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "NO_STEPS",
		Message:  fmt.Sprintf("scenario %s: must have at least one step", scenarioID),
		Cause:    ErrNoSteps,
	}
}

func ErrNoKeywordSetError(scenarioID string, stepIdx int) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "NO_KEYWORD_SET",
		Message:  fmt.Sprintf("scenario %s, step %d: step must have at least one keyword set", scenarioID, stepIdx),
		Cause:    ErrNoKeywordSet,
	}
}

func ErrMultipleKeywordsError(scenarioID string, stepIdx int) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "MULTIPLE_KEYWORDS",
		Message:  fmt.Sprintf("scenario %s, step %d: step must have exactly one keyword set", scenarioID, stepIdx),
		Cause:    ErrMultipleKeywords,
	}
}

func ErrNoGivenStepError(scenarioID string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "NO_GIVEN_STEP",
		Message:  fmt.Sprintf("scenario %s: must have at least one 'Given' step", scenarioID),
		Cause:    ErrNoGivenStep,
	}
}

func ErrNoWhenStepError(scenarioID string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "NO_WHEN_STEP",
		Message:  fmt.Sprintf("scenario %s: must have at least one 'When' step", scenarioID),
		Cause:    ErrNoWhenStep,
	}
}

func ErrNoThenStepError(scenarioID string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "NO_THEN_STEP",
		Message:  fmt.Sprintf("scenario %s: must have at least one 'Then' step", scenarioID),
		Cause:    ErrNoThenStep,
	}
}

func ErrNoExamplesError(scenarioID string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "NO_EXAMPLES",
		Message:  fmt.Sprintf("scenario %s: scenario outline must have at least one example", scenarioID),
		Cause:    ErrNoExamples,
	}
}

func ErrInvalidLevelError(scenarioID string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_LEVEL",
		Message: fmt.Sprintf(
			"scenario %s: level must be integration or e2e (unit scenarios are not allowed in BDD)",
			scenarioID,
		),
		Cause: ErrInvalidLevel,
	}
}

func ErrInvalidPriorityError(scenarioID string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_PRIORITY",
		Message:  fmt.Sprintf("scenario %s: priority must be P0, P1, P2, or P3", scenarioID),
		Cause:    ErrInvalidPriority,
	}
}

func ErrUncoveredCriterionError(acID string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "UNCOVERED_CRITERION",
		Message:  fmt.Sprintf("acceptance criterion %s is not covered by any test scenario", acID),
		Cause:    ErrUncoveredCriterion,
	}
}

func ErrNoStatementsError(scenarioID string, stepIdx int, keyword string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "NO_STATEMENTS",
		Message: fmt.Sprintf(
			"scenario %s, step %d: %s must have at least one statement",
			scenarioID,
			stepIdx,
			keyword,
		),
		Cause: ErrNoStatements,
	}
}

func ErrEmptyStatementError(
	scenarioID string,
	stepIdx int,
	keyword string,
	stmtIdx int,
) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "EMPTY_STATEMENT",
		Message: fmt.Sprintf(
			"scenario %s, step %d: %s statement[%d] cannot be empty",
			scenarioID,
			stepIdx,
			keyword,
			stmtIdx,
		),
		Cause: ErrEmptyStatement,
	}
}

func ErrInvalidFirstStmtError(
	scenarioID string,
	stepIdx int,
	keyword string,
) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_FIRST_STMT",
		Message: fmt.Sprintf(
			"scenario %s, step %d: %s first statement must be main (not 'and' or 'but')",
			scenarioID,
			stepIdx,
			keyword,
		),
		Cause: ErrInvalidFirstStmt,
	}
}

func ErrInvalidFollowingStmtError(
	scenarioID string,
	stepIdx int,
	keyword string,
	stmtIdx int,
) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_FOLLOWING_STMT",
		Message: fmt.Sprintf(
			"scenario %s, step %d: %s statement[%d] must be 'and' or 'but'",
			scenarioID,
			stepIdx,
			keyword,
			stmtIdx,
		),
		Cause: ErrInvalidFollowingStmt,
	}
}

func ErrInvalidModifierError(modifierType string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_MODIFIER",
		Message:  fmt.Sprintf("invalid modifier type: %s (must be 'and' or 'but')", modifierType),
		Cause:    ErrInvalidModifier,
	}
}

func ErrEmptyCoverageError(key string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "EMPTY_COVERAGE",
		Message:  "coverage value for " + key + " cannot be empty",
		Cause:    ErrEmptyCoverage,
	}
}

func ErrInvalidCoverageError(key string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_COVERAGE",
		Message:  "coverage value for " + key + " should be a percentage",
		Cause:    ErrInvalidCoverage,
	}
}

func ErrInvalidSteps(cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "INVALID_STEPS",
		Message:  "invalid steps parameter",
		Cause:    cause,
	}
}

// Test Execution Errors.
var (
	ErrBaselineTestsFailed     = errors.New("baseline tests failed")
	ErrGeneratedTestsPass      = errors.New("generated tests should fail but passed")
	ErrRunTests                = errors.New("failed to run tests")
	ErrValidateTests           = errors.New("test validation failed")
	ErrValidateScenarios       = errors.New("scenario validation failed")
	ErrMissingScenarioCoverage = errors.New("missing scenario coverage in tests")
	ErrUnfixedTestIssues       = errors.New("unfixed test quality issues")
	ErrReadResultFile          = errors.New("failed to read result file")
	ErrParseResultYAML         = errors.New("failed to parse result YAML")
)

func ErrBaselineTestsFailedError(output string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "BASELINE_TESTS_FAILED",
		Message:  "baseline tests must pass before test generation - tests are failing. Output: " + output,
		Cause:    ErrBaselineTestsFailed,
	}
}

func ErrGeneratedTestsPassError(output string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "GENERATED_TESTS_PASS",
		Message:  "generated tests should fail (TDD red phase) but all tests passed. Output: " + output,
		Cause:    ErrGeneratedTestsPass,
	}
}

func ErrRunTestsFailed(phase string, cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "RUN_TESTS_FAILED",
		Message:  "failed to run tests during " + phase + " phase",
		Cause:    errors.Join(ErrRunTests, cause),
	}
}

func ErrValidateTestsFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "VALIDATE_TESTS_FAILED",
		Message:  "test validation failed - tests contain quality issues",
		Cause:    errors.Join(ErrValidateTests, cause),
	}
}

func ErrValidateScenariosFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "VALIDATE_SCENARIOS_FAILED",
		Message:  "scenario validation failed",
		Cause:    errors.Join(ErrValidateScenarios, cause),
	}
}

func ErrMissingScenarioCoverageError(missingScenarios []string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "MISSING_SCENARIO_COVERAGE",
		Message:  fmt.Sprintf("missing test coverage for scenarios: %v", missingScenarios),
		Cause:    ErrMissingScenarioCoverage,
	}
}

func ErrUnfixedTestIssuesError(count int) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "UNFIXED_TEST_ISSUES",
		Message:  fmt.Sprintf("%d issues could not be automatically fixed", count),
		Cause:    ErrUnfixedTestIssues,
	}
}

func ErrReadResultFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "READ_RESULT_FILE_FAILED",
		Message:  "failed to read result file",
		Cause:    errors.Join(ErrReadResultFile, cause),
	}
}

func ErrParseResultYAMLFailed(cause error) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "PARSE_RESULT_YAML_FAILED",
		Message:  "failed to parse result YAML",
		Cause:    errors.Join(ErrParseResultYAML, cause),
	}
}

// ErrImplementFailed wraps implementation errors.
func ErrImplementFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "IMPLEMENT_FAILED",
		Message:  "implementation failed",
		Cause:    cause,
	}
}

// Checklist Validation Errors.
var (
	ErrLoadChecklistSystemPrompt = errors.New("failed to load checklist system prompt")
	ErrLoadChecklistUserPrompt   = errors.New("failed to load checklist user prompt")
	ErrChecklistAIEvaluation     = errors.New("AI evaluation failed")
	ErrFixApplierNoContent       = errors.New("no FILE_START/FILE_END content found")
	ErrFixPromptGeneration       = errors.New("fix prompt generation failed")
	ErrFixPromptRefinement       = errors.New("fix prompt refinement failed")
	ErrSaveStoryVersion          = errors.New("failed to save story version")
)

func ErrLoadChecklistSystemPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "LOAD_CHECKLIST_SYSTEM_PROMPT_FAILED",
		Message:  "failed to load checklist system prompt template",
		Cause:    errors.Join(ErrLoadChecklistSystemPrompt, cause),
	}
}

func ErrLoadChecklistUserPromptFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "LOAD_CHECKLIST_USER_PROMPT_FAILED",
		Message:  "failed to load checklist user prompt template",
		Cause:    errors.Join(ErrLoadChecklistUserPrompt, cause),
	}
}

func ErrChecklistAIEvaluationFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "CHECKLIST_AI_EVALUATION_FAILED",
		Message:  "AI evaluation of checklist prompt failed",
		Cause:    errors.Join(ErrChecklistAIEvaluation, cause),
	}
}

func ErrFixApplierNoContentFound(resultPath string) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "FIX_APPLIER_NO_CONTENT",
		Message:  "no FILE_START/FILE_END content found for path: " + resultPath,
		Cause:    ErrFixApplierNoContent,
	}
}

func ErrFixPromptGenerationFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "FIX_PROMPT_GENERATION_FAILED",
		Message:  "fix prompt generation failed",
		Cause:    errors.Join(ErrFixPromptGeneration, cause),
	}
}

func ErrFixPromptRefinementFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "FIX_PROMPT_REFINEMENT_FAILED",
		Message:  "fix prompt refinement failed",
		Cause:    errors.Join(ErrFixPromptRefinement, cause),
	}
}

func ErrSaveStoryVersionFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "SAVE_STORY_VERSION_FAILED",
		Message:  "failed to save story version",
		Cause:    errors.Join(ErrSaveStoryVersion, cause),
	}
}
