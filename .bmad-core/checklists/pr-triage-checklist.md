<!-- Powered by BMAD™ Core -->

# PR Triage Checklist

- [ ] Tools present: gh, go, jq
- [ ] PR number detected for current branch
- [ ] Conversations JSON fetched (ALL threads) to ./tmp/PR_CONVERSATIONS.json
- [ ] Auto-resolve executed for threads with all comments marked outdated (before analysis)
- [ ] Optional refresh of conversations JSON performed
- [ ] Architecture context read: tech-stack, coding-standards, source-tree, architecture.md, frontend-architecture.md
- [ ] Comprehensive relevance classification applied to remaining threads based on PR scope
- [ ] Human approval requested per relevant thread, processed one-by-one
- [ ] Preferred option presented and confirmation requested: "Default: <Preferred>. Proceed with default?"
- [ ] If declined, accept alternative: 1 Proceed fix / 2 Create ticket / 3 Custom
- [ ] Relevant items: if approved, fix implemented and tests run; thread replied and resolved (no commit yet)
- [ ] Not relevant now: issue created (linked), thread replied and resolved
- [ ] Console summary printed (auto-resolved, fixed, ticketed)
