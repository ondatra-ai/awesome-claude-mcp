# BDD 101: Writing Good Gherkin

**Source:** [Automation Panda - BDD 101: Writing Good Gherkin](https://automationpanda.com/2017/01/30/bdd-101-writing-good-gherkin/)
**Author:** Andy Knight (Automation Panda)
**Date:** January 30, 2017

---

## The Golden Gherkin Rule

**Treat other readers as you would want to be treated. Write Gherkin so that people who don't know the feature will understand it.**

---

## Proper Behavior

### The Cardinal Rule of BDD

**One Scenario, One Behavior!**

The biggest mistake BDD beginners make is writing Gherkin without a behavior-driven mindset. They often write feature files as if they are writing "traditional" procedure-driven functional tests: step-by-step instructions with actions and expected results.

### Bad Example (Procedure-Driven)

```gherkin
# BAD EXAMPLE! Do not copy.
Feature: Google Searching

  Scenario: Google Image search shows pictures
    Given the user opens a web browser
    And the user navigates to "https://www.google.com/"
    When the user enters "panda" into the search bar
    Then links related to "panda" are shown on the results page
    When the user clicks on the "Images" link at the top of the results page
    Then images related to "panda" are shown on the results page
```

**Problems with this scenario:**
- Has two When-Then pairs (covers two behaviors, not one)
- First two steps are purely setup and strongly imperative
- Given-When-Then steps must appear in order and cannot repeat
- A Given may not follow a When or Then
- A When may not follow a Then
- Any single When-Then pair denotes an individual behavior

### Good Example (Behavior-Driven)

```gherkin
Feature: Google Searching

  Scenario: Search from the search bar
    Given a web browser is at the Google home page
    When the user enters "panda" into the search bar
    Then links related to "panda" are shown on the results page

  Scenario: Image search
    Given Google search results for "panda" are shown
    When the user clicks on the "Images" link at the top of the results page
    Then images related to "panda" are shown on the results page
```

**Key Points:**
- Each scenario covers one independent behavior
- Steps are declarative, not imperative
- Second scenario can run independently (Given step establishes required state)
- Behavior scenarios represent requirements and acceptance criteria
- **Good Gherkin comes from good behavior**

### Respecting Step Type Integrity

**Do not arbitrarily reassign step types to make scenarios follow Given-When-Then ordering.**

- **Given** steps set up initial state
- **When** steps perform an action
- **Then** steps verify outcomes

Step types are meant to be guide rails for writing good behavior scenarios.

---

## Phrasing Steps

### Third-Person Point of View

**Write all steps in third-person point of view.**

Mixing first-person and third-person steps makes scenarios confusing.

### Subject-Predicate Action Phrases

**Write steps as a subject-predicate action phrase.**

Leaving parts of speech out makes steps ambiguous and more likely to be reused improperly.

#### Bad Example

```gherkin
# BAD EXAMPLE! Do not copy.
Feature: Google Searching

  Scenario: Google search result page elements
    Given the user navigates to the Google home page
    When the user entered "panda" at the search bar
    Then the results page shows links related to "panda"
    And image links for "panda"
    And video links for "panda"
```

**Problem:** The final two And steps lack subject-predicate phrase format. Are the links subjects (performing action) or objects (receiving action)?

### Tense Usage

**Use present tense for all step types.**

#### Bad Example

```gherkin
# BAD EXAMPLE! Do not copy.
Feature: Google Searching

  Scenario: Simple Google search
    Given the user navigates to the Google home page
    When the user entered "panda" at the search bar
    Then links related to "panda" will be shown on the results page
```

**Problems:**
- Given step indicates action ("navigates") when it should establish state
- When step uses past tense ("entered") instead of present
- Then step uses future tense ("will be shown") instead of present

#### Good Example

```gherkin
Feature: Google Searching

  Scenario: Simple Google search
    Given the Google home page is displayed
    When the user enters "panda" into the search bar
    Then links related to "panda" are shown on the results page
```

**Key Points:**
- Given step establishes state, not action
- When step uses present tense for current action
- Then step uses present tense (behaviors are present-tense aspects of the product)
- All steps written in third-person

---

## Good Titles

**The title is like the face of a scenario** – it's the first thing people read.

- Must communicate in one concise line what the behavior is
- Often logged by the automation framework
- Should be specific and descriptive

---

## Choices, Choices

### No "Or" Step Exists

**Gherkin does not have an "Or" step.**

When automated, every step is executed sequentially.

#### Bad Example

```gherkin
# BAD EXAMPLE! Do not copy.
Feature: SNES Mario Controls

  Scenario: Mario jumps
    Given a level is started
    When the player pushes the "A" button
    Or the player pushes the "B" button
    Then Mario jumps straight up
```

### Use Scenario Outline for Variations

**Use Scenario Outline sections to cover multiple variations of the same behavior.**

#### Good Example

```gherkin
Feature: SNES Mario Controls

  Scenario Outline: Mario jumps
    Given a level is started
    When the player pushes the "<letter>" button
    Then Mario jumps straight up

    Examples: Buttons
      | letter |
      | A      |
      | B      |
```

---

## The Known Unknowns

### Handling Changing Data

**Write scenarios defensively so that changes in the underlying data do not cause test runs to fail.**

**Think about data not as test data but as examples of behavior.**

#### Problem Example

```gherkin
Feature: Google Searching

  Scenario: Simple Google search
    Given a web browser is on the Google page
    When the search phrase "panda" is entered
    Then results for "panda" are shown
    And the following related results are shown
      | related       |
      | Panda Express |
      | giant panda   |
      | panda videos  |
```

**Problem:** Hard-coded results may change over time (e.g., Panda Express goes out of business).

#### Better Approach

```gherkin
Feature: Google Searching

  Scenario: Simple Google search
    Given a web browser is on the Google page
    When the search phrase "panda" is entered
    Then results for "panda" are shown
    And links related to "panda" are shown on the results page
```

**Solution:** Step definition implementation can intelligently verify that each result somehow relates to the search phrase.

### Hiding Data in Automation

**Step definitions can hide data in the automation when it doesn't need to be exposed.**

#### Example

```gherkin
Feature: Google Searching

  Scenario: Search result linking
    Given Google search results for "panda" are shown
    When the user clicks the first result link
    Then the page for the chosen result link is displayed
```

**Key Point:** The When step doesn't explicitly name the value of the result link. Behind the scenes, the step definition stores the value and passes it forward to the Then step.

---

## Handling Test Data

BDD is **specification by example** – scenarios should be descriptive of the behaviors they cover.

**Any data written into the Gherkin should support that descriptive nature.**

Some types of test data should be handled directly within Gherkin, but other types should not.

---

## Less is More

### Scenario Length

**Scenarios should be short and sweet.**

- Recommended: single-digit step count (<10)
- Long scenarios are hard to understand
- Often indicative of poor practices

### Declarative vs. Imperative Steps

**Imperative steps** state the mechanics of *how* an action should happen (procedure-driven).

#### Bad Example (Imperative)

```gherkin
When the user scrolls the mouse to the search bar
And the user clicks the search bar
And the user types the letter "p"
And the user types the letter "a"
And the user types the letter "n"
And the user types the letter "d"
And the user types the letter "a"
And the user types the ENTER key
```

**Declarative steps** state *what* action should happen without providing all information for how (behavior-driven).

#### Good Example (Declarative)

```gherkin
When the user enters "panda" at the search bar
```

**Key Point:** The scrolling and keystroking is implied and handled by automation in the step definition.

### Scenario Outline Best Practices

**Scenario outlines should focus on one behavior and use only the necessary variations.**

Questions to ask when facing an oversized scenario outline:

1. **Does each row represent an equivalence class of variations?**
   - Searching for "elephant" in addition to "panda" may not add much test value

2. **Does every combination of inputs need to be covered?**
   - N columns with M inputs each generates M^N possible combinations
   - Consider making each input appear only once, regardless of combination

3. **Do any columns represent separate behaviors?**
   - May be true if columns are never referenced together in the same step
   - If so, consider splitting apart the scenario outline by column

4. **Does the feature file reader need to explicitly know all of the data?**
   - Consider hiding some data in step definitions
   - Some data may be derivable from other data

---

## Style and Structure

**Good writing style improves communication.**

In a truly behavior-driven team, non-technical stakeholders will rely upon feature files just as much as engineers.

### Style Guidelines

**Content Organization:**
- Focus a feature on customer needs
- Limit one feature per feature file (makes it easy to find features)
- Limit the number of scenarios per feature (avoid thousand-line feature files)
  - A good measure is a dozen scenarios per feature
- Limit the number of steps per scenario to less than ten
- Limit the character length of each step (common limits: 80-120 characters)

**Language and Grammar:**
- Use proper spelling
- Use proper grammar
- Capitalize Gherkin keywords (Feature, Scenario, Given, When, Then, And, But)
- Capitalize the first word in titles
- Do not capitalize words in step phrases unless they are proper nouns
- Do not use punctuation (periods and commas) at the end of step phrases
- Use single spaces between words

**Formatting:**
- Indent the content beneath every section header
- Separate features and scenarios by two blank lines
- Separate examples tables by 1 blank line
- Do not separate steps within a scenario by blank lines
- Space table delimiter pipes ("|") evenly

**Tags:**
- Adopt a standard set of tag names (avoid duplicates)
- Write all tag names in lowercase
- Use hyphens ("-") to separate words in tag names
- Limit the length of tag names

### Bad Example (Poor Style)

```gherkin
# BAD EXAMPLE! Do not copy.
Feature: Google Searching

@AUTOMATE @Automated @automation @Sprint32GoogleSearchFeature
Scenario outline: GOOGLE STUFF
Given a Web Browser is on the Google page,
when The seach phrase "<phrase>" Enter,
Then "<phrase>" shown.
and The relatedd   results include "<related>".

Examples: animals
 |phrase|related|
| panda | Panda Express        |
| elephant    | elephant Man  |
```

**Problems:**
- Inconsistent capitalization
- Spelling errors ("seach", "relatedd")
- Poor tag naming and duplication
- Inconsistent spacing
- Missing proper formatting

### Key Point

**Gherkin files should look elegant.**

While automation code may look hairy in parts, feature files should be clean and professional.

---

## Gherkinize Those Behaviors!

With these best practices, you can write Gherkin feature files like a pro.

**Key Reminders:**
- Don't be afraid to try – nobody does things perfectly the first time
- Don't give up if you get stuck
- Always remember the **Golden Gherkin Rule**
- Always remember the **Cardinal Rule of BDD**

---

## Summary of Key Rules

### The Two Most Important Rules

1. **The Golden Gherkin Rule:** Write Gherkin so that people who don't know the feature will understand it

2. **The Cardinal Rule of BDD:** One Scenario, One Behavior!

### Essential Practices

1. Each scenario covers one unique, independent behavior
2. Given-When-Then steps must appear in order and cannot repeat
3. Write all steps in third-person point of view
4. Write steps as subject-predicate action phrases
5. Use present tense for all step types
6. Write declaratively (what), not imperatively (how)
7. Keep scenarios short (<10 steps)
8. Use Scenario Outline for behavior variations
9. Write scenarios defensively for changing data
10. Make Gherkin files look elegant and professional

### Step Type Integrity

- **Given:** Set up initial state (not actions)
- **When:** Perform an action (present tense)
- **Then:** Verify outcomes (present tense)
- **And/But:** Continue the previous step type

### What Gherkin Does NOT Have

- No "Or" step (use Scenario Outline instead)
- No conditional logic
- No loops
- Not a programming language – it's a specification language

---

## Additional Resources

For more information on BDD and Gherkin, see:
- [Automation Panda BDD Page](https://automationpanda.com/bdd/)
- [Should Gherkin Steps Use First-Person or Third-Person?](https://automationpanda.com/2017/01/26/should-gherkin-steps-use-first-person-or-third-person/)
- [Should Gherkin Steps use Past, Present, or Future Tense?](https://automationpanda.com/2018/05/17/should-gherkin-steps-use-past-present-or-future-tense/)
- [Good Gherkin Scenario Titles](https://automationpanda.com/2018/10/23/good-gherkin-scenario-titles/)
- [Are Gherkin Scenarios with Multiple When-Then Pairs Okay?](https://automationpanda.com/2018/05/07/are-gherkin-scenarios-with-multiple-when-then-pairs-okay/)
- [Handling Test Data in BDD](https://automationpanda.com/2017/08/05/handling-test-data-in-bdd/)

---

*This document summarizes the key teachings from Andy Knight's article "BDD 101: Writing Good Gherkin" published on Automation Panda. All examples and principles are adapted from the original source to serve as a comprehensive reference guide for writing effective Gherkin scenarios.*
