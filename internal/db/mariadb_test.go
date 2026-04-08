package db

import (
	"testing"
)

func TestMariaDBDriver_Kind(t *testing.T) {
	d := &MariaDBDriver{}
	if got := d.Kind(); got != KindMariaDB {
		t.Errorf("MariaDBDriver.Kind() = %q, want %q", got, KindMariaDB)
	}
}

func TestMariaDBDriver_Placeholder(t *testing.T) {
	d := &MariaDBDriver{}
	if got := d.Placeholder(1); got != "?" {
		t.Errorf("MariaDBDriver.Placeholder(1) = %q, want %q (delegates to MySQL)", got, "?")
	}
}

func TestMariaDBDriver_QuoteIdent(t *testing.T) {
	d := &MariaDBDriver{}
	if got := d.QuoteIdent("table1"); got != "`table1`" {
		t.Errorf("MariaDBDriver.QuoteIdent(%q) = %q, want %q (delegates to MySQL)", "table1", got, "`table1`")
	}
}
