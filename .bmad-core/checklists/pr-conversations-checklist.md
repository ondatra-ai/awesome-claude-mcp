<!-- Powered by BMAD™ Core -->

# PR Conversations Read Checklist

- [ ] Tools present: gh, go, jq
- [ ] PR number detected for current branch
- [ ] Conversations JSON fetched (ALL threads) to ./tmp/PR_CONVERSATIONS.json
- [ ] Auto-resolve executed for threads with all comments marked outdated (before analysis)
- [ ] Optional refresh of conversations JSON performed
- [ ] Comprehensive relevance classification applied to the remaining threads
- [ ] Report generated at ./tmp/PR_CONVERSATIONS.md following template (sections: Auto-Resolved Outdated, Still Relevant After Auto-Resolve)
- [ ] Outdated and relevant counts printed
