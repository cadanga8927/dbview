package db

import (
	"testing"
)

func TestDetectDriver(t *testing.T) {
	tests := []struct {
		dsn       string
		wantKind  DataSourceKind
		wantIsNew bool
	}{
		{"mariadb://user:pass@localhost:3306/mydb", KindMariaDB, true},
		{"cockroach://user:pass@localhost:26257/mydb", KindCockroachDB, true},
		{"cockroachdb://user:pass@localhost:26257/mydb", KindCockroachDB, true},
		{"sqlserver://user:pass@localhost:1433/mydb", KindMSSQL, true},
		{"mssql://user:pass@localhost:1433/mydb", KindMSSQL, true},
		{"mysql://user:pass@localhost:3306/mydb", KindMySQL, false},
		{"postgres://user:pass@localhost:5432/mydb", KindPostgreSQL, false},
		{"postgresql://user:pass@localhost:5432/mydb", KindPostgreSQL, false},
		{"mongodb://localhost:27017/mydb", KindMongoDB, false},
		{"mongodb+srv://user:pass@cluster/mydb", KindMongoDB, false},
		{"redis://localhost:6379", KindRedis, false},
		{"rediss://localhost:6379", KindRedis, false},
		{"cassandra://localhost:9042/testdb", KindCassandra, false},
		{"cassandra://user:pass@host:9042/keyspace", KindCassandra, false},
		{"./mydb.db", KindSQLite, false},
		{"/path/to/file.db", KindSQLite, false},
	}

	for _, tt := range tests {
		t.Run(tt.dsn, func(t *testing.T) {
			got, _ := DetectDriver(tt.dsn)
			if got != tt.wantKind {
				t.Errorf("DetectDriver(%q) = %q, want %q", tt.dsn, got, tt.wantKind)
			}
		})
	}
}

func TestDetectDriverPriority(t *testing.T) {
	mariadbKind, _ := DetectDriver("mariadb://user:pass@host/db")
	if mariadbKind != KindMariaDB {
		t.Errorf("mariadb:// should detect as MariaDB, got %q", mariadbKind)
	}
	mysqlKind, _ := DetectDriver("mysql://user:pass@host/db")
	if mysqlKind != KindMySQL {
		t.Errorf("mysql:// should detect as MySQL, got %q", mysqlKind)
	}
	cockroachKind, _ := DetectDriver("cockroachdb://user:pass@host/db")
	if cockroachKind != KindCockroachDB {
		t.Errorf("cockroachdb:// should detect as CockroachDB, got %q", cockroachKind)
	}
	pgKind, _ := DetectDriver("postgres://user:pass@host/db")
	if pgKind != KindPostgreSQL {
		t.Errorf("postgres:// should detect as PostgreSQL, got %q", pgKind)
	}
	mssqlKind, _ := DetectDriver("sqlserver://user:pass@host/db")
	if mssqlKind != KindMSSQL {
		t.Errorf("sqlserver:// should detect as MSSQL, got %q", mssqlKind)
	}
}

func TestOpenDriver_Unsupported(t *testing.T) {
	_, err := OpenDriver(t.Context(), "mysql://invalid:invalid@127.0.0.1:99999/nonexistent")
	if err == nil {
		t.Error("expected error opening invalid MySQL DSN")
	}
}

func TestReplaceLocalhost(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"postgres://user:pass@localhost:5432/db", "postgres://user:pass@127.0.0.1:5432/db"},
		{"mysql://user:pass@localhost:3306/db", "mysql://user:pass@127.0.0.1:3306/db"},
		{"sqlserver://user:pass@localhost:1433/db", "sqlserver://user:pass@127.0.0.1:1433/db"},
		{"postgres://user:pass@db.example.com:5432/db", "postgres://user:pass@db.example.com:5432/db"},
		{"redis://127.0.0.1:6379", "redis://127.0.0.1:6379"},
	}
	for _, tt := range tests {
		got := replaceLocalhost(tt.input)
		if got != tt.want {
			t.Errorf("replaceLocalhost(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestParseDSN(t *testing.T) {
	tests := []struct {
		dsn      string
		scheme   string
		host     string
		database string
	}{
		{"mysql://user:pass@localhost:3306/mydb", "mysql", "localhost:3306", "mydb"},
		{"postgres://user:pass@localhost:5432/testdb", "postgres", "localhost:5432", "testdb"},
		{"redis://localhost:6379", "redis", "localhost:6379", ""},
		{"mariadb://user:pass@localhost:3306/mydb", "mariadb", "localhost:3306", "mydb"},
		{"cockroachdb://user:pass@localhost:26257/mydb", "cockroachdb", "localhost:26257", "mydb"},
		{"sqlserver://user:pass@localhost:1433/mydb", "sqlserver", "localhost:1433", "mydb"},
	}
	for _, tt := range tests {
		t.Run(tt.dsn, func(t *testing.T) {
			scheme, host, database, err := ParseDSN(tt.dsn)
			if err != nil {
				t.Fatalf("ParseDSN(%q) error: %v", tt.dsn, err)
			}
			if scheme != tt.scheme {
				t.Errorf("scheme = %q, want %q", scheme, tt.scheme)
			}
			if host != tt.host {
				t.Errorf("host = %q, want %q", host, tt.host)
			}
			if database != tt.database {
				t.Errorf("database = %q, want %q", database, tt.database)
			}
		})
	}
}
