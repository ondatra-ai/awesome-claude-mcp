# BDD Scenario Writing Guidelines

## Purpose
This document defines BDD best practices for writing Given-When-Then scenarios in user stories, based on industry standards from Automation Panda and Cucumber.

**Last Updated**: 2025-10-03
**References**:
- [BDD 101: Writing Good Gherkin](https://automationpanda.com/2017/01/30/bdd-101-writing-good-gherkin/)
- [Cucumber BDD Documentation](https://cucumber.io/docs/bdd/)

---

## Golden Rules

1. **Write scenarios so non-technical stakeholders understand them**
2. **One Scenario, One Behavior**
3. **Integration and E2E scenarios only** - no unit-level scenarios in BDD

---

## Test Level Definitions

### Integration (INT) Scenarios
**What**: Direct API/protocol testing
**How**: Programmatic requests via Playwright Request API
**No UI Required**: Tests run without browser

**Examples**:
```gherkin
✅ Given: MCP server is running
   When: Test client sends WebSocket upgrade request
   Then: Server responds with connection accepted

✅ Given: Client has active connection
   When: Client sends invalid message format
   Then: Server responds with validation error
```

**Keywords**: Test client, API call, WebSocket message, HTTP request

---

### End-to-End (E2E) Scenarios
**What**: Complete user journey through UI or Claude.ai
**How**: Browser automation via Playwright Page API or AI interaction
**UI Required**: Tests run in browser or chat interface

**Examples**:
```gherkin
✅ Given: User is in Claude.ai chat
   When: User asks Claude to analyze code
   Then: Claude displays analysis results

✅ Given: User is on login page
   When: User enters credentials and clicks login
   Then: User sees dashboard
```

**Keywords**: User clicks, User sees, User enters, Claude responds

---

### MCP E2E Tests with Claude SDK

**What**: Tests where Claude decides which MCP tools to use
**How**: Claude SDK (`@anthropic-ai/sdk`) + MCP Client (`@modelcontextprotocol/sdk`)
**No Browser Required**: Playwright used only for test framework (assertions, structure)

**Architecture**:
```
Test → Claude SDK → Claude decides tool → MCP Client → MCP Server → Google Docs
```

**Examples**:
```gherkin
✅ Given: MCP server runs with document operation tools
   When: Claude API client performs complete workflow from initialize to tool call
   Then: Server processes entire flow and returns valid tool result

✅ Given: MCP server accepts HTTP connections on configured endpoint
   When: Claude API client executes document operation via MCP tools
   Then: Server returns tool result within 2 seconds round-trip time
```

**Keywords**: Claude API client, MCP server, tool result, document operation

**Key Differences from INT tests**:
- INT tests: Direct JSON-RPC calls to MCP server
- E2E tests: Claude decides which tools to call based on user prompt

**See**: [MCP E2E Testing Guide](./mcp-e2e-testing.md) for implementation details.

---

## Given-When-Then Format

### Structure Rules
- **Given**: State/preconditions (present tense, third person)
- **When**: Action/event trigger (present tense, third person)
- **Then**: Expected outcome (present tense, third person)

### Writing Style

✅ **DO**:
- Use declarative style (**WHAT** happens, not **HOW**)
- Use active voice ("Client connects" not "Connection is established")
- Use present tense throughout
- Focus on observable behavior
- Keep steps under 120 characters
- Use third-person perspective
- Write so a Product Owner understands

❌ **DON'T**:
- Use passive voice ("is initialized", "is configured", "is processed")
- Include technical implementation details
- Mention internal components (ConnectionManager, ResponseFormatter, MessageValidator)
- Use vague qualifiers ("properly", "correctly", "specific")
- Combine multiple behaviors in one scenario
- Write scenarios about unit-level code
- Use first-person ("I click", "we send")

---

## The "Product Owner Test"

Before accepting any scenario, ask these four questions:

1. **Would a Product Owner understand this without asking questions?**
   If NO → Reject scenario

2. **Does this describe value delivered to a user/system?**
   If NO → Reject scenario

3. **Can this be observed from outside the system?**
   If NO → Reject scenario

4. **Does this mention internal components?**
   If YES → Reject scenario

---

## Examples

### ❌ BAD: Technical, Unit-level
```yaml
given: "Message validator is initialized with MCP schema"
when: "Invalid message format is received"
then: "Validation error is returned with specific violation details"
level: "unit"
```

**Problems**:
- Mentions internal component (Message validator)
- Passive voice ("is initialized", "is received")
- Unit-level (not appropriate for BDD)
- Technical implementation details

---

### ✅ GOOD: Integration, Observable
```yaml
given: "Client has active WebSocket connection"
when: "Client sends message with invalid format"
then: "Server responds with validation error"
level: "integration"
```

**Why Good**:
- Observable external behavior
- No internal components mentioned
- Active voice
- Declarative style
- Product Owner understandable

---

### ❌ BAD: Multiple Behaviors
```yaml
when: "Multiple clients establish connections and exchange MCP messages"
then: "All messages are properly handled, validated, and responded to while maintaining connection stability"
```

**Problems**:
- Multiple behaviors (connection + message handling + stability)
- Vague qualifiers ("properly", "maintaining")
- Cannot clearly determine pass/fail

---

### ✅ GOOD: Single Behavior
```yaml
when: "Multiple clients connect simultaneously"
then: "Server accepts connections up to configured limit"
```

**Why Good**:
- Single testable behavior
- Clear, measurable outcome
- No vague qualifiers

---

## Forbidden Terms in Scenarios

### ❌ Component Names
- ConnectionManager, ResponseFormatter, MessageValidator
- DocumentService, AuthHandler, TokenValidator
- Any class or internal component name

### ❌ Implementation Terms
- initialize, instantiate, configure, register
- parse, serialize, deserialize
- allocate, cleanup, pool

### ❌ Internal Architecture
- middleware, adapter, wrapper, handler
- service, repository, factory
- thread, goroutine, process

### ❌ Vague Qualifiers
- properly, correctly, specific
- appropriate, suitable, valid (without criteria)

---

## ✅ Allowed Terms

### Actors
- Client, User, System, Server
- Administrator, External System
- Claude (for Claude.ai interactions)
- Claude API client (for MCP E2E tests via SDK)

### Actions
- connect, send, receive, respond, reject
- display, show, navigate, click
- authenticate, authorize

### Artifacts
- connection, message, request, response
- page, button, form, dialog
- error, result, data

---

## Validation Checklist

Every scenario must pass ALL these checks:

### Structure
☐ Has exactly one Given-When-Then sequence
☐ No multiple When-Then pairs
☐ Follows strict Given → When → Then order
☐ Tests ONE behavior only

### Language
☐ Active voice (no "is initialized", "is configured")
☐ Third-person perspective maintained
☐ Present tense throughout
☐ No passive constructions

### Content
☐ Declarative (WHAT not HOW)
☐ No component names
☐ No technical implementation details
☐ No vague qualifiers

### Quality
☐ Product Owner understandable
☐ Observable from outside system
☐ Each step < 120 characters
☐ Could be automated without seeing code

---

## Scenario ID Format

**Format**: `{story_id}-{LEVEL}-{sequence}`

**Examples**:
- `3.1-INT-001` - First integration scenario for story 3.1
- `3.1-INT-002` - Second integration scenario
- `3.1-E2E-001` - First E2E scenario for story 3.1

**Invalid**:
- `3.1-UNIT-001` ❌ (No unit scenarios in BDD)
- `3.1-001` ❌ (Missing level)

---

## Priority Assignment

**P0 (Critical)**: Security, data integrity, compliance, revenue-critical
**P1 (High)**: Core user journeys, frequently used features
**P2 (Medium)**: Secondary features, admin functions
**P3 (Low)**: Nice-to-have, rarely used features

---

## Common Mistakes to Avoid

### Mistake 1: Unit-Level Scenarios
❌ Writing scenarios for internal component behavior
✅ Write scenarios for external system behavior

### Mistake 2: Passive Voice
❌ "Given: Connection is established"
✅ "Given: Client establishes connection"

### Mistake 3: Implementation Details
❌ "When: Handler processes message through middleware"
✅ "When: Server processes client message"

### Mistake 4: Multiple Behaviors
❌ "Then: Message validated, processed, and response sent"
✅ Create 3 separate scenarios

### Mistake 5: Vague Outcomes
❌ "Then: System behaves correctly"
✅ "Then: Server responds with status 200"

---

## References

- [Automation Panda - BDD 101: Writing Good Gherkin](https://automationpanda.com/2017/01/30/bdd-101-writing-good-gherkin/)
- [Cucumber - BDD Documentation](https://cucumber.io/docs/bdd/)
- [Coding Standards - Testing](./coding-standards.md#testing-standards)
- [Tech Stack - Testing Strategy](./tech-stack.md#testing-strategy)
- [MCP E2E Testing with Claude SDK](./mcp-e2e-testing.md)
