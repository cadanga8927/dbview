package db

import (
	"context"
	"testing"
	"time"
)

func TestCassandraDriver_Kind(t *testing.T) {
	d := &CassandraDriver{}
	if d.Kind() != KindCassandra {
		t.Errorf("expected KindCassandra, got %v", d.Kind())
	}
}

func TestCassandraDriver_Placeholder(t *testing.T) {
	d := &CassandraDriver{}
	if d.Placeholder(1) != "?" {
		t.Errorf("expected ?, got %s", d.Placeholder(1))
	}
}

func TestCassandraDriver_QuoteIdent(t *testing.T) {
	d := &CassandraDriver{}
	got := d.QuoteIdent("my_table")
	if got != `"my_table"` {
		t.Errorf("expected \"my_table\", got %s", got)
	}
}

func TestCassandraDriver_Close_WhenNil(t *testing.T) {
	d := &CassandraDriver{}
	if err := d.Close(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCassandraDriver_DB(t *testing.T) {
	d := &CassandraDriver{}
	if d.DB() != nil {
		t.Error("expected nil")
	}
}

func TestCassandraDriver_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	d := &CassandraDriver{}
	err := d.Open(ctx, "cassandra://127.0.0.1:9042/testdb")
	if err != nil {
		t.Skipf("skipping: cassandra not available: %v", err)
	}
	defer func() { _ = d.Close() }()

	// Test Ping
	if err := d.Ping(ctx); err != nil {
		t.Errorf("ping: %v", err)
	}

	// Test ListTables
	tables, err := d.ListTables(ctx)
	if err != nil {
		t.Fatalf("list tables: %v", err)
	}
	if len(tables) < 3 {
		t.Errorf("expected at least 3 tables, got %d: %v", len(tables), tables)
	}
	t.Logf("tables: %v", tables)

	// Test LoadSchema for users
	schema, err := d.LoadSchema(ctx, "users")
	if err != nil {
		t.Fatalf("load schema: %v", err)
	}
	if len(schema) < 3 {
		t.Errorf("expected at least 3 columns, got %d", len(schema))
	}
	for _, c := range schema {
		t.Logf("  column: %s type=%s pk=%v", c.Name, c.Type, c.PK)
	}

	// Test RowCount
	count, err := d.RowCount(ctx, "users")
	if err != nil {
		t.Errorf("row count: %v", err)
	}
	if count != 25 {
		t.Errorf("expected 25 users, got %d", count)
	}

	// Test LoadTableData (page 1)
	cols, rows, total, err := d.LoadTableData(ctx, "users", 1, 10)
	if err != nil {
		t.Fatalf("load table data: %v", err)
	}
	if total != 25 {
		t.Errorf("expected total=25, got %d", total)
	}
	if len(rows) > 10 {
		t.Errorf("expected max 10 rows for page 1, got %d", len(rows))
	}
	t.Logf("columns: %v", cols)
	t.Logf("got %d rows (total=%d)", len(rows), total)
	if len(rows) > 0 {
		t.Logf("first row: %v", rows[0])
	}

	// Test LoadTableData (page 3 — should get remaining 5 rows)
	_, rows3, total3, err := d.LoadTableData(ctx, "users", 3, 10)
	if err != nil {
		t.Fatalf("load table data page 3: %v", err)
	}
	if total3 != 25 {
		t.Errorf("expected total=25, got %d", total3)
	}
	if len(rows3) != 5 {
		t.Errorf("expected 5 rows on page 3, got %d", len(rows3))
	}

	// Test LoadIndices
	indices, err := d.LoadIndices(ctx, "users")
	if err != nil {
		t.Logf("load indices: %v (non-fatal)", err)
	} else {
		t.Logf("indices: %d", len(indices))
	}

	// Test LoadFKs (should be nil)
	fks, err := d.LoadFKs(ctx, "users")
	if err != nil {
		t.Errorf("load fks: %v", err)
	}
	if len(fks) != 0 {
		t.Errorf("expected no FKs, got %d", len(fks))
	}

	// Test ExecuteQuery (SELECT)
	qCols, qRows, qAffected, err := d.ExecuteQuery(ctx, "SELECT * FROM testdb.users LIMIT 5")
	if err != nil {
		t.Fatalf("execute query: %v", err)
	}
	if len(qRows) != 5 {
		t.Errorf("expected 5 rows, got %d", len(qRows))
	}
	if qAffected != 5 {
		t.Errorf("expected affected=5, got %d", qAffected)
	}
	t.Logf("query columns: %v", qCols)

	// Test ExecuteQuery (TABLES)
	tCols, tRows, _, err := d.ExecuteQuery(ctx, "TABLES")
	if err != nil {
		t.Fatalf("execute TABLES: %v", err)
	}
	if len(tRows) < 3 {
		t.Errorf("expected at least 3 tables, got %d", len(tRows))
	}
	t.Logf("tables query: %v", tCols)
}
