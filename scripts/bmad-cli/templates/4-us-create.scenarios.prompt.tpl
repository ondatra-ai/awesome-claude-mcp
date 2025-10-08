<!-- Powered by BMAD™ Core -->

# Create Test Scenarios with Given-When-Then Format

## Purpose

Generate comprehensive test scenarios in Given-When-Then format for user story {{.Story.ID}}.

## Instructions

**CRITICAL**: Generate ONLY Integration (INT) and End-to-End (E2E) scenarios. NEVER generate unit-level scenarios.

---

### Step 1: Review Acceptance Criteria

For EACH acceptance criterion ({{range $i, $ac := .Story.AcceptanceCriteria}}{{if $i}}, {{end}}{{$ac.ID}}{{end}}), identify what system behavior proves it's implemented.

**Key Question**: "What can an external actor observe that proves this criterion?"

**Example**:
- AC-1: "WebSocket server implemented"
- ❌ Unit thinking: Test WebSocket class initialization
- ✅ System thinking: Test client successfully connects via WebSocket

---

### Step 2: Determine Test Level (INT vs E2E)

**Decision Question**: "Does testing this require UI or Claude.ai interaction?"

**NO** → Integration (INT)
- Test via Playwright Request API
- Direct protocol/API testing
- No browser needed
- Examples: HTTP requests, WebSocket messages, API responses

**YES** → End-to-End (E2E)
- Test via Playwright Browser API
- UI workflows or Claude chat
- Browser required
- Examples: User clicks button, sees result, completes auth flow

---

### Step 3: Apply BDD Best Practices

**Given-When-Then Structure:**
- **Given**: Describe state external actors can observe
  - Use present tense: "Server is running" not "Server was started"
  - Use active voice: "Client has connection" not "Connection is established"
  - Use third person: "Client connects" not "I connect"

- **When**: Describe action by external actor
  - Use present tense: "Client sends message" not "Client sent message"
  - Use active voice: "User clicks button" not "Button is clicked"
  - Use third person: "User enters credentials" not "I enter credentials"

- **Then**: Describe observable outcome
  - Use present tense: "Server responds" not "Server responded"
  - Use active voice: "System displays result" not "Result is displayed"
  - Be specific: "Server returns status 200" not "Server behaves correctly"

**The "Product Owner Test"** - Every scenario must pass:
1. Would a Product Owner understand without asking questions?
2. Does it describe observable system behavior?
3. Does it avoid mentioning internal components?
4. Is it written in active voice?

**Step Structure with Array Format:**

Each Given/When/Then is an array where:
- **First element**: Plain string (main statement)
- **Additional elements**: Objects with `and:` or `but:` keys

**Basic Scenario Format (Most Common - No Modifiers):**
```yaml
id: "3.1-INT-001"
acceptance_criteria: ["AC-1"]
steps:
  - given:
      - "Server is ready to accept WebSocket connections"
  - when:
      - "Client attempts to establish connection"
  - then:
      - "Server accepts connection"
level: "integration"
priority: "P0"
```

**Scenario with 'And' Modifiers in Given (Multiple Preconditions):**
```yaml
id: "3.1-INT-002"
acceptance_criteria: ["AC-1"]
steps:
  - given:
      - "Server is ready to accept connections"
      - and: "MCP endpoint is configured with authentication"
      - and: "Redis cache is available"
  - when:
      - "Client attempts to connect"
  - then:
      - "Server accepts connection"
level: "integration"
priority: "P0"
```

**Scenario with 'And' Modifiers in When (Multiple Actions):**
```yaml
id: "3.1-INT-003"
acceptance_criteria: ["AC-3"]
steps:
  - given:
      - "Client has active WebSocket connection"
  - when:
      - "Client sends authentication request"
      - and: "Client provides valid credentials"
  - then:
      - "Server returns authentication success"
level: "integration"
priority: "P0"
```

**Scenario with 'And' Modifiers in Then (Multiple Outcomes):**
```yaml
id: "3.1-INT-004"
acceptance_criteria: ["AC-4"]
steps:
  - given:
      - "Server is running normally"
  - when:
      - "Client sends valid request"
  - then:
      - "Server returns success response"
      - and: "Response includes correlation ID"
      - and: "Metrics are updated"
level: "integration"
priority: "P0"
```

