# Expectations for `us apply 99.3 --fix` (re-walk converges)

The fixture seeds `docs/requirements.yaml` so that:

- AC-1 (`99.3-001`) is referenced by TWO entries (`INT-900` and
  `INT-901`) — a deliberate duplicate that forces a collapse fix.
- AC-2 (`99.3-002`) is referenced by NO entry — its scenario must be
  added.

The expected behavior under the refined `--fix` semantics:

1. **Walk #1** —
   - AC-1 "present?" PASS; AC-1 "duplicate?" FAIL → F: collapses
     `INT-900` and `INT-901` into a single entry; re-check passes.
   - AC-2 "present?" FAIL → F: adds a new entry referencing
     `99.3-002`; AC-2 "duplicate?" PASS.
   - Walk #1 ends with every prompt passing but two fixes applied.
2. **Walk #2** — the new re-walk semantics kick in. AC-1 and AC-2 are
   each re-evaluated; both pass on the first prompt of each. Walk #2
   ends with `anyFixApplied=false` → fixpoint reached → canonical
   commit.

## What MUST be true after the run

1. `docs/requirements.yaml` contains **exactly one** scenario whose
   `user_stories[]` references `99.3-001` (the collapse worked).
2. `docs/requirements.yaml` contains **exactly one** scenario whose
   `user_stories[]` references `99.3-002` (AC-2 was added).
3. No two scenarios in `docs/requirements.yaml` share the same
   `merged_steps` (no duplicates remain after the collapse).
4. The pre-existing seed entries `INT-900` and `INT-901` are either:
   - collapsed into a single entry (id may be either, or a new id —
     up to the F: handler), OR
   - one of them is preserved verbatim while the other is removed.
   Either outcome is acceptable so long as rule 1 holds.
5. No files outside `docs/requirements.yaml` and the per-run `tmp/`
   directory are created or modified. In particular, no edits to
   `docs/stories/`, `bdd-cli/`, or `scripts/`.

## What MUST NOT happen

- The string `RE-WALK 3/` MUST NOT appear in stdout. (Walk #2 was the
  clean confirmation walk — it must not itself trigger a third.)
- The string `Hit max apply attempts` MUST NOT appear in stdout. (The
  cap from `config.max_apply_attempts: 5` did not fire.)
- The canonical file MUST NOT be byte-identical to the seed
  `input/docs/requirements.yaml` — at minimum, the AC-2 scenario must
  have been added, and the duplicate must have been collapsed.

## Tolerances

- Order of scenarios inside `docs/requirements.yaml` may differ from
  the seed.
- The collapsed entry's `description`, `service`, `last_updated`, and
  `merged_steps` may be any reasonable merger of `INT-900` and
  `INT-901` — exact formatting is up to the F: handler.
- The new AC-2 entry's id may be any `INT-NNN` or `E2E-NNN` value not
  reused from the seed; the F: handler's rules in `us-apply.yaml`
  govern.
- Files inside `tmp/` are scratch artifacts — ignore any changes
  there.

Reply `PASS` if all the MUST-be-true rules hold AND all the
MUST-NOT-happen rules hold. Otherwise reply
`FAIL: <one-sentence reason>` describing the first violation you find.
