# Unified Checklist Engine

## Core insight

Every CLI command is a walk over **entities × prompts**, with `--fix` turning
each failed (entity, prompt) cell into a Claude-driven fix step. There are
exactly two things that vary per command:

1. **Entity source** — what populates the entity list.
2. **Fix target** — what the fix step mutates.

Everything else (prompt loading, evaluation, fix-prompt generation, fix
application, version management, reporting) is shared.

## The matrix

| Command            | Entities                                     | Prompts (N)                  | Fix target                       |
|--------------------|----------------------------------------------|------------------------------|----------------------------------|
| `us create`        | 1 — story extracted from epic                | from `us-create.yaml`        | story file in `docs/stories/`    |
| `us refine`        | 1 — story loaded from `docs/stories/`        | from `us-refine.yaml`        | same story file                  |
| `us apply`         | M — ACs from one story file                  | from `us-apply.yaml`         | `docs/requirements.yaml`         |

The `1 × N` commands share `newUSChecklistCmd`. The `M × N` commands share
the `ExecuteStoryScenarioChecklist` walker for story-sourced scenarios. The
walker delegates to the same evaluator / fix-generator / fix-applier triple.

## The walk (pseudocode)

    entities = parser.parse(source)            // story file or requirements.yaml
    prompts  = checklistLoader.load(name)      // bdd-cli/checklists/<name>.yaml

    for entity in entities:
      for prompt in prompts:
        result = evaluator.evaluate(entity, prompt)
        if result == FAIL and --fix:
          fixPrompt = fixGen.generate(entity, prompt)
          fixApplier.apply(entity, fixPrompt)        // mutates fix target
          re-evaluate

For `1 × N` commands `entities` is a single-element list; the inner loop is
unchanged. The `--fix` branch is what carries the per-command difference: a
`us refine` fix rewrites the story; a `us apply` fix calls `Edit` on a scratch
copy of `docs/requirements.yaml`.

## Pluggable pieces

| Piece              | Subject-agnostic? | Where it lives                                                      |
|--------------------|-------------------|---------------------------------------------------------------------|
| ChecklistLoader    | yes               | `internal/infrastructure/checklist/checklist_loader.go`             |
| Evaluator          | yes (template-driven) | `internal/app/generators/validate/checklist_evaluator.go`         |
| FixPromptGenerator | yes (template-driven) | `internal/app/generators/validate/fix_prompt_generator.go`        |
| FixApplier         | yes; returns content, caller persists | `internal/app/generators/validate/fix_applier.go`           |
| EntityParser       | NO — one impl per source | story: `infrastructure/story/story_scenario_parser.go` |

Adding a new command is therefore: new checklist YAML + (optional) new entity
parser + (optional) new template set + thin command wiring. No new engine
code.

## Persistence note for `us apply`

`us apply` mutates a single shared file (`docs/requirements.yaml`). To avoid
leaving the registry in a partial state if the run aborts mid-walk:

1. At start, copy `docs/requirements.yaml` → `<tmpDir>/requirements.yaml`
   (scratch).
2. Every evaluator and fix-applier prompt reads from and writes to the
   scratch path; `docs/requirements.yaml` is untouched.
3. After the walk completes successfully, `os.Rename(scratch, original)`
   atomically.
4. On any abort, discard the scratch.

Because the walk is sequential (M scenarios × N prompts, one cell at a time)
there is no concurrency. Each cell sees the cumulative effect of all prior
fixes through the scratch file.

## Why `us apply` is M × N, not 1 × N

`us create` and `us refine` evaluate properties of a story *as a whole*
(role/goal/AC count/etc.) — those questions have one answer per story, so the
walk is `1 × N`.

`us apply` evaluates properties of *each AC inside* the story (does this
scenario already live in `requirements.yaml`? are there duplicates?) — those
are per-scenario questions, so the walk is `M × N`. The fix prompt for each
(scenario, prompt) cell is literally what performs the "copy" into
`requirements.yaml`.

## Replacing the old `us merge-scenarios`

The deleted `merge_scenarios_generator.go` (last in commit `1001bd5`) ran a
one-shot Claude call per scenario with no validation loop. The checklist
engine subsumes it: the validate prompts answer "is this already merged
correctly?", the fix prompt does the merge, and the loop gives idempotency
(`us apply` without `--fix` is a dry-run report; with `--fix` it actually
merges). No new generator code; no separate one-shot path.
