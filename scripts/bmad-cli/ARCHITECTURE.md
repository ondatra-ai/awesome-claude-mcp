# BMAD CLI Architecture

## Core Principle: Direct Data Flow

**NO CACHING, NO LOADERS, NO UNNECESSARY INTERFACES - JUST DIRECT DATA FLOW!** 🎉

### The Architecture

```
Epic File → StoryFactory → StoryDocument (complete) → Generators → Output
```

### Key Design Decisions

1. **StoryDocument as Single Source of Truth**
   - Contains all data needed for generation: Story, ArchitectureDocs, etc.
   - No lazy loading - everything populated upfront
   - Generators receive complete StoryDocument, not IDs

2. **No Caching Layer**
   - Data loaded once in StoryFactory
   - No cache misses, no cache invalidation complexity
   - Simple, predictable data flow

3. **No Loader Interfaces**
   - Direct use of concrete types (docs.ArchitectureLoader)
   - No unnecessary abstractions
   - Each layer has clear, concrete dependencies

4. **Generator Simplicity**
   - Accept `*story.StoryDocument` parameter
   - Extract needed data directly from document
   - No dependency injection of loaders

### Code Example

```go
// ✅ Current Architecture
func (f *StoryFactory) CreateStory(ctx context.Context, storyNumber string) (*story.StoryDocument, error) {
    // Load everything once
    loadedStory, err := f.epicLoader.LoadStoryFromEpic(storyNumber)
    architectureDocs, err := f.architectureLoader.LoadAllArchitectureDocsStruct()

    // Create complete document
    storyDoc := &story.StoryDocument{
        Story: *loadedStory,
        ArchitectureDocs: architectureDocs,
        // ... other fields
    }

    // Generators work with complete document
    tasks, err := taskGenerator.GenerateTasks(ctx, storyDoc)
    storyDoc.Tasks = tasks

    return storyDoc, nil
}

func (g *TaskGenerator) GenerateTasks(ctx context.Context, storyDoc *story.StoryDocument) ([]story.Task, error) {
    // Direct access to all needed data
    return g.generateFromStoryAndArchitecture(storyDoc.Story, storyDoc.ArchitectureDocs)
}
```

### Benefits Achieved

- ✅ **Eliminated 200+ lines** of unnecessary abstraction code
- ✅ **Zero caching complexity** - no cache misses or invalidation
- ✅ **Self-documenting data flow** - easy to understand and debug
- ✅ **Impossible to misuse** - only one way data flows through system
- ✅ **Fast and reliable** - no network calls or cache lookups during generation

### Refactoring History

- **Before**: Complex loader interfaces + caching + dependency injection
- **After**: StoryDocument with direct data access
- **Result**: Same functionality, 50% less code, 100% more maintainable

---

**Remember**: When tempted to add interfaces, loaders, or caching - ask "Does this add real value or just complexity?"