**Scenario with 'But' Modifiers (Contrasting/Negative Conditions):**
```yaml
id: "3.1-INT-005"
acceptance_criteria: ["AC-3"]
steps:
  - given:
      - "Server is running with rate limiting enabled"
      - but: "No requests have been made yet"
  - when:
      - "Client sends invalid request"
  - then:
      - "Server returns error response"
      - but: "Connection remains active"
      - but: "No alarm is triggered"
level: "integration"
priority: "P0"
```

**Important Rules:**
- **All three keywords** (Given, When, Then) are arrays
- First element must be a plain string (main statement)
- Additional elements are objects with `and:` or `but:` keys
- Use `and:` for additional preconditions, actions, or outcomes
- Use `but:` for contrasting or negative conditions (rare)
- Most scenarios should have 0-2 additional elements per step
- Plain strings for main statements = cleaner, more readable YAML

**Data-Driven Testing with Scenario Outlines:**

Use Scenario Outlines when testing the same behavior with different inputs:

```yaml
id: "3.1-INT-003"
acceptance_criteria: ["AC-2"]
scenario_outline: true
steps:
  - given:
      - "Server is running on port <port>"
  - when:
      - "Client sends <method> request to <endpoint>"
  - then:
      - "Response code should be <status>"
      - "Response contains <field>"
examples:
  - port: 8080
    method: "GET"
    endpoint: "/version"
    status: 200
    field: "version"
  - port: 8080
    method: "POST"
    endpoint: "/version"
    status: 405
    field: "error"
level: "integration"
priority: "P0"
```

**When to Use Scenario Outlines:**
- ✅ Same behavior, different data (HTTP methods, status codes, validation rules)
- ✅ Boundary testing (min/max values, edge cases)
- ✅ Error conditions with different inputs
- ❌ Different behaviors (use separate scenarios instead)
- ❌ Complex setup variations (keep examples simple)

---

### Step 4: Self-Validation Checklist (MANDATORY)

Before outputting, validate EACH scenario against ALL checks:

#### Structure Validation
☐ Has exactly one Given-When-Then sequence
☐ No multiple When-Then pairs
☐ Strict Given → When → Then order
☐ Tests single behavior only

#### Language Validation
☐ Active voice (no "is initialized", "is configured", "is processed")
☐ Third-person perspective maintained
☐ Present tense throughout
☐ No passive constructions

#### Content Validation
☐ Declarative style (WHAT not HOW)
☐ No component names (ConnectionManager, MessageValidator, etc.)
☐ No technical implementation details (middleware, handlers, etc.)
☐ No vague qualifiers ("properly", "correctly", "specific")

#### Quality Validation
☐ Product Owner understandable
☐ Observable from outside system
☐ Each step < 120 characters
☐ Integration or E2E level (NOT unit)

**If ANY check fails**: REGENERATE that scenario and re-validate

---

### Step 5: Examples - Learn from These

#### ❌ BAD: Unit-level, Technical
```yaml
id: "3.1-UNIT-001"  # ← FORBIDDEN LEVEL
given: "Message validator is initialized with MCP schema"  # ← Component name, passive
when: "Invalid message format is received"  # ← Passive voice
then: "Validation error is returned with details"  # ← Passive voice
level: "unit"  # ← NEVER USE THIS
```

**Problems**:
- Unit-level (forbidden in BDD)
- Mentions internal component (Message validator)
- Passive voice throughout
- Not Product Owner understandable

---

#### ✅ GOOD: Integration, BDD-compliant
```yaml
id: "3.1-INT-001"
acceptance_criteria: ["AC-1"]
steps:
  - given:
      - "Client has active WebSocket connection"  # ← External state, active
  - when:
      - "Client sends message with invalid format"  # ← External actor, active
  - then:
      - "Server responds with validation error"  # ← Observable outcome, active
level: "integration"
priority: "P0"
```

**Why Good**:
- Integration level (appropriate for BDD)
- No internal components mentioned
- Active voice throughout
- Observable external behavior
- Product Owner understandable

---

