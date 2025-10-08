# BDD & Gherkin Best Practices Framework

**Comprehensive Guide for Writing Effective BDD Scenarios**

This document synthesizes authoritative BDD practices from Automation Panda and Cucumber documentation, providing a complete framework for writing high-quality Gherkin scenarios.

---

## Table of Contents

1. [Core BDD Principles](#core-bdd-principles)
2. [The Two Golden Rules](#the-two-golden-rules)
3. [BDD Three-Practice Framework](#bdd-three-practice-framework)
4. [Proper Behavior Structure](#proper-behavior-structure)
5. [Step Writing Guidelines](#step-writing-guidelines)
6. [Declarative vs Imperative Style](#declarative-vs-imperative-style)
7. [Anti-Patterns to Avoid](#anti-patterns-to-avoid)
8. [Scenario Outline Best Practices](#scenario-outline-best-practices)
9. [Handling Test Data](#handling-test-data)
10. [Style and Structure Guidelines](#style-and-structure-guidelines)
11. [Validation Checklist](#validation-checklist)

---

## Core BDD Principles

### What is BDD?

BDD is a software development process that closes the gap between business and technical people by:

- **Encouraging collaboration** across roles to build shared understanding
- **Working in rapid iterations** to increase feedback and flow of value
- **Producing living documentation** that is automatically checked against system behavior

### BDD Philosophy

- Focus collaborative work around **concrete, real-world examples**
- Use examples to guide from concept through implementation
- Treat documentation as a **collaborative asset**
- Continuously evolve shared understanding

### Key Principles

1. **Behavior over Implementation**: Describe *what*, not *how*
2. **Concrete over Abstract**: Mention names, dates, amounts relevant to problem domain
3. **Observable Behavior**: Only test what's visible externally
4. **Technology-Agnostic**: Write examples that could exist without computers (imagine it's 1922)
5. **Collaboration First**: Discovery conversations are more important than automation

---

## The Two Golden Rules

### 1. The Golden Gherkin Rule

**"Write Gherkin so that people who don't know the feature will understand it."**

- Treat readers as you want to be treated
- Make scenarios understandable to non-technical stakeholders
- Scenarios should be self-explanatory

### 2. The Cardinal Rule of BDD

**"One Scenario, One Behavior!"**

- Each scenario tests exactly one independent behavior
- Never combine multiple behaviors in a single scenario
- Avoid multiple When-Then pairs

---

## BDD Three-Practice Framework

### 1. Discovery: What it *could* do

**Goal**: Have right conversations at right time

- Use structured **discovery workshops**
- Focus on real-world examples from users' perspective
- Identify valuable behaviors and scope
- Reveal gaps in understanding
- Defer low-priority functionality

**Key Quote**:
> "The hardest single part of building a software system is deciding precisely what to build." — Fred Brooks

### 2. Formulation: What it *should* do

**Goal**: Document examples in executable format

- Create structured documentation readable by humans and computers
- Establish **shared language** for talking about system
- Get feedback from whole team
- Use problem-domain terminology throughout

### 3. Automation: What it *actually does*

**Goal**: Guide development with executable specifications

- Automate examples one at a time
- Use examples as guide-rails for development
- Reduce manual regression testing burden
- Enable safe maintenance and changes

**Important**: These practices must be done in order. You cannot skip discovery and jump to automation.

---

## Proper Behavior Structure

### Step Type Integrity

**Given-When-Then must appear in order and cannot repeat:**

- **Given**: Set up initial state (not actions)
  - Establishes preconditions
  - Uses present tense for observable state

- **When**: Perform an action (present tense)
  - Triggers the behavior
  - External actor performs action

- **Then**: Verify outcomes (present tense)
  - Observable results
  - Expected behavior validation

- **And/But**: Continue previous step type
  - Use sparingly
  - Must logically extend the previous step

### Structure Rules

✅ **Allowed**:
- Given → When → Then
- Given → And → When → Then
- Given → When → Then → And

❌ **Forbidden**:
- Multiple When-Then pairs (two behaviors)
- Given after When or Then
- When after Then
- Repeating Given-When-Then sequences

### Bad Example (Two Behaviors)

```gherkin
# BAD EXAMPLE! Do not copy.
Scenario: Google Image search shows pictures
  Given the user opens a web browser
  And the user navigates to "https://www.google.com/"
  When the user enters "panda" into the search bar
  Then links related to "panda" are shown on the results page
  When the user clicks on the "Images" link              # ❌ Second When
  Then images related to "panda" are shown               # ❌ Second Then
```

**Problems**:
- Two When-Then pairs = Two behaviors
- Should be split into two scenarios

### Good Example (One Behavior Each)

```gherkin
Scenario: Search from the search bar
  Given a web browser is at the Google home page
  When the user enters "panda" into the search bar
  Then links related to "panda" are shown on the results page

Scenario: Image search
  Given Google search results for "panda" are shown
  When the user clicks on the "Images" link
  Then images related to "panda" are shown on the results page
```

---

## Step Writing Guidelines

### Point of View

**Always use third-person perspective:**

✅ Good:
- "the user logs in"
- "the customer enters payment details"
- "Bob submits the form"

❌ Bad:
- "I log in" (first person)
- "you enter details" (second person)

### Tense

**Always use present tense for all step types:**

✅ Good:
- Given: "the Google home page is displayed"
- When: "the user enters 'panda' into the search bar"
- Then: "links related to 'panda' are shown"

❌ Bad:
- "the user navigates" (Given should establish state, not action)
- "the user entered" (past tense)
- "links will be shown" (future tense)

### Subject-Predicate Action Phrases

**Write steps as complete subject-predicate phrases:**

✅ Good:
- "the results page shows links related to 'panda'"
- "the user clicks the submit button"

❌ Bad:
- "links related to 'panda'" (incomplete - missing subject or verb)
- "image links for 'panda'" (ambiguous)

### Given Step Best Practices

**Given establishes state, NOT actions:**

✅ Good:
- "the Google home page is displayed"
- "the user has a valid account"
- "the server is running on port 8080"

❌ Bad:
- "the user navigates to the home page" (action)
- "the user logs in" (action - this is a When)

---

## Declarative vs Imperative Style

### Imperative Style (Bad)

**Describes HOW actions happen step-by-step:**

```gherkin
# BAD EXAMPLE - Imperative
When the user scrolls the mouse to the search bar
And the user clicks the search bar
And the user types the letter "p"
And the user types the letter "a"
And the user types the letter "n"
And the user types the letter "d"
And the user types the letter "a"
And the user types the ENTER key
```

**Problems**:
- Too detailed
- Brittle (breaks when UI changes)
- Hard to read
- Procedure-driven, not behavior-driven

### Declarative Style (Good)

**Describes WHAT should happen:**

```gherkin
# GOOD EXAMPLE - Declarative
When the user enters "panda" at the search bar
```

**Benefits**:
- Concise and clear
- Resilient to implementation changes
- Focuses on behavior, not mechanics
- Easy to understand

### Real-World Comparison

#### Imperative (Bad):
```gherkin
Feature: Subscribers see different articles

Scenario: Free subscribers see only free articles
  Given users with free subscription can access "FreeArticle1"
  When I type "freeFrieda@example.com" in the email field
  And I type "validPassword123" in the password field
  And I press the "Submit" button
  Then I see "FreeArticle1" on the home page
  And I do not see "PaidArticle1" on the home page
```

#### Declarative (Good):
```gherkin
Feature: Subscribers see different articles

Scenario: Free subscribers see only free articles
  Given Free Frieda has a free subscription
  When Free Frieda logs in with her valid credentials
  Then she sees a Free article
```

**Key Differences**:
- Declarative hides implementation details
- More resilient to UI changes
- Easier to understand behavior
- Step definitions handle the "how"

---

## Anti-Patterns to Avoid

### 1. No "Or" Step

**Gherkin does not have "Or" logic:**

❌ Bad:
```gherkin
When the player pushes the "A" button
Or the player pushes the "B" button
```

✅ Good - Use Scenario Outline:
```gherkin
Scenario Outline: Mario jumps
  Given a level is started
  When the player pushes the "<button>" button
  Then Mario jumps straight up

  Examples:
    | button |
    | A      |
    | B      |
```

### 2. Passive Voice

**Never use passive constructions:**

❌ Bad Patterns:
- "Server is ready to accept connections"
- "Connection is established"
- "Request is processed"
- "Data is stored"
- "System is configured"

✅ Good - Active Voice:
- "Server accepts connections"
- "Client establishes connection"
- "Server processes request"
- "Database stores data"
- "System requires configuration"

**The "Remove State Verb" Test**:
- If step uses "is/are/was/were/has/have" → Check if it's passive
- ❌ "Server is ready" → Cannot remove "is" = Passive
- ✅ "Server accepts connections" → No state verb = Active

### 3. Implementation Details

**Never mention internal components:**

❌ Bad:
- ConnectionManager
- MessageValidator
- ResponseFormatter
- middleware
- handlers
- internal state

✅ Good:
- Server
- Client
- System
- Service
- User

### 4. Vague Qualifiers

**Avoid imprecise language:**

❌ Bad:
- "properly configured"
- "correctly processed"
- "specific error message"

✅ Good:
- "configured to accept SSL connections"
- "processed and returns status 200"
- "error message 'Invalid credentials'"

### 5. Hard-Coded Test Data

**Write defensively for changing data:**

❌ Bad:
```gherkin
Then the following related results are shown
  | related       |
  | Panda Express |
  | giant panda   |
  | panda videos  |
```

**Problem**: Specific results may change over time

✅ Good:
```gherkin
Then links related to "panda" are shown on the results page
```

**Solution**: Step definition intelligently verifies results relate to search term

---

## Scenario Outline Best Practices

### When to Use

**Use Scenario Outline when testing same behavior with different data:**

```gherkin
Scenario Outline: Search for different terms
  Given a web browser is at the Google home page
  When the user searches for "<search_term>"
  Then results for "<search_term>" are displayed

  Examples:
    | search_term |
    | panda       |
    | elephant    |
    | tiger       |
```

### Questions to Ask

**When facing oversized scenario outlines:**

1. **Does each row represent an equivalence class?**
   - Don't test trivial variations
   - Focus on meaningful differences

2. **Does every combination need coverage?**
   - N columns with M inputs = M^N combinations
   - Consider reducing to single appearance per input

3. **Do columns represent separate behaviors?**
   - If columns never referenced together, split scenarios

4. **Does reader need explicit data?**
   - Consider hiding data in step definitions
   - Some data may be derivable

---

## Handling Test Data

### Principles

- BDD is **specification by example**
- Data should support descriptive nature
- Think of data as **examples of behavior**, not test data

### Hiding Data in Automation

**Step definitions can hide unnecessary details:**

```gherkin
Scenario: Search result linking
  Given Google search results for "panda" are shown
  When the user clicks the first result link
  Then the page for the chosen result link is displayed
```

**Key Point**: Step definition stores and passes result link value behind the scenes

### When to Show Data

**Show data when it's meaningful to behavior:**

✅ Show:
- Boundary values
- Examples illustrating rules
- Data that defines the behavior

❌ Hide:
- Implementation details
- Technical IDs
- Derived values
- Constantly changing data

---

## Style and Structure Guidelines

### Content Organization

**Feature Files:**
- Focus feature on customer needs
- Limit one feature per file (easy to find)
- ~12 scenarios per feature (avoid thousand-line files)

**Scenarios:**
- <10 steps per scenario (single-digit recommended)
- <80-120 characters per step
- Short and sweet is best

**Titles:**
- Specific and descriptive
- Communicate behavior in one line
- First thing people read

### Language and Grammar

**Spelling and Grammar:**
- ✅ Use proper spelling
- ✅ Use proper grammar
- ✅ No punctuation at end of steps (no periods/commas)
- ✅ Single spaces between words

**Capitalization:**
- ✅ Capitalize Gherkin keywords: Feature, Scenario, Given, When, Then, And, But
- ✅ Capitalize first word in titles
- ✅ Do NOT capitalize words in steps unless proper nouns

**Examples:**

✅ Good:
```gherkin
Feature: Google Searching

  Scenario: Simple search
    Given the Google home page is displayed
    When the user enters "panda" into the search bar
    Then links related to "panda" are shown
```

❌ Bad:
```gherkin
feature: Google Searching

  SCENARIO: SIMPLE SEARCH
    given the Google Home Page is Displayed.
    When The User Enters "panda" Into the search bar,
    then Links Related to "panda" Are Shown.
```

### Formatting

**Indentation and Spacing:**
- Indent content beneath every section header
- 2 blank lines between features and scenarios
- 1 blank line between example tables
- NO blank lines between steps within scenario
- Space table pipes (|) evenly

**Tags:**
- Standard set of tag names
- All lowercase
- Use hyphens (-) for multi-word tags
- Keep tag names short

✅ Good tags:
```gherkin
@smoke @authentication @regression
```

❌ Bad tags:
```gherkin
@AUTOMATE @Automated @automation @Sprint32GoogleSearchFeature
```

---

## Validation Checklist

### Before Outputting Any Scenario

**Structure Checks:**
- ☐ Has exactly one Given-When-Then sequence
- ☐ No multiple When-Then pairs
- ☐ Strict Given → When → Then order
- ☐ Tests single behavior only

**Language Checks:**
- ☐ Active voice throughout (no "is initialized")
- ☐ Third-person perspective maintained
- ☐ Present tense in all steps
- ☐ No passive constructions

**Content Checks:**
- ☐ Declarative style (WHAT not HOW)
- ☐ No component names mentioned
- ☐ No technical implementation details
- ☐ No vague qualifiers

**Quality Checks:**
- ☐ Product Owner understandable
- ☐ Observable from outside system
- ☐ Each step < 120 characters
- ☐ Integration or E2E level (never unit)

### The "Product Owner Test"

**Every scenario MUST pass all 4 questions:**

1. Would a Product Owner understand without asking questions?
   - If NO → REJECT

2. Does it describe observable system behavior?
   - If NO → REJECT

3. Does it avoid internal components?
   - If NO → REJECT

4. Is it written in active voice?
   - If NO → REJECT

### Scoring System

**Each scenario scores 0-10 points:**
- +1 for active voice in each step (max 3)
- +1 for declarative style
- +1 for no technical terms
- +1 for single behavior
- +1 for < 100 chars per step
- +1 for Product Owner understandable
- +2 for no component names
- +1 for INT/E2E level

**Minimum Score**: 8/10 required

**If scenario scores < 8**: REGENERATE

---

## Summary: Essential Practices

### Top 10 Rules

1. **Golden Rule**: Write so people who don't know feature will understand
2. **Cardinal Rule**: One Scenario, One Behavior
3. **Third-person**: All steps in third-person POV
4. **Present tense**: All steps use present tense
5. **Active voice**: Never use passive constructions
6. **Declarative**: Describe WHAT, not HOW
7. **Short scenarios**: <10 steps recommended
8. **Subject-predicate**: Complete action phrases
9. **Step order**: Given → When → Then (no repeating)
10. **Observable**: Only test externally visible behavior

### When in Doubt

**Ask yourself:**
- Would this make sense to someone unfamiliar with the system?
- Am I describing behavior or implementation?
- Could this be written without computers (1922 test)?
- Is this one behavior or multiple behaviors?

### Remember

> "Good Gherkin comes from good behavior."

**Quality over speed**: Take time to validate. Regenerate until perfect.

---

## Additional Resources

### From Automation Panda:
- [BDD 101: Writing Good Gherkin](https://automationpanda.com/2017/01/30/bdd-101-writing-good-gherkin/)
- [Should Gherkin Steps Use First-Person or Third-Person?](https://automationpanda.com/2017/01/26/should-gherkin-steps-use-first-person-or-third-person/)
- [Should Gherkin Steps use Past, Present, or Future Tense?](https://automationpanda.com/2018/05/17/should-gherkin-steps-use-past-present-or-future-tense/)
- [Handling Test Data in BDD](https://automationpanda.com/2017/08/05/handling-test-data-in-bdd/)

### From Cucumber:
- [Behaviour-Driven Development](https://cucumber.io/docs/bdd/)
- [Writing Better Gherkin](https://cucumber.io/docs/bdd/better-gherkin/)
- [Discovery Workshop](https://cucumber.io/docs/bdd/discovery-workshop/)
- [Example Mapping](https://cucumber.io/docs/bdd/example-mapping/)

---

*This framework synthesizes best practices from industry-leading sources to provide comprehensive guidance for writing effective BDD scenarios using Gherkin. All principles are derived from authoritative sources including Andy Knight's Automation Panda blog and official Cucumber documentation.*
