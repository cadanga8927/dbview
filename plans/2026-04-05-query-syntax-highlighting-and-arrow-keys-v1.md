# Query View: Syntax Highlighting & Arrow Key Improvements

## Objective

Add SQL syntax highlighting to the query input field and ensure arrow key navigation (left/right for cursor movement, up/down for history) works correctly in all edge cases — empty history, beginning/end of line, and beginning/end of history.

## Implementation Plan

### Phase 1: SQL Syntax Highlighting

- [x] **1.1 Create a SQL highlighter function** in a new file `internal/highlight/sql.go`. Implement a simple tokenizer that splits SQL text into tokens and classifies them as:
  - **Keywords** (SELECT, FROM, WHERE, INSERT, UPDATE, DELETE, DROP, ALTER, CREATE, TABLE, INDEX, AND, OR, NOT, IN, LIKE, JOIN, ON, AS, ORDER, BY, GROUP, HAVING, LIMIT, OFFSET, SET, VALUES, INTO, NULL, IS, DISTINCT, COUNT, SUM, AVG, MIN, MAX, EXISTS, BETWEEN, CASE, WHEN, THEN, ELSE, END, UNION, ALL, PRIMARY, KEY, FOREIGN, REFERENCES, CASCADE, INNER, LEFT, RIGHT, OUTER, FULL, CROSS, NATURAL, ASC, DESC, IF, BEGIN, COMMIT, ROLLBACK, TRANSACTION, ATTACH, DETACH, REPLACE, TRUNCATE, WITH, EXPLAIN, PRAGMA, VACUUM, ANALYZE)
  - **Strings** (single-quoted `'...'`)
  - **Numbers** (integer and float literals)
  - **Operators** (`=`, `<>`, `!=`, `<`, `>`, `<=`, `>=`, `+`, `-`, `*`, `/`, `%`, `||`, `(`, `)`, `,`, `;`, `.`)
  - **Identifiers** (everything else — table/column names)
  - **Comments** (`--` single-line comments)

- [x] **1.2 Define highlight color mapping** using the existing `theme.Colors` struct. The mapping:
  - Keywords → `cl.Accent` (purple/blue depending on theme)
  - Strings → `cl.Ok` (green)
  - Numbers → `cl.Warn` (yellow/gold)
  - Operators → `cl.Dim` (gray)
  - Comments → `cl.Dim` (gray)
  - Identifiers → `cl.White` (default text color)

- [x] **1.3 Handle the cursor position correctly** within highlighted text. Cursor is inserted after tokenization.

- [x] **1.4 Update `renderQuery()` in `internal/app/view.go`** to use the highlighter.

- [x] **1.5 Handle multi-line or very long queries** gracefully.

### Phase 2: Arrow Key Enhancements

- [x] **2.1 Add `home` key support** in `updateQuery()`.

- [x] **2.2 Add `end` key support** in `updateQuery()`.

- [x] **2.3 Add `ctrl+left` for word-backward navigation** in `updateQuery()`.

- [x] **2.4 Add `ctrl+right` for word-forward navigation** in `updateQuery()`.

- [x] **2.5 Verify up/down history behavior** at all boundaries — no changes needed, already correct.

- [x] **2.6 Add `delete` (forward delete) key support** in `updateQuery()`.

- [x] **2.7 Add `ctrl+a` (home) and `ctrl+e` (end) support** for Emacs-style line editing.

### Phase 3: Driver-Aware Highlighting

- [x] **3.1 Add MongoDB query keyword highlighting** in the highlighter.

- [x] **3.2 Add Redis command highlighting** in the highlighter.

- [x] **3.3 Make the highlighter driver-aware** by passing the driver kind to the highlight function.

## Files Modified

| File | Change |
|------|--------|
| `internal/highlight/sql.go` (new) | SQL tokenizer and highlighter function with driver-aware keyword sets |
| `internal/app/view.go:484-493` | Replace plain-text query rendering with highlighted rendering |
| `internal/app/update.go:96-139` | Add helper functions: deleteRuneAtPos, wordBackward, wordForward |
| `internal/app/update.go:1093-1106` | Add home, end, ctrl+left, ctrl+right, delete, ctrl+a, ctrl+e key handlers |
