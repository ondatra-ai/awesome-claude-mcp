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

```
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

### 1. Load and Filter

```
ChecklistLoader.Load()
  → reads YAML file, unmarshals into Checklist struct

ChecklistLoader.ExtractPromptsForStage(checklist, stageID)
  → finds the stage by ID
  → iterates all sections → all validation_prompts
  → skips prompts where skip field is non-empty
  → wraps each prompt in PromptWithContext (adds stage/section metadata, default_docs)
  → returns []PromptWithContext
```

### 2. Evaluate Each Prompt

For each prompt, `ChecklistEvaluator.evaluatePrompt()` runs:

```
1. Resolve document keys
   prompt-level docs (or default_docs) → config lookup → file paths
   e.g. "architecture" → config["documents.architecture"] → "docs/architecture.md"

2. Render templates
   System prompt: loaded from templates.prompts.checklist_system (cached)
   User prompt:   loaded from templates.prompts.checklist
     Injected data: { Subject, SubjectID, Question, Rationale, ResultPath, Docs, FixTemplate }

3. Call Claude (think mode)
   ExecutePromptWithSystem(systemPrompt, userPrompt, mode=think)
   → Claude analyzes the artifact and writes a YAML result between FILE_START/FILE_END markers

4. Parse response
   Extract content between FILE_START <path> and FILE_END markers
   Parse YAML: { answer: "...", fix_prompt: "..." (optional) }

5. Compare answer
   compareAnswers(expected, actual, acCount) → PASS | WARN | FAIL
```

All prompts, responses, and parsed results are saved to a temp directory for debugging.

### 3. Answer Comparison Rules

The `compareAnswers` function handles multiple expected-answer patterns:

| Expected pattern | Example | Logic |
|---|---|---|
| Exact string | `"yes"`, `"no"`, `"need"` | Case-insensitive `expected == actual` |
| Range | `"3-7"` | `min ≤ actual ≤ max` → PASS; off by 1 → WARN |
| Greater-or-equal | `"≥2"` or `">=2"` | `actual ≥ threshold` → PASS |
| Less-or-equal | `"≤10"` or `"<=10"` | `actual ≤ threshold` → PASS |
| AC count equality | `"= total AC count"` | `actual == len(acceptanceCriteria)` → PASS |
| Percentage | `"≥50%"` | `actual ≥ threshold` → PASS; within 10% → WARN |
| Percentage of total | `"≥50% of total"` | Converts actual count to `(count × 100) / acCount`, then compares |

Comparison order matters — AC-count and percentage checks run before generic `≥`/`≤` to avoid misparse.

### 4. Report Generation

After all prompts are evaluated (or after the first FAIL in fix mode), `ChecklistReport.CalculateSummary()` computes:

```
TotalPrompts, PassCount, WarnCount, FailCount, SkipCount, PassRate
```

The report is rendered as a table in the terminal.

### 5. Fix Mode (--fix)

When `--fix` is passed, the system enters an interactive loop:

```
┌─────────────────────────────────────────────────┐
│  1. EvaluateUntilFailure                        │
│     → stops at first FAIL                       │
│                                                 │
│  2. FixPromptGenerator.Generate                 │
│     → sends failed check + fix template to AI   │
│     → AI returns EITHER:                        │
│       a) fix_prompt (actionable instructions)   │
│       b) clarify_questions (needs user input)   │
│                                                 │
│  3. If questions → ask user, loop (max 5 iters) │
│                                                 │
│  4. Display fix prompt to user                  │
│     User chooses: [A]pply / [R]efine / [E]xit   │
│                                                 │
│  5. If Apply:                                   │
│     → FixApplier sends fix prompt to AI         │
│     → AI returns updated artifact               │
│     → Save as next version                      │
│     → Re-run validation from step 1             │
│                                                 │
│  6. If Refine (max 3 iterations):               │
│     → Collect user feedback                     │
│     → Regenerate fix prompt                     │
│     → Back to step 4                            │
│                                                 │
│  7. If Exit → stop                              │
└─────────────────────────────────────────────────┘
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
| `scripts/bmad-cli/internal/domain/models/checklist/` | Domain models: Checklist, Stage, Section, Prompt, ValidationResult |
| `scripts/bmad-cli/internal/infrastructure/checklist/checklist_loader.go` | Loads YAML, extracts stage-specific prompts |
| `scripts/bmad-cli/internal/app/generators/validate/checklist_evaluator.go` | AI evaluation, answer comparison |
| `scripts/bmad-cli/internal/app/generators/validate/fix_prompt_generator.go` | Fix prompt generation with clarification loop |
| `scripts/bmad-cli/internal/app/generators/validate/fix_applier.go` | Applies fix prompts to produce updated artifacts |
| `scripts/bmad-cli/internal/app/commands/us_validation_command.go` | User story validation command |
| `scripts/bmad-cli/internal/app/commands/req_validation_command.go` | Test validation command |
