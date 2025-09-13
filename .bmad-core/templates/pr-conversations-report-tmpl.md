<!-- Powered by BMAD™ Core -->

# All Conversations for PR #{{pr_number}}:

## ❌ Auto-Resolved Outdated:

{{#each outdated}}
### **{{file}}:{{line}}**
Id: {{id}}
Author: {{author}}
Description: {{description}}
----
{{body}}
----
Status: OUTDATED: All comments were marked outdated and resolved automatically.

{{/each}}

## ✅ Still Relevant After Auto-Resolve:

{{#each relevant}}
### **{{file}}:{{line}}**
Id: {{id}}
Author: {{author}}
Description: {{description}}
----
{{body}}
----
Status: RELEVANT: At least one comment not marked outdated.
Recommendation: Review intent and address accordingly.
Decision:

{{/each}}
