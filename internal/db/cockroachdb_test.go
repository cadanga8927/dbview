package db

import (
	"testing"
)

func TestCockroachDBDriver_Kind(t *testing.T) {
	d := &CockroachDBDriver{}
	if got := d.Kind(); got != KindCockroachDB {
		t.Errorf("CockroachDBDriver.Kind() = %q, want %q", got, KindCockroachDB)
	}
}

func TestCockroachDBDriver_Placeholder(t *testing.T) {
	d := &CockroachDBDriver{}
	if got := d.Placeholder(1); got != "$1" {
		t.Errorf("CockroachDBDriver.Placeholder(1) = %q, want %q (delegates to PostgreSQL)", got, "$1")
	}
	if got := d.Placeholder(3); got != "$3" {
		t.Errorf("CockroachDBDriver.Placeholder(3) = %q, want %q (delegates to PostgreSQL)", got, "$3")
	}
}

func TestCockroachDBDriver_QuoteIdent(t *testing.T) {
	d := &CockroachDBDriver{}
	if got := d.QuoteIdent("table1"); got != `"table1"` {
		t.Errorf("CockroachDBDriver.QuoteIdent(%q) = %q, want %q (delegates to PostgreSQL)", "table1", got, `"table1"`)
	}
}
