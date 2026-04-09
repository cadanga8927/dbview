package db

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"gopkg.in/inf.v0"
)

// CassandraDriver implements Driver for Apache Cassandra / ScyllaDB.
// It maps keyspaces to databases and CQL tables to the TUI table view.
type CassandraDriver struct {
	session  *gocql.Session
	keyspace string
	cluster  *gocql.ClusterConfig
}

func (d *CassandraDriver) Open(ctx context.Context, dsn string) error {
	// Parse cassandra://user:pass@host:port/keyspace
	u := dsn
	if after, ok := strings.CutPrefix(u, "cassandra://"); ok {
		u = after
	}

	var host, keyspace, username, password string
	port := 9042

	// Extract user:pass@ if present
	if atIdx := strings.Index(u, "@"); atIdx >= 0 {
		userPass := u[:atIdx]
		hostPart := u[atIdx+1:]
		if colonIdx := strings.Index(userPass, ":"); colonIdx >= 0 {
			username = userPass[:colonIdx]
			password = userPass[colonIdx+1:]
		} else {
			username = userPass
		}
		u = hostPart
	}

	// Extract keyspace from path
	if slashIdx := strings.Index(u, "/"); slashIdx >= 0 {
		keyspace = u[slashIdx+1:]
		u = u[:slashIdx]
	}

	// Extract host:port
	host = u
	if colonIdx := strings.Index(host, ":"); colonIdx >= 0 {
		portStr := host[colonIdx+1:]
		host = host[:colonIdx]
		if p, err := parsePort(portStr); err == nil {
			port = p
		}
	}

	if host == "" {
		host = "127.0.0.1"
	}

	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second
	cluster.DisableInitialHostLookup = true

	if username != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: password,
		}
	}

	if keyspace == "" {
		return fmt.Errorf("cassandra: keyspace is required in DSN (cassandra://host:port/keyspace)")
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("cassandra connect %s:%d/%s: %w", host, port, keyspace, err)
	}

	d.session = session
	d.keyspace = keyspace
	d.cluster = cluster
	return nil
}

func (d *CassandraDriver) Close() error {
	if d.session != nil {
		d.session.Close()
	}
	return nil
}

func (d *CassandraDriver) Kind() DataSourceKind { return KindCassandra }

func (d *CassandraDriver) ListTables(ctx context.Context) ([]string, error) {
	query := "SELECT table_name FROM system_schema.tables WHERE keyspace_name = ?"
	iter := d.session.Query(query, d.keyspace).Iter()
	var name string
	var tables []string
	for iter.Scan(&name) {
		tables = append(tables, name)
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("list tables: %w", err)
	}
	sort.Strings(tables)
	return tables, nil
}

func (d *CassandraDriver) LoadSchema(ctx context.Context, table string) ([]ColInfo, error) {
	query := `SELECT column_name, kind, type FROM system_schema.columns WHERE keyspace_name = ? AND table_name = ?`
	iter := d.session.Query(query, d.keyspace, table).Iter()

	var colName, kind, typ string
	var cols []ColInfo
	var pkCols []ColInfo
	cid := 0
	for iter.Scan(&colName, &kind, &typ) {
		c := ColInfo{
			CID:     cid,
			Name:    colName,
			Type:    typ,
			NotNull: kind == "partition_key",
			PK:      kind == "partition_key" || kind == "clustering",
		}
		if kind == "partition_key" || kind == "clustering" {
			pkCols = append(pkCols, c)
		} else {
			cols = append(cols, c)
		}
		cid++
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("load schema: %w", err)
	}

	// Partition keys first, then clustering columns, then regular columns
	result := append(pkCols, cols...)
	return result, nil
}

func (d *CassandraDriver) LoadFKs(ctx context.Context, table string) ([]FKInfo, error) {
	return nil, nil // Cassandra has no foreign keys
}

func (d *CassandraDriver) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, fmt.Errorf("use CassandraQuery for Cassandra data access")
}

