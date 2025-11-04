# Make Failing Tests Pass

Your task is to implement the production code needed to make all failing tests for user story **{{.StoryID}}** pass.

## Story Context

**Story**: {{.StoryID}} - {{.StoryTitle}}

**User Story**:
- As a: {{.AsA}}
- I want: {{.IWant}}
- So that: {{.SoThat}}

## Test Files to Fix

The following test files have been created and are currently failing:

{{range .TestFiles}}
- `{{.}}`
{{end}}

## Your Approach

1. **Run the tests first** to see what's failing:
   ```bash
   {{.TestCommand}}
   ```

2. **Read the failing tests** to understand what needs to be implemented

3. **Implement the production code** following these guidelines:
   - Follow the project's coding standards
   - Reference the story tasks for implementation details
   - Create necessary files in the correct locations
   - Implement complete functionality (no stubs or TODOs)

4. **Run tests again** to verify they pass

5. **Iterate** until all tests pass - keep running tests and fixing issues until everything is green

## Success Criteria

- All tests in the test files listed above pass
- Code follows project coding standards
- No TODO comments or placeholder code remains

Begin by running the tests to see what needs to be implemented.
