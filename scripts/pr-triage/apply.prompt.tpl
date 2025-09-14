# Apply Changes Prompt Template

As @dev implement suggested changes.

Context:
- PR Number: {{PR_NUMBER}}
- Conversation Location: {{LOCATION}}
- Conversation URL: {{URL}}
- Conversation Text:
{{CONVERSATION_TEXT}}

Guidance:
- Prefer minimal, localized changes within the specified file/location.
- Keep changes aligned with architecture and coding standards (naming, formatting, module boundaries).
- If critical uncertainty exists (e.g., naming conventions, missing context), set action to "clarify-then-implement" and craft a precise question in reply_comment; do not produce a patch.
- If the change is beneficial but out-of-scope or cross-cutting, set action to "create-followup-ticket" and propose a ticket title in reply_comment; do not produce a patch.
- If you detect architectural/security/performance implications beyond scope, set action to "escalate-architecture" and explain briefly in reply_comment; do not produce a patch.
- Only include tests_to_run that are relevant and fast (e.g., file-specific unit tests, linters).

Output example (format guide, not content; DO NOT copy verbatim):
"""
  diff --git a/infrastructure/terraform/environments/variables.tf b/infrastructure/terraform/environments/variables.tf
  index 0000000..0000000 100644
  --- a/infrastructure/terraform/environments/variables.tf
  +++ b/infrastructure/terraform/environments/variables.tf
  @@ -6,6 +6,12 @@ variable "environment" {
    description = "Environment name (dev|staging|prod)"
    type        = string
  +  validation {
  +    condition     = contains(["dev", "staging", "prod"], var.environment)
  +    error_message = "environment must be one of: dev, staging, prod."
  +  }
   }
"""
