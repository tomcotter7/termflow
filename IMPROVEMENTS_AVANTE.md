# termflow Improvement Ideas

> Focus: Making termflow **dead simple** to use

---


## ⌨️ Simpler Keybindings

2. **Number Keys for Priority** - `1`, `2`, `3` to set priority instantly without entering edit mode
3. **Double-tap Escape** - Single `Esc` cancels current action, double-tap `Esc` quits app (prevents accidental exits)

---

## 🎯 Simplified Workflow

1. **Smart Defaults** - New tasks default to "today" priority if added in morning, "tomorrow" if added after 6pm
2. **Auto-Archive** - Tasks in "done" for 7+ days automatically archive (configurable)
- [X] **Inbox Column** - Optional "inbox" column for quick brain dumps, sort later
4. **Batch Operations** - Select multiple tasks with `v` (visual mode) then promote/delete all at once

---

## 👀 Cleaner Interface

1. **Minimal Mode** - Hide all chrome, just show tasks (toggle with `m`)
2. **Focus Mode** - Dim all columns except current one
3. **Compact View** - Single-line task display for boards with many tasks

---

## 🔧 Zero-Config Experience

1. **Portable Mode** - `termflow --portable` stores everything in current directory

---

## 📱 Quick Actions

1. **Fuzzy Task Jump** - Press `/` and type to fuzzy-find any task across all columns
2. **Time Tracking Light** - Simple start/stop timer per task, no complexity

---

## 🔄 Natural Language Input

1. **Smart Parsing** - `"Fix bug tomorrow !!"` auto-sets due date and priority
- [X] **Relative Dates** - Support "today", "tomorrow", "next week", "friday"
4. **Tags Inline** - `#work` or `@john` parsed automatically from task text

---

## 💾 Data & Sync

1. **Plain Text Storage** - Human-readable YAML/JSON so users can edit manually if needed
2. **Git-Friendly** - Format that diffs well for version control


## 🐛 Error Prevention

1. **Confirm Destructive Actions** - "Delete task? (y/n)" for delete operations
2. **Undo Last Action** - `u` to undo (at minimum last delete)
3. **Task Recovery** - Deleted tasks go to trash, recoverable for 24 hours
4. **Auto-Backup** - Silent daily backup of task data

---

## 📊 Lightweight Analytics

1. **Weekly Summary** - Optional: "You completed 12 tasks this week" on Monday
2. **Streak Counter** - Days in a row with at least one completed task
3. **Velocity Hint** - Subtle indicator if you're completing more/fewer tasks than usual

---
