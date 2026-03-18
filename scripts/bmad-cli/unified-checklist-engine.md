# Unified Checklist Engine

## Core Insight

Every CLI command can be implemented as: **checklist + prompts + entity filter + fix mode**.

The existing checklist engine (ChecklistEvaluator → FixPromptGenerator → FixApplier) already handles both single-entity (`us create`) and multi-entity (`req validate`) cases. The only command that bypasses this engine is `us merge-scenarios`, which uses a separate one-shot generator. This proposal unifies all commands under the checklist engine.

## The Model

Every command is defined by three pieces of config:

1. **Checklist YAML** — what to validate (validation prompts with expected answers)
2. **System + user prompt** — how to evaluate and how to fix
3. **Entity filter** — which entities to operate on (JSON filter / parser)

The `--fix` flag controls behavior:
- Without `--fix`: report-only mode (dry run showing pass/fail status)
- With `--fix`: validate → fail → fix loop (the fix prompt does the actual work)

## Unified Runner (Pseudocode)

```
entities = load_entities(entity_filter, source_file)
prompts  = load_checklist(checklist_yaml, stage)

for each entity in entities:
    for each prompt in prompts:
        result = evaluate(entity, prompt)
        if FAIL and --fix:
            generate fix_prompt(entity, failed_check)
            apply fix
            re-evaluate
```

This is exactly what `us create` and `req validate` already do. The difference is only in what `entities` and `prompts` are.

## How Existing Commands Map

| Command | Checklist | Entities | Fix action |
|---------|-----------|----------|------------|
| `us create 4.1 --fix` | user-story-description-checklist.yaml (stage: story_creation) | `[story 4.1]` (single entity from epic) | Rewrite story text |
| `us refine 4.1 --fix` | user-story-description-checklist.yaml (stage: refinement) | `[story 4.1]` (single entity from story file) | Improve story text |
| `req validate --fix` | test-validation-checklist.yaml (stage: test_validation) | All scenarios from requirements.yaml | Rewrite test code |
| `us merge-scenarios 4.1` | **NEW** merge-scenarios-checklist.yaml | Scenarios from story 4.1 | Merge scenario into requirements.yaml |

## Merge Scenarios as a Checklist (New)

Currently `merge_scenarios_generator.go` calls Claude directly per scenario with no validation loop. Reframed as a checklist:

- **Validate prompt**: "Is scenario X already correctly merged into requirements.yaml with proper requirement ID, bidirectional mapping, and story lineage?"
- **Expected answer**: "Yes"
- **Fix prompt**: "Merge this scenario into requirements.yaml following the merge rules (ID assignment, field mapping, conflict resolution)"

This turns a one-shot action into the same validate→fix loop. Benefits:
- `us merge-scenarios 4.1` (no --fix) = dry run showing which scenarios are/aren't merged
- `us merge-scenarios 4.1 --fix` = actually merges them
- Same engine, same reporting, same interactive refinement

## What Needs to Change

1. **Create merge-scenarios-checklist.yaml** with validation prompts for merge status
2. **Create merge-scenarios system/user prompt templates** for the fix action
3. **Create an entity parser** that extracts scenarios from a story (analogous to ScenarioParser for requirements.yaml)
4. **Remove `merge_scenarios_generator.go`** — replaced by the checklist engine
5. **Update `us_merge_scenarios_command.go`** to use ChecklistEvaluator + FixPromptGenerator + FixApplier instead of calling MergeScenariosGenerator directly

## Longer-Term Vision

Any new command becomes just config:
- Write a checklist YAML
- Write system/user prompt templates
- Define the entity filter
- Wire it up as a thin command that calls the unified engine

No new generator code needed. The engine handles the validate/fix loop, versioning, interactive refinement, and reporting for all commands.
