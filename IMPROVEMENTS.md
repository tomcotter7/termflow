## 🎨 Appearance Improvements

1. [X] **Task Information Density** - Show due dates and priority indicators (e.g., `!`, `!!`, `!!!`) inline with task names in the board view
2. [X] **Column Task Counts** - Add task counts to column headers: `todo (5)`, `inprogress (2)`
3. **Visual Priority** - Use color gradients or icons to distinguish priority levels at a glance
4. **Progress Bar** - Add a simple progress indicator showing done/total tasks
5. **Better Help Layout** - Group shortcuts by category (navigation, actions, modes) in a cleaner grid format
6. **Completion Dates** - Show when tasks were completed in the done column

## 🏗️ Codebase Improvements

1. [X] **Extract Constants** - Status strings (`"todo"`, `"inprogress"`, etc.) are scattered as magic strings. Create a `const` block:
   ```go
   const (
       StatusTodo       = "todo"
       StatusInProgress = "inprogress"
       StatusInReview   = "in-review"
       StatusDone       = "done"
   )
   ```

2. [X] **Handle `rand.Read` Error** - In `editmode.go`, `randomId()` ignores the error:
   ```go
   rand.Read(b) // error ignored!
   ```

3. [X] **Unify Form Types** - `Form` and `TextInputs` structs have duplicated logic for focus management. Consider a single abstraction.

4. **Split Large Functions** - `normalModeView()` is ~150 lines. Extract helpers like `renderTaskCell()`, `renderHeaders()`, `renderFooter()`.

5. **Add Tests** - No tests visible. At minimum, add tests for `formatTasks()`, `priorityOrdering()`, and storage operations.

6. **Use Interfaces for Storage** - `Handler` could implement an interface for easier testing/mocking.

7. [X] **Consistent Error Handling** - Some places set `m.err` and switch to `ErrorMode`, others use `log.Fatal`. Standardize the approach.

## ✨ Feature Improvements

1. **Search/Filter** - Press `/` to filter tasks by keyword across all columns
2. **Quick Priority Adjust** - `+`/`-` keys to bump priority up/down without entering edit mode
3. **Undo** - Track last action and allow `u` to undo delete/promote/regress
5. **Vim Navigation** - 
    - [X] `gg`/`G` for top/bottom
    - [X] `0`/`$` for first/last column
6. **Tags/Labels** - Support `#tag` syntax in descriptions for categorization
7. **Archive** - Instead of delete, archive completed tasks for historical reference
8. **View .plan Files** - Command to view previously generated .plan files from within the app
9. **Task Templates** - Save common task structures for quick creation
10. **Export** - Export board to markdown/CSV for sharing

## 🐛 Minor Issues

- The `showModetitleStyle` uses deprecated `.Copy()` - use `.Inherit()` or just create new styles
- `trimToLength()` in `showmode.go` is defined but never used
- Consider using `crypto/rand` with proper error handling or switch to `math/rand/v2` for non-crypto IDs

