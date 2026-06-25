package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadMigration(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "00001_test.sql")
	content := "-- +goose Up\nCREATE TABLE test(id int);\n-- +goose Down\nDROP TABLE test;"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	upSQL, downSQL, err := readMigration(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(upSQL, "CREATE TABLE") {
		t.Fatalf("unexpected up SQL: %s", upSQL)
	}
	if !strings.Contains(downSQL, "DROP TABLE") {
		t.Fatalf("unexpected down SQL: %s", downSQL)
	}
}
