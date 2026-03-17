# Checklist Validation Algorithm

This document describes how BMAD CLI validates artifacts (user stories and generated tests) against YAML-based checklists.

## Overview

The checklist validation system is an AI-driven quality gate. A YAML checklist defines questions with expected answers. For each question, the CLI sends the artifact and the question to Claude, parses the AI's answer, and compares it to the expected answer. The result is PASS, WARN, FAIL, or SKIP per question.

There are two checklist files, both validated against the same Yamale schema at build time:

| Checklist file | Purpose | CLI entry point |
|---|---|---|
| `bdd-cli/user-story-description-checklist.yaml` | Validates user stories across 4 stages | `bmad-cli us create/refine/ready` |
| `bdd-cli/test-validation-checklist.yaml` | Validates generated Playwright tests | `bmad-cli req generate_tests` |

## YAML Structure

```yaml
version: "3.0"
default_docs: [architecture, prd]   # optional — doc keys loaded for every prompt
stages:
  - id: story_creation
    name: Story Creation
    sections:
      - id: format
        name: Format
        validation_prompts:
          - Q: "Does the story follow the template?"
            A: "yes"
            rationale: "optional — why this matters"
            skip: ""             # non-empty → skipped
            docs: [prd]          # overrides default_docs for this prompt
            F: "fix template…"   # template for generating fix prompt on failure
```

Key fields on each prompt:

| Field | Role |
|---|---|
| `Q` | Question sent to Claude about the artifact |
| `A` | Expected answer — supports exact match, ranges, comparisons (see below) |
| `rationale` | Included in the AI prompt for context |
| `skip` | If non-empty, the prompt is excluded from evaluation |
| `docs` | Document keys resolved to file paths via config; overrides `default_docs` |
| `F` | Markdown template used by the fix-prompt generator when this check fails |

## Algorithm

### High-Level Flow

```mermaid
flowchart TD
    A[CLI Command<br>us create / req generate_tests] --> B[Load Checklist YAML]
    B --> C[Extract Prompts for Stage]
    C --> D{Fix mode?}
    D -->|No| E[Evaluate ALL prompts]
    D -->|Yes| F[Evaluate until first FAIL]
    E --> G[Render Report Table]
    F --> G
    G --> H{All passed?}
    H -->|Yes| I[Save artifact & advance stage]
    H -->|No, not fix mode| J[Show failures & exit]
    H -->|No, fix mode| K[Enter Fix Loop]
    K --> F
```

### 1. Load and Filter

```mermaid
flowchart LR
    A[YAML File] -->|yaml.Unmarshal| B[Checklist struct]
    B -->|ExtractPromptsForStage| C{For each section}
    C --> D{For each prompt}
    D --> E{skip non-empty?}
    E -->|Yes| F[Skip]
    E -->|No| G[Wrap in PromptWithContext<br>add stage/section metadata]
    G --> H["[]PromptWithContext"]
```

### 2. Evaluate Each Prompt

```mermaid
sequenceDiagram
    participant E as ChecklistEvaluator
    participant C as Config
    participant T as Template Engine
    participant AI as Claude API
    participant P as Parser

    E->>C: Resolve doc keys → file paths
    E->>T: Render system prompt (cached)
    E->>T: Render user prompt with data
    Note over T: Subject, Question, Rationale,<br>ResultPath, Docs, FixTemplate
    E->>AI: ExecutePromptWithSystem<br>(think mode)
    AI-->>E: Response with FILE_START/FILE_END
    E->>P: Extract YAML between markers
    P-->>E: {answer, fix_prompt?}
    E->>E: compareAnswers(expected, actual, acCount)
    E-->>E: PASS | WARN | FAIL
```

All prompts, responses, and parsed results are saved to a temp directory for debugging.

### 3. Answer Comparison Rules

```mermaid
flowchart TD
    A[compareAnswers] --> B{Contains 'total' + 'ac'?}
    B -->|Yes| C["AC count check<br>actual == acCount"]
    B -->|No| D{Contains '%'?}
    D -->|Yes| E{Contains 'of total'?}
    E -->|Yes| F["Convert count → %<br>(count × 100) / acCount<br>then compare"]
    E -->|No| G["Direct % compare<br>≥ threshold → PASS<br>within 10% → WARN"]
    D -->|No| H{Starts with ≥ / >=?}
    H -->|Yes| I["actual ≥ threshold → PASS"]
    H -->|No| J{Starts with ≤ / <=?}
    J -->|Yes| K["actual ≤ threshold → PASS"]
    J -->|No| L{Contains '-'?}
    L -->|Yes| M["Range: min ≤ actual ≤ max → PASS<br>off by 1 → WARN"]
    L -->|No| N["Exact match<br>case-insensitive"]
```

Comparison order matters — AC-count and percentage checks run before generic `≥`/`≤` to avoid misparse.

### 4. Report Generation

After all prompts are evaluated (or after the first FAIL in fix mode), `ChecklistReport.CalculateSummary()` computes:

```
TotalPrompts, PassCount, WarnCount, FailCount, SkipCount, PassRate
```

The report is rendered as a table in the terminal.

### 5. Fix Mode (--fix)

```mermaid
flowchart TD
    A[EvaluateUntilFailure] --> B{First FAIL found}
    B --> C[FixPromptGenerator.Generate]
    C --> D{AI response type}
    D -->|clarify_questions| E[Ask user for answers]
    E --> F{Iteration < 5?}
    F -->|Yes| C
    F -->|No| Z[Give up]
    D -->|fix_prompt| G[Display fix prompt]
    G --> H{User choice}
    H -->|Apply| I[FixApplier sends prompt to AI]
    I --> J[AI returns updated artifact]
    J --> K[Save as next version]
    K --> A
    H -->|Refine| L{Refinement < 3?}
    L -->|Yes| M[Collect user feedback]
    M --> N[Regenerate fix prompt]
    N --> G
    L -->|No| G
    H -->|Exit| O[Stop]
```

### 6. On All Checks Passed

**User stories:** The story's `stage` field is advanced to the next stage (e.g., `story_creation` → `refinement`), and the story YAML is saved to `docs/stories/`.

**Tests:** The fixed test file is written back to disk at its `TestFilePath`.

## Build-Time Schema Validation

Both checklist YAML files are validated against `bdd-cli/user-story-description-checklist-schema.yaml` using [Yamale](https://github.com/23andMe/Yamale) during `make lint-docs`:

```makefile
yamale -s bdd-cli/user-story-description-checklist-schema.yaml \
         bdd-cli/user-story-description-checklist.yaml \
         bdd-cli/test-validation-checklist.yaml
```

This ensures structural consistency (correct field names, types, required vs optional) before any runtime use.

## Key Source Files

| File | Role |
|---|---|
| `internal/domain/models/checklist/` | Domain models: Checklist, Stage, Section, Prompt, ValidationResult |
| `internal/infrastructure/checklist/checklist_loader.go` | Loads YAML, extracts stage-specific prompts |
| `internal/app/generators/validate/checklist_evaluator.go` | AI evaluation, answer comparison |
| `internal/app/generators/validate/fix_prompt_generator.go` | Fix prompt generation with clarification loop |
| `internal/app/generators/validate/fix_applier.go` | Applies fix prompts to produce updated artifacts |
| `internal/app/commands/us_validation_command.go` | User story validation command |
| `internal/app/commands/req_validation_command.go` | Test validation command |
