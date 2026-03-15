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

func ErrLoadPromptsFailed(cause error) error {
	return &AppError{
		Category: CategoryAI,
		Code:     "LOAD_PROMPTS_FAILED",
		Message:  "failed to load prompts",
		Cause:    errors.Join(ErrLoadPrompts, cause),
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
	ErrDocumentNotConfigured      = errors.New("document path not configured")
	ErrLoadDocument               = errors.New("failed to load document")
	ErrLoadTemplateFile           = errors.New("failed to load template file")
	ErrInitializeConfig           = errors.New("failed to initialize config")
	ErrCreateAIClient             = errors.New("failed to create AI client")
	ErrLoadStoryFromEpic          = errors.New("failed to load story from epic file")
	ErrFindStoryFile              = errors.New("failed to find story file")
	ErrStoryFileNotFound          = errors.New("story file not found")
	ErrMultipleStoryFiles         = errors.New("multiple story files found")
	ErrStoryFileNotAccessible     = errors.New("story file not accessible")
	ErrReadStoryFile              = errors.New("failed to read story file")
	ErrParseStoryYAML             = errors.New("failed to parse story YAML")
	ErrInvalidStorySlug           = errors.New("invalid story slug")
	ErrParseStoryNumber           = errors.New("failed to parse story number")
	ErrLoadEpic                   = errors.New("failed to load epic")
	ErrStoryIndexOutOfBounds      = errors.New("story index out of bounds")
	ErrInvalidStoryNumberFormat   = errors.New("invalid story number format")
	ErrSearchEpicFiles            = errors.New("failed to search for epic files")
	ErrNoEpicFile                 = errors.New("no epic file found")
	ErrMultipleEpicFiles          = errors.New("multiple epic files found")
	ErrReadEpicFile               = errors.New("failed to read epic file")
	ErrParseEpicYAML              = errors.New("failed to parse epic YAML")
	ErrRegexError                 = errors.New("regex error")
	ErrApplyThreadAction          = errors.New("apply failed for thread")
	ErrTriageError                = errors.New("triage error")
	ErrGetStorySlug               = errors.New("failed to get story slug")
	ErrCloneRequirementsFile      = errors.New("failed to clone requirements file")
	ErrLoadStory                  = errors.New("failed to load story")
	ErrMergeScenarios             = errors.New("failed to merge scenarios")
	ErrReplaceRequirements        = errors.New("failed to replace requirements")
	ErrGenerateTests              = errors.New("failed to generate tests")
	ErrImplementFeatures          = errors.New("failed to implement features")
	ErrLoadUserPromptFailed       = errors.New("failed to load user prompt")
	ErrLoadSystemPromptFailed     = errors.New("failed to load system prompt")
	ErrCreateOutputDirectory      = errors.New("failed to create directory")
	ErrReadRequirementsFile       = errors.New("failed to read requirements.yaml")
	ErrWriteOutputFile            = errors.New("failed to write output file")
	ErrReadOriginalFile           = errors.New("failed to read original")
	ErrCreateBackupFile           = errors.New("failed to create backup")
	ErrReadMergedFile             = errors.New("failed to read merged")
	ErrReplaceFile                = errors.New("failed to replace")
	ErrLoadUserPrompt             = errors.New("failed to load user prompt")
	ErrLoadSystemPrompt           = errors.New("failed to load system prompt")
	ErrMergeScenario              = errors.New("failed to merge scenario")
	ErrParsePendingScenarios      = errors.New("failed to parse pending scenarios")
	ErrReadRequirements           = errors.New("failed to read requirements file")
	ErrUnmarshalRequirements      = errors.New("failed to unmarshal requirements YAML")
	ErrArchUpdateNoContent        = errors.New("no content found in architecture update response")
	ErrModifierMustHaveOneKey     = errors.New("modifier must have exactly one key")
	ErrInvalidStepStatementFormat = errors.New("invalid step statement format")
	ErrHEADDetached               = errors.New("HEAD is detached - please checkout a branch first")
	ErrWorkingTreeDirty           = errors.New(
		"working tree has uncommitted changes - please commit or stash them first",
	)
	ErrNotGitRepository = errors.New("current directory is not a git repository")
	ErrMainBehindOrigin = errors.New(
		"main branch is behind origin/main - please pull the latest changes first",
	)
	ErrInvalidStoryFilename     = errors.New("invalid story filename format")
	ErrInvalidStoryFilenameSlug = errors.New("invalid story filename: slug cannot be empty")
)

func ErrLoadTemplateFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "LOAD_TEMPLATE_FILE_FAILED",
		Message:  "failed to load template file",
		Cause:    errors.Join(ErrLoadTemplateFile, cause),
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
	ErrACMissingSteps       = errors.New("acceptance criterion has no steps")
)

// Filesystem Errors.
var (
	ErrCreateDirectory       = errors.New("failed to create directory")
	ErrCheckWorkingDirectory = errors.New("failed to check working directory")
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

// Parser Errors (for claudecode/internal/parser).
var (
	ErrBufferOverflow    = errors.New("buffer overflow")
	ErrParseContentBlock = errors.New("failed to parse content block")
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
	ErrClaudeStreamClosed        = errors.New("claude stream closed before receiving messages")
	ErrClaudeStreamNoOutput      = errors.New("claude process produced no output")
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

func ErrACMissingStepsError(acID string) error {
	return &AppError{
		Category: CategoryParsing,
		Code:     "AC_MISSING_STEPS",
		Message:  fmt.Sprintf("AC %s has no steps", acID),
		Cause:    ErrACMissingSteps,
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
	ErrStageMismatch             = errors.New("story stage mismatch")
	ErrWriteStoryFile            = errors.New("failed to write story file")
	ErrNoPromptsForStage         = errors.New("no prompts found for stage")
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

func ErrStageMismatchError(storyID, currentStage, commandName, requiredStage string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "STAGE_MISMATCH",
		Message: fmt.Sprintf(
			"story %s is at stage %q, but %s requires stage %q — run the previous stage command first",
			storyID, currentStage, commandName, requiredStage,
		),
		Cause: ErrStageMismatch,
	}
}

func ErrWriteStoryFileFailed(cause error) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "WRITE_STORY_FILE_FAILED",
		Message:  "failed to write story file",
		Cause:    errors.Join(ErrWriteStoryFile, cause),
	}
}

func ErrNoPromptsForStageFailed(stageID string) error {
	return &AppError{
		Category: CategoryInfrastructure,
		Code:     "NO_PROMPTS_FOR_STAGE",
		Message:  "no prompts found for stage: " + stageID,
		Cause:    ErrNoPromptsForStage,
	}
}
