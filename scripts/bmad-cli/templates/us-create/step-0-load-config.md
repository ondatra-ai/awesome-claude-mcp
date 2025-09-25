<!-- Powered by BMAD™ Core -->

# Step 0: Load Configuration and Check Workflow

## Purpose

Load and validate the core project configuration required for story creation. This step establishes the foundational context needed for all subsequent steps.

## Input Context

```yaml
input:
  epicStoryId: "3.1"  # Format: {epicNum}.{storyNum}
  options: {}  # Additional command-line options
```

## Process

### 1. Load Core Configuration

- **Target File:** `.bmad-core/core-config.yaml`
- **Required Check:** Verify file exists in project root
- **Failure Action:** If missing, HALT with error message:
  ```
  ERROR: core-config.yaml not found. This file is required for story creation.

  Solutions:
  1) Copy from GITHUB bmad-core/core-config.yaml and configure for your project
  2) Run the BMad installer against your project to upgrade automatically

  Please add and configure core-config.yaml before proceeding.
  ```

### 2. Extract Key Configurations

Extract and validate the following required configurations:

#### Development Settings
- `devStoryLocation`: Where story files are created (e.g., "scripts/bmad-cli/templates")
- `prdSharded`: Boolean indicating if PRD is sharded across multiple files
- `prd.*`: PRD-related configuration

#### Architecture Settings
- `architectureVersion`: Version of architecture documentation (check for >= v4)
- `architectureSharded`: Boolean indicating if architecture is sharded
- `architectureFile`: Monolithic architecture file path (if not sharded)
- `architectureShardedLocation`: Directory containing sharded architecture files

#### Workflow Settings
- `workflow.*`: Any workflow-specific configurations

### 3. Parse Target Story

From input `epicStoryId` (e.g., "3.1"):
- Extract `epicNum` (e.g., 3)
- Extract `storyNum` (e.g., 1)
- Validate format matches pattern: `\d+\.\d+`

### 4. Validate Configuration Integrity

Check that all required paths exist and are accessible:
- Verify `devStoryLocation` directory exists or can be created
- If `prdSharded` is true, verify PRD location exists
- If `architectureSharded` is true, verify `architectureShardedLocation` exists
- If `architectureSharded` is false, verify `architectureFile` exists

## Output Context

```yaml
context:
  config:
    devStoryLocation: "scripts/bmad-cli/templates"
    prdSharded: true
    prdLocation: "docs/epics/jsons"
    architectureVersion: "v4"
    architectureSharded: true
    architectureShardedLocation: "docs/architecture"
    # ... additional config fields

  target:
    epicNum: 3
    storyNum: 1
    epicStoryId: "3.1"

  validation:
    configValid: true
    pathsAccessible: true
    errors: []
    warnings: []
```

## Error Handling

### Missing Configuration File
- **Action:** HALT execution
- **Message:** Provide clear instructions for obtaining core-config.yaml

### Invalid Configuration Values
- **Action:** HALT execution with validation details
- **Message:** List specific configuration issues and required fixes

### Inaccessible Paths
- **Action:** HALT execution
- **Message:** List which required paths are missing or inaccessible

## Success Criteria

- Core configuration loaded and validated
- Target story identifiers parsed correctly
- All required paths verified as accessible
- Configuration context prepared for next step

## Next Step

On success, pass context to `step-1-identify-story.md`
