# BMad Framework - Agents and Commands Reference

## Overview
BMad (Business Management and Development) is a comprehensive framework for software development lifecycle management, featuring specialized AI agents that handle different aspects of project development.

## Agent Directory

| Agent | ID | Emoji | Primary Role | # Commands |
|-------|-----|-------|--------------|------------|
| **BMad Orchestrator** | bmad-orchestrator | 🎭 | Master coordinator and workflow management | 11 |
| **BMad Master** | bmad-master | 🧙 | Universal task executor across all domains | 10 |
| **Product Owner (Sarah)** | po | 📝 | Backlog management and story refinement | 9 |
| **Product Manager (John)** | pm | 📋 | Product strategy and documentation | 11 |
| **Business Analyst (Mary)** | analyst | 📊 | Research, analysis, and strategic ideation | 9 |
| **Architect (Winston)** | architect | 🏗️ | System design and technical architecture | 11 |
| **QA/Test Architect (Quinn)** | qa | 🧪 | Quality assurance and test strategy | 7 |
| **Full Stack Developer (James)** | dev | 💻 | Code implementation and development | 5 |
| **Scrum Master (Bob)** | sm | 🏃 | Story creation and agile process guidance | 4 |
| **UX Expert (Sally)** | ux-expert | 🎨 | UI/UX design and user experience | 3 |

## Commands by Agent

### 🎭 BMad Orchestrator
| Command | Description |
|---------|-------------|
| `*help` | Show available agents and workflows |
| `*agent [name]` | Transform into specialized agent |
| `*chat-mode` | Start conversational mode |
| `*checklist [name]` | Execute checklist |
| `*doc-out` | Output full document |
| `*kb-mode` | Load full BMad knowledge base |
| `*party-mode` | Group chat with all agents |
| `*status` | Show current context and progress |
| `*task [name]` | Run specific task |
| `*yolo` | Toggle skip confirmations mode |
| `*exit` | Return to BMad or exit session |

### 🧙 BMad Master
| Command | Description |
|---------|-------------|
| `*help` | Show commands list |
| `*create-doc {template}` | Execute document creation task |
| `*doc-out` | Output full document |
| `*document-project` | Execute project documentation task |
| `*execute-checklist {checklist}` | Run checklist execution |
| `*kb` | Toggle knowledge base mode |
| `*shard-doc {document} {destination}` | Shard document task |
| `*task {task}` | Execute any available task |
| `*yolo` | Toggle Yolo Mode |
| `*exit` | Exit with confirmation |

### 📝 Product Owner (Sarah)
| Command | Description |
|---------|-------------|
| `*help` | Show commands list |
| `*correct-course` | Execute course correction task |
| `*create-epic` | Create epic for brownfield projects |
| `*create-story` | Create user story from requirements |
| `*doc-out` | Output full document |
| `*execute-checklist-po` | Run PO master checklist |
| `*shard-doc {document} {destination}` | Shard document task |
| `*validate-story-draft {story}` | Validate story draft |
| `*yolo` | Toggle Yolo Mode |
| `*exit` | Exit with confirmation |

### 📋 Product Manager (John)
| Command | Description |
|---------|-------------|
| `*help` | Show commands list |
| `*correct-course` | Execute course correction task |
| `*create-brownfield-epic` | Create brownfield epic |
| `*create-brownfield-prd` | Create brownfield PRD |
| `*create-brownfield-story` | Create brownfield story |
| `*create-epic` | Create epic for brownfield projects |
| `*create-prd` | Create product requirements document |
| `*create-story` | Create user story |
| `*doc-out` | Output full document |
| `*shard-prd` | Shard PRD document |
| `*yolo` | Toggle Yolo Mode |
| `*exit` | Exit |

### 📊 Business Analyst (Mary)
| Command | Description |
|---------|-------------|
| `*help` | Show commands list |
| `*brainstorm {topic}` | Facilitate structured brainstorming |
| `*create-competitor-analysis` | Create competitive analysis |
| `*create-project-brief` | Create project brief |
| `*doc-out` | Output document in progress |
| `*elicit` | Run advanced elicitation task |
| `*perform-market-research` | Create market research document |
| `*research-prompt {topic}` | Create deep research prompt |
| `*yolo` | Toggle Yolo Mode |
| `*exit` | Exit as Business Analyst |

### 🏗️ Architect (Winston)
| Command | Description |
|---------|-------------|
| `*help` | Show commands list |
| `*create-backend-architecture` | Create backend architecture |
| `*create-brownfield-architecture` | Create brownfield architecture |
| `*create-front-end-architecture` | Create frontend architecture |
| `*create-full-stack-architecture` | Create full-stack architecture |
| `*doc-out` | Output full document |
| `*document-project` | Execute project documentation |
| `*execute-checklist {checklist}` | Run architect checklist (default) |
| `*research {topic}` | Execute deep research prompt |
| `*shard-prd` | Shard architecture document |
| `*yolo` | Toggle Yolo Mode |
| `*exit` | Exit as Architect |

