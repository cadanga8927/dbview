package db

import (
	"context"
	"strings"
)

// MariaDBDriver implements Driver for MariaDB databases.
// It embeds MySQLDriver because MariaDB is wire-compatible with the MySQL
// protocol and the go-sql-driver/mysql package works with MariaDB servers.
// Only methods that need MariaDB-specific behavior are overridden;
// everything else delegates to the embedded MySQLDriver.
type MariaDBDriver struct {
	MySQLDriver
}

func (d *MariaDBDriver) Open(ctx context.Context, dsn string) error {
	dsn = strings.TrimPrefix(dsn, "mariadb://")
	dsn = replaceLocalhost(dsn)
	return d.MySQLDriver.Open(ctx, "mysql://"+dsn)
}

func (d *MariaDBDriver) Kind() DataSourceKind { return KindMariaDB }
