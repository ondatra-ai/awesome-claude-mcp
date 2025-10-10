<!-- Powered by BMAD™ Core -->

# Implement Playwright Test from BDD Scenario

## Your Task

Generate a Playwright test for scenario **{{.ScenarioID}}** and update the requirements registry.

---

## Scenario Details

**ID**: `{{.ScenarioID}}`
**Description**: {{.Description}}
**Level**: {{.Level}}
**Category**: {{.Category}}
**Priority**: {{.Priority}}

**Test Steps (Given-When-Then)**:
```
{{.FormatSteps}}
```

---

## Step 1: Determine Test File Path

Based on the scenario metadata:
- Level: `{{.Level}}`
- Category: `{{.Category}}`

**Target file**: `tests/{{.Level}}/{{.Category}}.spec.ts`

---

## Step 2: Read Existing Test File (if exists)

```
Read tests/{{.Level}}/{{.Category}}.spec.ts
```

If the file exists, analyze:
- Import statements pattern
- test.describe structure
- Existing test naming convention
- Assertion style (expect patterns)
- Variable naming conventions

If the file doesn't exist, you'll create it with proper structure.

---

## Step 3: Generate Test Code

### For Integration Tests (API)
Use Playwright Request API (`{ request }`):

```typescript
test('{{.ScenarioID}}: {{.Description}}', async ({ request }) => {
  // Given: [map Given steps to setup]

  // When: [map When steps to API call]
  const response = await request.get(`${backendUrl}/endpoint`);

  // Then: [map Then steps to assertions]
  expect(response.status()).toBe(200);
  const data = await response.json();
  expect(data).toHaveProperty('key', 'value');
});
```

### For E2E Tests (Browser)
Use Playwright Browser API (`{ page }`):

```typescript
test('{{.ScenarioID}}: {{.Description}}', async ({ page }) => {
  // Given: [map Given steps to navigation/setup]
  await page.goto(frontendUrl);

  // When: [map When steps to interactions]
  await page.click('button');

  // Then: [map Then steps to visibility/content checks]
  await expect(page.locator('.result')).toBeVisible();
  await expect(page.locator('.result')).toHaveText('Expected');
});
```

---

## Step 4: Add Test to File

### If file exists:
- Append test inside the existing `test.describe()` block
- Match the existing code style and formatting
- Preserve all existing tests

### If file doesn't exist:
Create complete file structure:
```typescript
import { test, expect } from '@playwright/test';
import { getEnvironmentConfig } from '../config/environments';

const { backendUrl } = getEnvironmentConfig(process.env.E2E_ENV);

test.describe('{{.Category | title}} {{.Level | title}} Tests', () => {

  test('{{.ScenarioID}}: {{.Description}}', async ({ request }) => {
    // Implementation here
  });

});
```

---

## Step 5: Update Requirements Registry

Update `{{.RequirementsFile}}` for scenario `{{.ScenarioID}}`:

```yaml
{{.ScenarioID}}:
  # ... keep existing fields ...
  implementation_status:
    status: "implemented"  # Change from "pending"
    file_path: "tests/{{.Level}}/{{.Category}}.spec.ts"
  # ... keep other fields unchanged ...
```

Use the Edit tool to update only the `implementation_status` section.

---

## Mapping Guide: Given-When-Then → Playwright Code

### Given (Preconditions)
- "service runs" → Setup base URL/config
- "user authenticated" → Add auth headers/cookies
- "data exists in database" → Pre-populate via API
- "page loaded" → `await page.goto(url)`

### When (Actions)
- "client sends GET request" → `await request.get(url)`
- "user clicks button" → `await page.click('button')`
- "form submitted" → `await page.fill() + page.click('submit')`
- "API called" → `await request.post(url, { data })`

### Then (Assertions)
- "returns 200 status" → `expect(response.status()).toBe(200)`
- "element visible" → `await expect(locator).toBeVisible()`
- "contains text" → `await expect(locator).toHaveText('text')`
- "property exists" → `expect(data).toHaveProperty('key')`

---

## Important Rules

1. **Test ID in name**: Always prefix test name with `{{.ScenarioID}}:`
2. **Match existing style**: If file exists, follow its patterns exactly
3. **Comments**: Add Given/When/Then comments for clarity
4. **No mocks**: Generate real API/browser interactions
5. **Update requirements**: Always update implementation_status after test creation

---

## Output Requirements

After completing the implementation, provide:

```
Test Implementation Summary:
- Test ID: {{.ScenarioID}}
- Description: {{.Description}}
- File: tests/{{.Level}}/{{.Category}}.spec.ts
- Status: [CREATED | APPENDED]
- Requirements Updated: YES
```

**If any errors occur**, provide clear error message with:
- What step failed
- Error details
- Suggested fix
