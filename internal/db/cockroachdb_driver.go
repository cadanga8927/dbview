package db

import (
	"context"
	"strings"
)

// CockroachDBDriver implements Driver for CockroachDB databases.
// It embeds PostgreSQLDriver because CockroachDB is wire-compatible with
// the PostgreSQL protocol and the lib/pq package works with CockroachDB
// servers. Only methods that need CockroachDB-specific behavior are
// overridden; everything else delegates to the embedded PostgreSQLDriver.
type CockroachDBDriver struct {
	PostgreSQLDriver
}

func (d *CockroachDBDriver) Open(ctx context.Context, dsn string) error {
	if after, ok := strings.CutPrefix(dsn, "cockroachdb://"); ok {
		dsn = "postgres://" + after
	} else {
		dsn = "postgres://" + strings.TrimPrefix(dsn, "cockroach://")
	}
	dsn = replaceLocalhost(dsn)
	return d.PostgreSQLDriver.Open(ctx, dsn)
}

func (d *CockroachDBDriver) Kind() DataSourceKind { return KindCockroachDB }
