<!-- Powered by BMAD™ Core -->

# All Conversations for PR #{{pr_number}}:

## ❌ OUTDATED (Fixed by previous changes):

{{#each outdated}}
### **{{file}}:{{line}}**
Id: {{id}}
Author: {{author}}
Description: {{description}}
----
{{body}}
----
Status: OUTDATED: All comments marked outdated or code no longer exists.

{{/each}}

## ✅ STILL RELEVANT (Need to be fixed):

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
