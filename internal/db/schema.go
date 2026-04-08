package db

import (
	"time"
)

// DataSourceKind identifies the type of database backend.
type DataSourceKind string

const (
	KindSQLite      DataSourceKind = "sqlite"
	KindMySQL       DataSourceKind = "mysql"
	KindMariaDB     DataSourceKind = "mariadb"
	KindPostgreSQL  DataSourceKind = "postgresql"
	KindCockroachDB DataSourceKind = "cockroachdb"
	KindMSSQL       DataSourceKind = "mssql"
	KindMongoDB     DataSourceKind = "mongodb"
	KindRedis       DataSourceKind = "redis"
	KindCassandra   DataSourceKind = "cassandra"
)

// TableKind identifies the data model of a table/collection.
type TableKind string

const (
	TableKindRelational TableKind = "relational" // SQL tables
	TableKindDocument   TableKind = "document"   // MongoDB collections
	TableKindKeyValue   TableKind = "keyvalue"   // Redis key groups
)

// FieldKind identifies the nature of a field/column.
type FieldKind string

const (
	FieldKindColumn FieldKind = "column" // SQL column
	FieldKindField  FieldKind = "field"  // MongoDB document field
	FieldKindKey    FieldKind = "key"    // Redis key component
	FieldKindValue  FieldKind = "value"  // Redis value
	FieldKindNested FieldKind = "nested" // nested object/array
)

// DataSource represents the connected database.
type DataSource struct {
	Kind DataSourceKind
	Name string // display name (filename, db name, etc.)
}

// TableSchema describes a table/collection/keyspace.
type TableSchema struct {
	Name        string
	Kind        TableKind
	RowCount    int
	SizeBytes   int64
	Fields      []FieldDefinition
	ForeignKeys []ForeignKeyDefinition
	Indexes     []IndexDefinition
	// MongoDB-specific
	SampledTypes map[string]string // field → inferred type from sampling
	// Redis-specific
	KeyType string // string, list, set, hash, zset, stream
}

// FieldDefinition describes a single column/field/property.
type FieldDefinition struct {
	Name         string
	Kind         FieldKind
	DataType     string // "INTEGER", "TEXT", "VARCHAR(255)", "object", "array", etc.
	Nullable     bool
	PrimaryKey   bool
	AutoIncr     bool
	DefaultValue string
	// MongoDB-specific
	Required bool // inferred from schema validation or sampling
	// Redis-specific
	TTL time.Duration
}

// ForeignKeyDefinition describes a foreign key relationship.
type ForeignKeyDefinition struct {
	FromColumn string
	ToTable    string
	ToColumn   string
	OnDelete   string // CASCADE, SET NULL, RESTRICT, etc.
	OnUpdate   string
}

// IndexDefinition describes a database index.
type IndexDefinition struct {
	Name    string
	Columns []string
	Unique  bool
	Primary bool
}

// ColInfoToFields converts ColInfo slice to FieldDefinition slice.
func ColInfoToFields(cols []ColInfo) []FieldDefinition {
	fields := make([]FieldDefinition, len(cols))
	for i, c := range cols {
		fields[i] = FieldDefinition{
			Name:       c.Name,
			Kind:       FieldKindColumn,
			DataType:   c.Type,
			Nullable:   !c.NotNull,
			PrimaryKey: c.PK,
			AutoIncr:   IsAutoIncrement(c),
			DefaultValue: func() string {
				if c.Dflt.Valid {
					return c.Dflt.String
				}
				return ""
			}(),
		}
	}
	return fields
}

// FKInfoToFKDefs converts FKInfo slice to ForeignKeyDefinition slice.
func FKInfoToFKDefs(fks []FKInfo) []ForeignKeyDefinition {
	defs := make([]ForeignKeyDefinition, len(fks))
	for i, fk := range fks {
		defs[i] = ForeignKeyDefinition{
			FromColumn: fk.From,
			ToTable:    fk.Table,
			ToColumn:   fk.To,
		}
	}
	return defs
}
