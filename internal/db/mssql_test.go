package db

import (
	"testing"
)

func TestMSSQLDriver_Kind(t *testing.T) {
	d := &MSSQLDriver{}
	if got := d.Kind(); got != KindMSSQL {
		t.Errorf("MSSQLDriver.Kind() = %q, want %q", got, KindMSSQL)
	}
}

func TestMSSQLDriver_Placeholder(t *testing.T) {
	d := &MSSQLDriver{}
	tests := []struct {
		idx  int
		want string
	}{
		{1, "@p1"},
		{2, "@p2"},
		{3, "@p3"},
		{10, "@p10"},
	}
	for _, tt := range tests {
		got := d.Placeholder(tt.idx)
		if got != tt.want {
			t.Errorf("Placeholder(%d) = %q, want %q", tt.idx, got, tt.want)
		}
	}
}

func TestMSSQLDriver_QuoteIdent(t *testing.T) {
	d := &MSSQLDriver{}
	tests := []struct {
		input string
		want  string
	}{
		{"table1", "[table1]"},
		{"column_name", "[column_name]"},
		{"has]bracket", "[has]]bracket]"},
		{"multiple]]brackets", "[multiple]]]]brackets]"},
		{"schema.tbl", "[schema.tbl]"},
	}
	for _, tt := range tests {
		got := d.QuoteIdent(tt.input)
		if got != tt.want {
			t.Errorf("QuoteIdent(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMSSQLDriver_Close_WhenNil(t *testing.T) {
	d := &MSSQLDriver{}
	if err := d.Close(); err != nil {
		t.Errorf("Close() on nil db should not error, got: %v", err)
	}
}

func TestMSSQLDriver_DB_WhenNil(t *testing.T) {
	d := &MSSQLDriver{}
	if got := d.DB(); got != nil {
		t.Errorf("DB() on nil db should return nil, got %v", got)
	}
}
