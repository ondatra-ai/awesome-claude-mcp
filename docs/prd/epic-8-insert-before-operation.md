# Epic 8: Insert Before Operation

**Goal:** Enable content insertion before matched anchor text

## User Stories

### Story 8.1: Insert Before Command Handler
**As a** Claude User
**I want** to insert content before specific text
**So that** I can add context to existing content

**Acceptance Criteria:**
- Command handler for insert_before created
- Anchor text matching working
- Content inserted before match
- First match only (MVP)
- Success response provided
- Position tracking accurate

### Story 8.2: Anchor Position Detection
**As a** Developer/Maintainer
**I want** to find anchor text position
**So that** I can insert content before it

**Acceptance Criteria:**
- Anchor search implemented
- Exact position determined
- Before position calculated
- Paragraph boundaries respected
- Format boundaries considered
- Performance acceptable

### Story 8.3: Before Insertion Logic
**As a** Developer/Maintainer
**I want** to insert content before anchor
**So that** document flows naturally

**Acceptance Criteria:**
- Content inserted at correct position
- Anchor text unchanged
- Spacing handled properly
- Formatting applied correctly
- Document structure maintained
- No content overwritten
