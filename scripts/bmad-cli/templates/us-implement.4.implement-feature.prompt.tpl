# Implement Feature - Attempt {{.Attempt}} of {{.MaxAttempts}}

Your task is to make all failing tests pass for user story **{{.StoryID}}** by implementing production code.

## Story Context

**Story**: {{.StoryID}} - {{.StoryTitle}}

**User Story**:
- As a: {{.AsA}}
- I want: {{.IWant}}
- So that: {{.SoThat}}

## Architecture Context

**CRITICAL**: Before implementing, understand the project architecture to place code correctly:

1. Read for references the following documents:
  - Read(`docs/architecture.md`) - Three-service architecture overview
  - Read(`docs/architecture/source-tree.md`) - Detailed service structure and file locations
  - Read(`docs/architecture/coding-standards.md`) - Coding standards and conventions
  - Read(`docs/architecture/tech-stack.md`) - Technology stack per service

2. **Understand Service Structure**:
   - **Frontend service**: `services/frontend/` - Next.js React application
   - **Backend API service**: `services/backend/` - Go REST API for user management
   - **MCP Protocol service**: `services/mcp-service/` - Go MCP WebSocket server

3. **Place Code in Correct Service**:
   - MCP WebSocket handlers → `services/mcp-service/internal/server/`
   - REST API endpoints → `services/backend/internal/`
   - React UI components → `services/frontend/`
   - **DO NOT put MCP code in backend service** - they are separate services

4. **Verify Against Story**:
   - Read the story file (`docs/stories/{{.StoryID}}-*.yaml`) to see specified file paths in tasks
   - Cross-reference task subtasks with architecture to ensure alignment

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
