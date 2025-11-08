# Implement Feature - Attempt {{.Attempt}} of {{.MaxAttempts}}

Your task is to make all failing tests pass for user story **{{.StoryID}}** by implementing production code.

## Story Context

**Story**: {{.StoryID}} - {{.StoryTitle}}

**User Story**:
- As a: {{.AsA}}
- I want: {{.IWant}}
- So that: {{.SoThat}}

## Current Test Results

I've run the tests and they are **FAILING**. Here's the output:

```
{{.TestOutput}}
```

## Your Task

**Fix the production code** to make these tests pass. Do NOT modify the tests themselves.

### Steps:

1. **Analyze the test failures above** - understand what's broken

2. **Read the failing test files** to understand exact requirements

3. **Implement/fix the production code**:
   - Follow the project's coding standards
   - Create necessary files in the correct locations
   - Implement complete functionality (no stubs or TODOs)
   - **Fix code, NOT tests**

4. **Verify your changes** - you can run tests again if needed:
   ```bash
   {{.TestCommand}}
   ```

## Important Rules

- ✅ Analyze the test failures shown above
- ✅ Fix production code to make tests pass
- ✅ Implement complete, working code
- ❌ Do NOT modify tests (unless genuinely wrong)
- ❌ Do NOT leave TODO comments or placeholder code

## Success Criteria

- All tests pass when running `{{.TestCommand}}`
- Code follows project coding standards
- Tests were NOT modified (only production code was changed)

**This is attempt {{.Attempt}} of {{.MaxAttempts}}**. After you finish, tests will run again automatically.
