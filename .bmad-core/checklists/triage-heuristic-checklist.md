<!-- Powered by BMAD™ Core -->

# Triage Heuristic Checklist

## Heuristic Session

- [ ] Locate and review the exact code under discussion
  - [ ] Read ~30 lines around the change for context; scan related functions/modules
  - [ ] If applicable, run the project or tests to reproduce the behavior

- [ ] Read all
  - [ ] Summarize the suggested change(s) and intended outcome in one sentence
  - [ ] Note assumptions, open questions, and edge cases mentioned by reviewers

- [ ] Verify whether the reported issue is already fixed in the current code

- [ ] Verify alignment with project standards (no conflicts)
  - [ ] System architecture: docs/architecture.md
  - [ ] Source tree and locations: docs/architecture/source-tree.md
  - [ ] Coding standards and patterns: docs/architecture/coding-standards.md
  - [ ] Approved tech/versions: docs/architecture/tech-stack.md
  - [ ] Record outcome: OK / CONFLICTS FOUND (list conflicts if any)

- [ ] Analyze pros and cons of the proposed change
  - [ ] Benefits (e.g., correctness, performance, security, maintainability, DX)
  - [ ] Costs/risks (e.g., complexity, regressions, technical debt, performance/security impact)
  - [ ] Scope fit: is this within current story/PR objectives?

- [ ] Research and propose a better or confirming solution (if needed)
  - [ ] Compare 1–2 viable alternatives and select a recommendation
  - [ ] Validate the recommendation against the four documents above
  - [ ] Provide a brief code-level pointer (file/function) or sketch if helpful

- [ ] Decide whether the change can be postponed (all of the following should be true)
  - [ ] It is not a blocker
  - [ ] It is out of scope
  - [ ] It has low impact/ROI
  - [ ] It requires broader refactoring that touches more than one file