// CassandraQuery executes a CQL query and returns column names and string rows.
func (d *CassandraDriver) CassandraQuery(ctx context.Context, query string, args ...interface{}) (cols []string, rows [][]string, err error) {
	q := d.session.Query(query, args...)
	iter := q.Iter()

	colInfos := iter.Columns()
	cols = make([]string, len(colInfos))
	for i, ci := range colInfos {
		cols[i] = ci.Name
	}
	colCount := len(cols)
	if colCount == 0 {
		return nil, nil, nil
	}

	// gocql requires type-specific scan destinations. Build dest slice from column TypeInfo.
	dest := make([]interface{}, colCount)
	for i, ci := range colInfos {
		dest[i] = newScanDest(ci.TypeInfo)
	}

	for iter.Scan(dest...) {
		row := make([]string, colCount)
		for i := range dest {
			row[i] = formatScanValue(dest[i])
		}
		rows = append(rows, row)
	}

	if cerr := iter.Close(); cerr != nil {
		return nil, nil, fmt.Errorf("cassandra query: %w", cerr)
	}

	return cols, rows, nil
}

func (d *CassandraDriver) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, fmt.Errorf("use CassandraExec for Cassandra mutations")
}

// CassandraExec executes a CQL mutation and returns a simple result.
func (d *CassandraDriver) CassandraExec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	if err := d.session.Query(query, args...).Exec(); err != nil {
		return 0, err
	}
	return 1, nil
}

func (d *CassandraDriver) Placeholder(idx int) string { return "?" }

func (d *CassandraDriver) QuoteIdent(name string) string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(name, `"`, `""`))
}

func (d *CassandraDriver) Ping(ctx context.Context) error {
	if d.session == nil {
		return fmt.Errorf("not connected")
	}
	return d.session.Query("SELECT now() FROM system.local").Exec()
}

func (d *CassandraDriver) DB() *sql.DB { return nil }

func (d *CassandraDriver) LoadTableData(ctx context.Context, table string, page, pageSize int) (cols []string, rows [][]string, total int, err error) {
	schema, _ := d.LoadSchema(ctx, table)
	cols = ColNames(schema)

	// Count total rows
	var count int
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", d.QuoteIdent(d.keyspace), d.QuoteIdent(table))
	iter := d.session.Query(countQ).Iter()
	for iter.Scan(&count) {
	}
	_ = iter.Close()
	total = count

	// Fetch page — Cassandra uses LIMIT only (no OFFSET), so we fetch all up to the needed range
	// and skip the first (page-1)*pageSize rows
	offset := (page - 1) * pageSize
	limit := offset + pageSize

	q := fmt.Sprintf("SELECT %s FROM %s.%s LIMIT %d",
		strings.Join(cols, ", "),
		d.QuoteIdent(d.keyspace), d.QuoteIdent(table), limit)

	queryCols, allRows, qerr := d.CassandraQuery(ctx, q)
	if qerr != nil {
		return nil, nil, 0, qerr
	}

	// Update cols from actual query if they differ
	if len(queryCols) > 0 {
		cols = queryCols
	}

	if offset >= len(allRows) {
		return cols, nil, total, nil
	}

	return cols, allRows[offset:], total, nil
}

func (d *CassandraDriver) RowCount(ctx context.Context, table string) (int, error) {
	var count int
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", d.QuoteIdent(d.keyspace), d.QuoteIdent(table))
	iter := d.session.Query(q).Iter()
	for iter.Scan(&count) {
	}
	_ = iter.Close()
	return count, nil
}

func (d *CassandraDriver) LoadIndices(ctx context.Context, table string) ([]IndexInfo, error) {
	// Cassandra indexes
	var indices []IndexInfo

	// Secondary indexes
	query := "SELECT index_name, options FROM system_schema.indexes WHERE keyspace_name = ? AND table_name = ?"
	iter := d.session.Query(query, d.keyspace, table).Iter()
	var name, options string
	for iter.Scan(&name, &options) {
		indices = append(indices, IndexInfo{
			Name:   name,
			Unique: false,
		})
	}
	_ = iter.Close()
	return indices, nil
}

// InsertRow inserts one row into a Cassandra table.
func (d *CassandraDriver) InsertRow(ctx context.Context, table string, cols []string, vals []string) (int64, error) {
	if len(cols) == 0 {
		return 0, fmt.Errorf("no columns provided")
	}

	quotedCols := make([]string, len(cols))
	placeholders := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	for i, col := range cols {
		quotedCols[i] = d.QuoteIdent(col)
		placeholders[i] = "?"
		val := ""
		if i < len(vals) {
			val = vals[i]
		}
		args[i] = parseCassandraValue(val)
	}

	q := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s)",
		d.QuoteIdent(d.keyspace), d.QuoteIdent(table),
		strings.Join(quotedCols, ", "), strings.Join(placeholders, ", "))

	if err := d.session.Query(q, args...).Exec(); err != nil {
		return 0, err
	}
	return 1, nil
}

