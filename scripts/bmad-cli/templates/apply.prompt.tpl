# Apply Changes Prompt Template

As @dev, implement the suggested changes directly in the repository.

Context:
- PR Number: {{PR_NUMBER}}
- Conversation Location: {{LOCATION}}
- Conversation URL: {{URL}}
- Conversation Text:
{{CONVERSATION_TEXT}}

Instructions (Apply Mode):
- Edit files in this workspace to implement the requested change.
- Make the smallest, standards-aligned modification at the specified location.
- Do not output diffs or plans; apply the change directly.
- After applying, print exactly one short line summarizing what you changed.
- Do not include code fences, YAML, or multi-line output.

Notes:
- Keep changes localized to the referenced file unless strictly necessary.
- Preserve formatting and surrounding context.
