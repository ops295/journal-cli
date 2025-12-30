package todo

import (
    "path/filepath"
    "strings"
    "testing"
    "time"

    "journal-cli/internal/fs"
)

func TestGetBacklog(t *testing.T) {
    dir := t.TempDir()
    // Create a sample previous day file
    md := `---
date: 2025-12-29
template: daily-human-dev
---

## ✅ Todos – Today
- [ ] Unchecked task
- [x] Checked task

## 🔁 Backlog
- [ ] Backlogged task
`

    path := filepath.Join(dir, "2025-12-29.md")
    if err := fs.WriteFile(path, []byte(md)); err != nil {
        t.Fatalf("failed to write test file: %v", err)
    }

    // Ensure fs.Exists returns true
    if !fs.Exists(path) {
        t.Fatalf("expected file to exist: %s", path)
    }

    items, err := GetBacklog(path)
    if err != nil {
        t.Fatalf("GetBacklog returned error: %v", err)
    }

    if len(items) != 2 {
        t.Fatalf("expected 2 backlog items (1 todo + 1 backlog), got %d", len(items))
    }

    // Basic content checks
    found := map[string]bool{}
    for _, it := range items {
        found[it.Text] = true
    }
    if !found["Unchecked task"] || !found["Backlogged task"] {
        t.Fatalf("unexpected backlog items: %v", found)
    }
}

func TestGetPreviousJournalPath(t *testing.T) {
    base := "/tmp/journal"
    // Use a known date
    // Previous path should end with 2025-12-29.md
    // We don't assert separator specifics, just suffix
    got := GetPreviousJournalPath(base, fsTime())
    if !strings.HasSuffix(got, "2025-12-29.md") {
        t.Fatalf("unexpected previous path: %s", got)
    }
}

// fsTime returns a fixed date used by tests
func fsTime() time.Time {
    return time.Date(2025, 12, 30, 0, 0, 0, 0, time.UTC)
}
