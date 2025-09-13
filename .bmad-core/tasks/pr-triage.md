1) Read conversations to `tmp/CONV.json`.
   - Run `go run scripts/list-pr-conversations/main.go` (no arguments needed). It MUST write the full JSON array of conversations directly to `tmp/CONV.json` and MUST NOT print the JSON to standard output.

2) Create `tmp/CONV_ID.txt`.
   - This file stores processed conversation IDs (one per line). It is empty by default.

3) Identify the next conversation to process.
   - Write a Go script that reads `tmp/CONV.json` and `tmp/CONV_ID.txt`, finds the first conversation ID present in JSON that is not listed in `tmp/CONV_ID.txt`, and writes the complete conversation object to `tmp/CONV_CURRENT.json`. If none remain, write `{ "id": "No More Converations" }` to `tmp/CONV_CURRENT.json`.

4) Read `tmp/CONV_CURRENT.json`.
   - If it contains `{ "id": "No More Converations" }`, stop. Otherwise, proceed with heuristic analysis for that conversation.
