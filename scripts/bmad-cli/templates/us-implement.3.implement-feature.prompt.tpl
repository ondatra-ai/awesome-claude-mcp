# Implement Feature

Your task is to make all failing tests pass for user story **{{.StoryID}}** by implementing production code.

## Story Context

**Story**: {{.StoryID}} - {{.StoryTitle}}

**User Story**:
- As a: {{.AsA}}
- I want: {{.IWant}}
- So that: {{.SoThat}}

## Your Task

1. **Run the tests** to see what's failing:
   ```bash
   {{.TestCommand}}
   ```

2. **Analyze the failures** - understand what production code needs to be implemented

3. **Read the failing test files** to understand the exact requirements

4. **Implement the production code** following these guidelines:
   - Follow the project's coding standards
   - Reference the story tasks for implementation details
   - Create necessary files in the correct locations
   - Implement complete functionality (no stubs or TODOs)
   - **Fix code, NOT tests** - tests define the requirements

5. **Run tests again** to verify they pass

6. **Iterate** until all tests pass - keep running tests and fixing code until everything is green

## Important Rules

- ✅ Run tests first to see failures
- ✅ Fix production code to make tests pass
- ✅ Run tests frequently to get feedback
- ✅ Iterate until all tests are green
- ❌ Do NOT modify tests (unless genuinely wrong)
- ❌ Do NOT leave TODO comments or placeholder code

## Success Criteria

- All tests pass when running `{{.TestCommand}}`
- Code follows project coding standards
- No TODO comments or placeholder code remains
- Tests were NOT modified (only production code was changed)

Begin by running the tests to see what needs to be fixed.
