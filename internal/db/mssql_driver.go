package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
)

// MSSQLDriver implements Driver for Microsoft SQL Server databases.
// This is a fully independent implementation because MSSQL uses a distinct
// wire protocol, placeholder syntax (@p1, @p2, ...), identifier quoting
// ([brackets]), and T-SQL pagination (OFFSET/FETCH) that differ from all
// other supported backends.
type MSSQLDriver struct {
	db *sql.DB
}

func (d *MSSQLDriver) Open(ctx context.Context, dsn string) error {
	if after, ok := strings.CutPrefix(dsn, "mssql://"); ok {
		dsn = "sqlserver://" + after
	}
	dsn = strings.TrimPrefix(dsn, "sqlserver://")
	dsn = replaceLocalhost(dsn)
	var err error
	d.db, err = sql.Open("sqlserver", dsn)
	if err != nil {
		return fmt.Errorf("open mssql: %w", err)
	}
	return d.db.PingContext(ctx)
}

func (d *MSSQLDriver) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *MSSQLDriver) Kind() DataSourceKind { return KindMSSQL }

func (d *MSSQLDriver) ListTables(ctx context.Context) ([]string, error) {
	rows, err := d.db.QueryContext(ctx,
		`SELECT TABLE_NAME
		 FROM INFORMATION_SCHEMA.TABLES
		 WHERE TABLE_TYPE = 'BASE TABLE'
		   AND TABLE_SCHEMA = SCHEMA_NAME()
		 ORDER BY TABLE_NAME`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var tables []string
	for rows.Next() {
		var name string
		if rows.Scan(&name) == nil {
			tables = append(tables, name)
		}
	}
	return tables, nil
}

func (d *MSSQLDriver) LoadSchema(ctx context.Context, table string) ([]ColInfo, error) {
	rows, err := d.db.QueryContext(ctx,
		`SELECT ORDINAL_POSITION, COLUMN_NAME, DATA_TYPE,
		        IS_NULLABLE,
		        CASE WHEN COLUMNPROPERTY(OBJECT_ID(TABLE_SCHEMA + '.' + TABLE_NAME), COLUMN_NAME, 'IsIdentity') = 1 THEN 1 ELSE 0 END AS is_identity,
		        CASE WHEN EXISTS (
		            SELECT 1 FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
		            JOIN INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
				      ON kcu.CONSTRAINT_NAME = tc.CONSTRAINT_NAME
					     AND kcu.TABLE_SCHEMA = tc.TABLE_SCHEMA
		            WHERE kcu.TABLE_SCHEMA = SCHEMA_NAME()
		              AND kcu.TABLE_NAME = @p1
		              AND kcu.COLUMN_NAME = c.COLUMN_NAME
		              AND tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
		        ) THEN 1 ELSE 0 END AS is_pk,
		        COLUMN_DEFAULT
		 FROM INFORMATION_SCHEMA.COLUMNS c
		 WHERE c.TABLE_SCHEMA = SCHEMA_NAME()
		   AND c.TABLE_NAME = @p1
		 ORDER BY c.ORDINAL_POSITION`, table)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var cols []ColInfo
	for rows.Next() {
		var c ColInfo
		var ordinal int
		var nullable string
		var isPK, isIdentity int
		if rows.Scan(&ordinal, &c.Name, &c.Type, &nullable, &isIdentity, &isPK, &c.Dflt) == nil {
			c.NotNull = nullable == "NO"
			c.PK = isPK == 1
			c.CID = ordinal - 1
			_ = isIdentity
			cols = append(cols, c)
		}
	}
	return cols, nil
}

func (d *MSSQLDriver) LoadFKs(ctx context.Context, table string) ([]FKInfo, error) {
	rows, err := d.db.QueryContext(ctx,
		`SELECT fk.name AS constraint_name,
		        OBJECT_NAME(fk.referenced_object_id) AS foreign_table,
		        COL_NAME(fkc.parent_object_id, fkc.parent_column_id) AS column_name,
		        COL_NAME(fkc.referenced_object_id, fkc.referenced_column_id) AS foreign_column
		 FROM sys.foreign_keys fk
		 JOIN sys.foreign_key_columns fkc
		   ON fk.object_id = fkc.constraint_object_id
		 WHERE fk.parent_object_id = OBJECT_ID(SCHEMA_NAME() + '.' + @p1)
		 ORDER BY fk.name, fkc.constraint_column_id`, table)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var fks []FKInfo
	for rows.Next() {
		var f FKInfo
		if rows.Scan(&f.ID, &f.Table, &f.From, &f.To) == nil {
			fks = append(fks, f)
		}
	}
	return fks, nil
}

func (d *MSSQLDriver) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func (d *MSSQLDriver) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

func (d *MSSQLDriver) Placeholder(idx int) string {
	return fmt.Sprintf("@p%d", idx)
}

func (d *MSSQLDriver) QuoteIdent(name string) string {
	return fmt.Sprintf("[%s]", strings.ReplaceAll(name, "]", "]]"))
}

func (d *MSSQLDriver) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

func (d *MSSQLDriver) DB() *sql.DB { return d.db }

func (d *MSSQLDriver) LoadTableData(ctx context.Context, table string, page, pageSize int) (cols []string, rows [][]string, total int, err error) {
	schema, _ := d.LoadSchema(ctx, table)
	cols = ColNames(schema)
	colExpr := ColSelectExpr(schema)

	r, qerr := d.db.QueryContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", d.QuoteIdent(table)))
	if qerr == nil {
		r.Next()
		_ = r.Scan(&total)
		_ = r.Close()
	}

	offset := (page - 1) * pageSize
	q := fmt.Sprintf(
		"SELECT %s FROM %s ORDER BY (SELECT NULL) OFFSET %d ROWS FETCH NEXT %d ROWS ONLY",
		colExpr, d.QuoteIdent(table), offset, pageSize,
	)
	r, qerr = d.db.QueryContext(ctx, q)
	if qerr != nil {
		return nil, nil, 0, qerr
	}
	defer func() { _ = r.Close() }()

	realCols, _ := r.Columns()
	data, _, _ := ScanRows(r, len(realCols))
	return cols, data, total, nil
}

func (d *MSSQLDriver) RowCount(ctx context.Context, table string) (int, error) {
	var n int
	rows, err := d.db.QueryContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", d.QuoteIdent(table)))
	if err != nil {
		return 0, err
	}
	rows.Next()
	_ = rows.Scan(&n)
	_ = rows.Close()
	return n, nil
}

func (d *MSSQLDriver) LoadIndices(ctx context.Context, table string) ([]IndexInfo, error) {
	rows, err := d.db.QueryContext(ctx,
		`SELECT i.name AS index_name,
		        i.is_unique,
		        STRING_AGG(c.name, ',') WITHIN GROUP (ORDER BY ic.key_ordinal) AS columns
		 FROM sys.indexes i
		 JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
		 JOIN sys.columns c ON ic.object_id = c.object_id AND ic.column_id = c.column_id
		 WHERE i.object_id = OBJECT_ID(SCHEMA_NAME() + '.' + @p1)
		   AND i.is_primary_key = 0
		 GROUP BY i.name, i.is_unique
		 ORDER BY i.name`, table)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var indices []IndexInfo
	for rows.Next() {
		var idx IndexInfo
		var colList string
		if rows.Scan(&idx.Name, &idx.Unique, &colList) == nil {
			if colList != "" {
				idx.Columns = strings.Split(colList, ",")
			}
			indices = append(indices, idx)
		}
	}
	return indices, nil
}