#### ❌ BAD: Multiple Behaviors
```yaml
when: "Multiple clients connect and exchange messages"  # ← Two behaviors
then: "All messages handled correctly and connections stable"  # ← Multiple outcomes + vague
```

**Problems**:
- Tests two behaviors (connection + message handling)
- Vague outcomes ("correctly", "stable")
- Cannot clearly determine pass/fail

---

#### ✅ GOOD: Single Behavior
```yaml
steps:
  - given:
      - "Server is running with connection limit"
  - when:
      - "Multiple clients connect simultaneously"  # ← One behavior
  - then:
      - "Server accepts connections up to configured limit"  # ← Specific outcome
```

**Why Good**:
- Single testable behavior
- Specific, measurable outcome
- Clear pass/fail criteria

---

#### ❌ BAD: Passive Voice, Components
```yaml
given: "Connection manager is initialized"  # ← Passive + component name
when: "Request is processed by handler"  # ← Passive + component name
then: "Response is formatted and returned"  # ← Passive
```

**Problems**:
- Passive voice throughout
- Mentions internal components
- Implementation-focused

---

#### ✅ GOOD: Active Voice, External
```yaml
steps:
  - given:
      - "Server is ready to accept requests"  # ← Active, external state
  - when:
      - "Client sends request to server"  # ← Active, external actor
  - then:
      - "Server returns formatted response"  # ← Active, observable
```

**Why Good**:
- Active voice throughout
- No internal components
- External observable behavior

---

### Forbidden Terms (Auto-Reject if Present)

If ANY of these appear in Given/When/Then → REJECT scenario:

**Component Names:**
- ConnectionManager, MessageValidator, ResponseFormatter
- DocumentService, AuthHandler, TokenValidator
- Any class or internal component name

**Implementation Terms:**
- initialize, instantiate, configure, register
- parse, serialize, deserialize
- allocate, cleanup, pool, thread

**Architecture Terms:**
- middleware, adapter, handler, wrapper
- service, repository, factory, builder

**Vague Qualifiers:**
- properly, correctly, specific
- appropriate, suitable, valid (without criteria)

---

### Required Terms (Use These)

**Actors:**
- Client, User, System, Server
- Administrator, External System
- Claude (for Claude.ai interactions)

**Actions:**
- connect, send, receive, respond, reject
- display, show, navigate, click
- authenticate, authorize

**Artifacts:**
- connection, message, request, response
- page, button, form
- error, result, data

---

### Pre-Output Validation Loop

**BEFORE writing scenarios to file, execute this validation:**

```
FOR EACH generated scenario:

  APPLY Product Owner Test:
    Q1: Would PO understand without questions?
    Q2: Does it describe external behavior?
    Q3: No internal components mentioned?
    Q4: Active voice throughout?

    IF ANY answer is NO:
      REGENERATE scenario
      REPEAT Product Owner Test

  APPLY Validation Checklist:
    [Run through all 16 checks from Step 4]

    IF ANY check fails:
      REGENERATE scenario
      REPEAT validation

  CALCULATE Score:
    [Use 0-10 scoring system]

    IF score < 8:
      REGENERATE scenario
      REPEAT validation

ONLY output scenarios that pass ALL validation
```

---

### Validation Example Walkthrough

**Scenario Under Review:**
```yaml
given: "Connection manager is initialized"
```

**Validation Process:**

**Product Owner Test:**
- Q1: Would PO understand "connection manager"? → NO ❌
- Q2: Describes external behavior? → NO ❌
- Q3: Mentions internal component? → YES ❌
- Q4: Active voice? → NO (passive "is initialized") ❌

**Result**: FAIL - All 4 questions failed

**Regenerate As:**
```yaml
given:
  - "Server is ready to accept connections"
```

**Re-validate:**
- Q1: Would PO understand "server" and "connections"? → YES ✅
- Q2: Describes external behavior? → YES ✅
- Q3: Mentions internal component? → NO ✅
- Q4: Active voice? → YES ✅

**Result**: PASS - Proceed to checklist validation

---

## Output Format

**CRITICAL**: Save text content to file: `./tmp/{{.Story.ID}}-scenarios.yaml`

**Follow EXACTLY this format:**

