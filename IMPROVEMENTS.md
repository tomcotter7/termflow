1. **Add Tests** - No tests visible. At minimum, add tests for `formatTasks()`, `priorityOrdering()`, and storage operations.
2. **Use Interfaces for Storage** - `Handler` could implement an interface for easier testing/mocking.

## ✨ Feature Improvements

1. **Search/Filter** - Press `/` to filter tasks by keyword across all columns
2. **Quick Priority Adjust** - `+`/`-` keys to bump priority up/down without entering edit mode
3. **Task Templates** - Save common task structures for quick creation
4. **CLI Access** - Allow for adding/promoting/editing tasks/notes directly from the CLI for agents.
5. **Smart Defaults** - New tasks default to "today" priority if added in morning, "tomorrow" if added after 6pm
6. **Time Tracking Light** - Simple start/stop timer per task, no complexity
7. **Confirm Destructive Actions** - "Delete task? (y/n)" for delete operations

*Which of these is the best?* I would be worried about making super small tasks just to hit the counter/sumarry

8. **Weekly Summary** - Optional: "You completed 12 tasks this week" on Monday
9. **Streak Counter** - Days in a row with at least one completed task
10. **Velocity Hint** - Subtle indicator if you're completing more/fewer tasks than usual

## 🐛 Minor Issues

- The `showModetitleStyle` uses deprecated `.Copy()` - use `.Inherit()` or just create new styles
- `trimToLength()` in `showmode.go` is defined but never used
- Consider using `crypto/rand` with proper error handling or switch to `math/rand/v2` for non-crypto IDs