// ExecuteQuery supports CQL commands for the query view.
func (d *CassandraDriver) ExecuteQuery(ctx context.Context, query string) (cols []string, rows [][]string, affected int64, err error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return nil, nil, 0, fmt.Errorf("empty query")
	}

	upper := strings.ToUpper(q)

	// DESCRIBE/SHOW tables
	if upper == "TABLES" || upper == "SHOW TABLES" || upper == "DESCRIBE TABLES" {
		names, lerr := d.ListTables(ctx)
		if lerr != nil {
			return nil, nil, 0, lerr
		}
		out := make([][]string, 0, len(names))
		for _, n := range names {
			out = append(out, []string{n})
		}
		return []string{"table"}, out, int64(len(out)), nil
	}

	// SELECT queries
	if strings.HasPrefix(upper, "SELECT") {
		cols, rows, qerr := d.CassandraQuery(ctx, q)
		if qerr != nil {
			return nil, nil, 0, qerr
		}
		return cols, rows, int64(len(rows)), nil
	}

	// Mutation queries (INSERT, UPDATE, DELETE, TRUNCATE, etc.)
	n, execErr := d.CassandraExec(ctx, q)
	if execErr != nil {
		return nil, nil, 0, execErr
	}
	return nil, nil, n, nil
}

func parseCassandraValue(raw string) interface{} {
	v := strings.TrimSpace(raw)
	if v == "" || strings.EqualFold(v, "NULL") {
		return nil
	}
	if strings.EqualFold(v, "true") {
		return true
	}
	if strings.EqualFold(v, "false") {
		return false
	}

	// Try parsing as UUID-like string or plain string — let gocql handle type conversion
	// via its own marshaling when binding parameters.
	return v
}

// newScanDest returns a pointer to the correct Go type for a given gocql TypeInfo.
func newScanDest(ti gocql.TypeInfo) interface{} {
	switch ti.Type() {
	case gocql.TypeBoolean:
		return new(bool)
	case gocql.TypeTinyInt:
		return new(int8)
	case gocql.TypeSmallInt:
		return new(int16)
	case gocql.TypeInt:
		return new(int)
	case gocql.TypeBigInt, gocql.TypeCounter:
		return new(int64)
	case gocql.TypeFloat:
		return new(float32)
	case gocql.TypeDouble:
		return new(float64)
	case gocql.TypeDecimal:
		return new(inf.Dec)
	case gocql.TypeTimestamp:
		return new(time.Time)
	case gocql.TypeDate:
		return new(time.Time)
	case gocql.TypeBlob:
		return new([]byte)
	case gocql.TypeUUID, gocql.TypeTimeUUID:
		return new(gocql.UUID)
	case gocql.TypeInet:
		return new(string)
	default:
		return new(string)
	}
}

// formatScanValue converts a scanned gocql value pointer to its string representation.
func formatScanValue(ptr interface{}) string {
	switch v := ptr.(type) {
	case *bool:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("%t", *v)
	case *int8:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("%d", *v)
	case *int16:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("%d", *v)
	case *int:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("%d", *v)
	case *int64:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("%d", *v)
	case *float32:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("%v", *v)
	case *float64:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("%v", *v)
	case *inf.Dec:
		if v == nil {
			return "NULL"
		}
		return v.String()
	case *time.Time:
		if v == nil {
			return "NULL"
		}
		return v.Format(time.RFC3339)
	case *[]byte:
		if v == nil || *v == nil {
			return "NULL"
		}
		return FormatValue(*v)
	case *gocql.UUID:
		if v == nil {
			return "NULL"
		}
		return v.String()
	case *string:
		if v == nil {
			return "NULL"
		}
		return *v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func parsePort(s string) (int, error) {
	var port int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			port = port*10 + int(c-'0')
		} else {
			return 0, fmt.Errorf("invalid port: %s", s)
		}
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port out of range: %d", port)
	}
	return port, nil
}