=== FILE_START: ./tmp/{{.Story.ID}}-scenarios.yaml ===
```yaml
scenarios:
  test_scenarios:
    - id: "{{.Story.ID}}-INT-001"
      acceptance_criteria: ["AC-1"]
      steps:
        - given:
            - "Clear description of initial state"
        - when:
            - "Specific action or event occurs"
        - then:
            - "Expected outcome that can be verified"
      level: "integration"
      priority: "P0"

    - id: "{{.Story.ID}}-E2E-001"
      acceptance_criteria: ["AC-1", "AC-2"]
      steps:
        - given:
            - "Complete system operational from user perspective"
        - when:
            - "User performs complete journey"
        - then:
            - "End-to-end flow completes successfully"
      level: "e2e"
      priority: "P1"
```
=== FILE_END: ./tmp/{{.Story.ID}}-scenarios.yaml ===

**Scenario ID Rules:**
- Format: `{story_id}-{LEVEL}-{sequence}`
- Examples: `3.1-INT-001`, `3.1-E2E-002`
- FORBIDDEN: `3.1-UNIT-001` (no unit scenarios)

**Level Values:**
- `integration` - Direct API/protocol testing (no UI)
- `e2e` - Complete user journey (with UI or Claude chat)
- FORBIDDEN: `unit` (not allowed in BDD scenarios)

**Priority Values:**
- `P0` - Critical: Security, data integrity, compliance
- `P1` - High: Core user journeys, frequent features
- `P2` - Medium: Secondary features, admin functions
- `P3` - Low: Nice-to-have, rare features

---

**COMPLETION SIGNAL:**

After writing the YAML file, respond with ONLY:
```
SCENARIOS_GENERATION_COMPLETE
```

**Do NOT add**:
- Explanations
- Implementation notes
- Conversations
- Additional commentary

---

## User Story
```yaml
id: "{{.Story.ID}}"
title: {{.Story.Title}}
as_a: {{.Story.AsA}}
i_want: {{.Story.IWant}}
so_that: {{.Story.SoThat}}
status: {{.Story.Status}}
acceptance_criteria:
{{- range .Story.AcceptanceCriteria}}
    - id: {{.ID}}
      description: {{.Description}}
{{- end}}
```

## Tasks
{{range $i, $task := .Tasks}}
### Task {{add $i 1}}: {{$task.Name}}
- Acceptance Criteria: [{{range $j, $ac := $task.AcceptanceCriteria}}{{if $j}}, {{end}}{{$ac}}{{end}}]
- Status: {{$task.Status}}
- Subtasks:
{{- range $task.Subtasks}}
  - {{.}}
{{- end}}
{{end}}

## Testing Requirements
```yaml
test_location: "{{.Testing.TestLocation}}"
frameworks:
{{- range .Testing.Frameworks}}
  - "{{.}}"
{{- end}}
requirements:
{{- range .Testing.Requirements}}
  - "{{.}}"
{{- end}}
coverage:
{{- range $key, $value := .Testing.Coverage}}
  {{$key}}: "{{$value}}"
{{- end}}
```

## Development Context

{{if .DevNotes.PreviousStoryInsights}}
**Previous Story Insights:**
{{.DevNotes.PreviousStoryInsights}}
{{end}}

**Technology Stack:**
{{if .DevNotes.TechnologyStack.Description}}{{.DevNotes.TechnologyStack.Description}}{{end}}

**Performance Requirements:**
{{if .DevNotes.PerformanceRequirements}}{{.DevNotes.PerformanceRequirements}}{{end}}

**CRITICAL REQUIREMENT**: Every acceptance criterion ({{range $i, $ac := .Story.AcceptanceCriteria}}{{if $i}}, {{end}}{{$ac.ID}}{{end}}) MUST be referenced in at least one scenario's acceptance_criteria list.

Remember:
1. EVERY acceptance criterion must be covered
2. Use clear, testable Given-When-Then format
3. Generate ONLY INT or E2E scenarios (NEVER unit)
4. Apply validation checklist to every scenario
5. Regenerate scenarios that fail validation
6. Follow the exact YAML format
7. End with SCENARIOS_GENERATION_COMPLETE