### 🧪 QA/Test Architect (Quinn)
| Command | Description |
|---------|-------------|
| `*help` | Show commands list |
| `*gate {story}` | Execute quality gate decision |
| `*nfr-assess {story}` | Assess non-functional requirements |
| `*review {story}` | Comprehensive story review with gate decision |
| `*risk-profile {story}` | Generate risk assessment matrix |
| `*test-design {story}` | Create comprehensive test scenarios |
| `*trace {story}` | Map requirements to tests using Given-When-Then |
| `*exit` | Exit as Test Architect |

**⚠️ Note:** QA Architect is only authorized to update the "QA Results" section of story files.

### 💻 Full Stack Developer (James)
| Command | Description |
|---------|-------------|
| `*help` | Show commands list |
| `*develop-story` | Implement story tasks with full development workflow |
| `*explain` | Explain recent development work in detail |
| `*review-qa` | Apply QA fixes and recommendations |
| `*run-tests` | Execute linting and tests |
| `*exit` | Exit as Developer |

**⚠️ Note:** Developer only updates specific story sections (Tasks/Subtasks checkboxes, Dev Agent Record, File List, Status).

### 🏃 Scrum Master (Bob)
| Command | Description |
|---------|-------------|
| `*help` | Show commands list |
| `*correct-course` | Execute course correction |
| `*draft` | Execute story creation task |
| `*story-checklist` | Execute story draft checklist |
| `*exit` | Exit as Scrum Master |

**⚠️ Note:** Cannot implement stories or modify code; focuses solely on story preparation.

### 🎨 UX Expert (Sally)
| Command | Description |
|---------|-------------|
| `*help` | Show commands list |
| `*create-front-end-spec` | Create frontend specification |
| `*generate-ui-prompt` | Generate AI frontend prompts |
| `*exit` | Exit as UX Expert |

## Common Commands Across Multiple Agents

| Command | Available in Agents |
|---------|-------------------|
| `*help` | All agents |
| `*exit` | All agents |
| `*yolo` | BMad Orchestrator, BMad Master, PO, PM, Analyst, Architect |
| `*doc-out` | BMad Orchestrator, BMad Master, PO, PM, Analyst, Architect |
| `*correct-course` | PO, PM, SM |
| `*create-story` | PO, PM |
| `*create-epic` | PO, PM |

## Task Categories

### Document Creation
- `advanced-elicitation` - Enhanced content quality through structured techniques
- `create-doc` - Template-driven document creation
- `document-project` - Project documentation
- `shard-doc` - Document sharding/splitting

### Story Management
- `create-next-story` - Sequential story creation
- `create-brownfield-story` - Existing project story creation
- `validate-next-story` - Story validation
- `review-story` - Comprehensive story review

### Quality Assurance
- `qa-gate` - Quality gate decisions
- `nfr-assess` - Non-functional requirements assessment
- `test-design` - Test scenario creation
- `trace-requirements` - Requirements traceability
- `apply-qa-fixes` - Apply QA recommendations

### Research & Analysis
- `create-deep-research-prompt` - Deep research prompt generation
- `facilitate-brainstorming-session` - Interactive brainstorming
- `advanced-elicitation` - Enhanced elicitation techniques

### Process Management
- `execute-checklist` - Checklist execution
- `correct-course` - Course correction
- `kb-mode-interaction` - Knowledge base interaction

## Usage Notes

1. All commands require the `*` prefix when used (e.g., `*help`, `*create-doc`)
2. Parameters in curly braces `{param}` are required
3. Parameters in square brackets `[param]` are optional
4. BMad Orchestrator serves as the entry point and can transform into any other agent
5. BMad Master has universal access to all tasks without persona transformation
6. Each agent has specific restrictions and permissions as noted in their descriptions

## Quick Start

1. Start with BMad Orchestrator using `*help` to see available workflows
2. Use `*agent [name]` to transform into a specialized agent
3. Execute commands specific to that agent's role
4. Use `*exit` to return to BMad Orchestrator or exit the session

## Agent Restrictions Summary

- **QA Architect**: Only updates "QA Results" section in story files
- **Developer**: Only modifies specific story sections (Tasks, Dev Records, File List, Status)
- **Scrum Master**: Cannot implement code or modify existing code
- **All Agents**: Must follow their specific role boundaries and command restrictions
