package fs

import (
    "os"
    "path/filepath"
    "testing"
)

func TestEnsureWriteReadExists(t *testing.T) {
    dir := t.TempDir()
    d := filepath.Join(dir, "subdir")
    if err := EnsureDir(d); err != nil {
        t.Fatalf("EnsureDir failed: %v", err)
    }

    p := filepath.Join(d, "file.txt")
    data := []byte("hello world")
    if err := WriteFile(p, data); err != nil {
        t.Fatalf("WriteFile failed: %v", err)
    }

    got, err := ReadFile(p)
    if err != nil {
        t.Fatalf("ReadFile failed: %v", err)
    }

    if string(got) != string(data) {
        t.Fatalf("content mismatch: %s", string(got))
    }

    if !Exists(p) {
        t.Fatalf("Exists should be true for %s", p)
    }

    // Non-existent file
    if Exists(filepath.Join(dir, "nope.txt")) {
        t.Fatalf("Exists should be false for non-existent file")
    }

    // cleanup test file
    os.Remove(p)
}
