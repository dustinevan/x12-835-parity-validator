package db

import "testing"

func TestHighestMigrationVersion(t *testing.T) {
	version, err := HighestMigrationVersion()
	if err != nil {
		t.Fatal(err)
	}
	if version != 1 {
		t.Fatalf("highest migration version: got %d, want 1", version)
	}
}

func TestMigrationVersion(t *testing.T) {
	tests := []struct {
		sql     string
		version int
		ok      bool
	}{
		{"-- Version: 12\n-- Description: x\nSELECT 1;", 12, true},
		{"--version:3\n", 3, true},
		{"-- Description: no version\n", 0, false},
		{"-- Version: abc\n", 0, false},
	}
	for _, test := range tests {
		version, ok := migrationVersion(test.sql)
		if version != test.version || ok != test.ok {
			t.Errorf("migrationVersion(%q): got %d, %v; want %d, %v", test.sql, version, ok, test.version, test.ok)
		}
	}
}
