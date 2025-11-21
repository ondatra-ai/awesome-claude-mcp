<!-- Powered by BMAD Core -->

# Validate Test Quality

## Your Task

Validate all test files against best practices and fix any issues found.

---

## Step 0: Understand the Validation Criteria

**CRITICAL**: Tests must fail fast and explicitly. Any hidden problem must be immediately surfaced.

### Anti-Patterns to Detect and Fix

1. **Conditional Logic That Hides Failures**
   ```typescript
   // BAD: Hides failures with silent conditionals
   if (!handshakeComplete && message.id === 1 && message.result) {
     handshakeComplete = true;
     return;  // Silently skips if condition fails
   }
   ```

2. **Missing Error Checks**
   ```typescript
   // BAD: Doesn't check for errors
   if (result.result) {
     // Process result - but what if result.error exists?
   }
   ```

3. **Permissive Status Acceptance**
   ```typescript
   // BAD: Accepts both success and error
   expect([200, 400]).toContain(response.status());
   ```

4. **Complex State Machines Without Assertions**
   ```typescript
   // BAD: State machine that fails silently
   const result = { success: false, handshakeComplete: false };
   // ... 100 lines of complex logic ...
   expect(result.success).toBeTruthy(); // Unclear what failed
   ```

### Correct Patterns to Enforce

1. **Explicit Error Checking**
   ```typescript
   // GOOD: Check for errors explicitly
   expect(message.error, 'Server should not return error').toBeUndefined();
   expect(message.result, 'Response must have result field').toBeDefined();
   ```

2. **Custom Messages for Context**
   ```typescript
   // GOOD: Clear context in assertion
   expect(tools.find(t => t.name === 'replace_all'),
     'replace_all tool should be in tools list').toBeDefined();
   ```

3. **Strict Status Checks**
   ```typescript
   // GOOD: Strict success assertion
   expect(response.status(), 'Expected successful status').toBe(200);
   ```

4. **Immediate Assertions After Operations**
   ```typescript
   // GOOD: Assert immediately after each operation
   const response = await request.get(url);
   expect(response.status(), 'Initialize request should succeed').toBe(200);
   const data = await response.json();
   expect(data.connectionId, 'Response should include connectionId').toBeDefined();
   ```

---

## Step 1: Scan Test Files

Use the Glob tool to find all test files:
```
Glob pattern: {{.TestFilesGlob}}
```

---

## Step 2: Analyze Each Test File

For each test file found:

1. **Read the file content**
2. **Identify tests that expect SUCCESS** (not error-testing scenarios)
3. **Check for anti-patterns**:
   - `if (result.result)` without checking `result.error`
   - `if (!stateFlag && condition)` patterns that silently skip
   - Permissive status acceptance like `[200, 400].includes(status)`
   - Missing custom messages in expect() calls
   - Complex state tracking without intermediate assertions

---

## Step 3: Fix Identified Issues

For each issue found, use the Edit tool to:

1. **Add explicit error checks** before processing results:
   ```typescript
   expect(result.error, 'Server should not return error').toBeUndefined();
   expect(result.result, 'Server must return result object').toBeDefined();
   ```

2. **Add custom messages** to all assertions:
   ```typescript
   expect(response.status(), 'Request should succeed with 200').toBe(200);
   ```

3. **Remove permissive patterns** that accept both success and failure:
   ```typescript
   // Change this:
   if (result.result) { ... } else if (result.error) { ... }

   // To this:
   expect(result.error, 'Should not have error').toBeUndefined();
   expect(result.result, 'Must have result').toBeDefined();
   // Then process result
   ```

---

## Step 4: Validate Error-Testing Scenarios

**IMPORTANT**: Some tests legitimately expect errors. Do NOT fix these:

- Tests with descriptions containing "error", "invalid", "reject", "fail"
- Tests that explicitly check `expect(result).toHaveProperty('error')`
- Tests for negative scenarios (e.g., INT-047, INT-048, INT-051, INT-052)

These tests should:
- Still have explicit assertions
- Check for expected error codes/messages
- NOT silently accept both success and error

---

## Step 5: Save Validation Results

**CRITICAL**: You MUST save the validation results to a YAML file using the Write tool.

Use the Write tool to create: `{{.TmpDir}}/validate-tests-result.yaml`

### Result Schema

```yaml
# If ALL issues were fixed or no issues found:
result: ok
data:
  files_scanned: <number>
  issues_found: <number>
  issues_fixed: <number>
  unfixed_issues: []

# If ANY issues could NOT be automatically fixed:
result: fail
data:
  files_scanned: <number>
  issues_found: <number>
  issues_fixed: <number>
  unfixed_issues:
    - file: "path/to/file.spec.ts"
      line: 42
      description: "Brief description of the issue"
      suggested_fix: "How to fix it manually"
```

### Example Output Files

**Success case (all issues fixed):**
```yaml
result: ok
data:
  files_scanned: 5
  issues_found: 3
  issues_fixed: 3
  unfixed_issues: []
```

**Failure case (some issues remain):**
```yaml
result: fail
data:
  files_scanned: 5
  issues_found: 4
  issues_fixed: 2
  unfixed_issues:
    - file: "tests/integration/mcp-protocol.spec.ts"
      line: 156
      description: "Complex state machine with multiple branches requires manual refactoring"
      suggested_fix: "Split into separate test cases with explicit state assertions"
    - file: "tests/e2e/editor-flow.spec.ts"
      line: 89
      description: "Dynamic timeout value cannot be statically validated"
      suggested_fix: "Add explicit timeout assertion with custom message"
```

---

## Completion Signal

After writing the YAML result file, respond with:
```
VALIDATE_TESTS_COMPLETE
```

Do not add any explanations after this signal.
