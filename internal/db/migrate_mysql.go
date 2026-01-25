package db

import (
	"database/sql"
	"fmt"
)

func MigrateMySQL(db *sql.DB) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS tasks (
  id         VARCHAR(64)  PRIMARY KEY,
  title      VARCHAR(200) NOT NULL,
  done       TINYINT(1)   NOT NULL DEFAULT 0,
  created_at DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`
	if _, err := db.Exec(ddl); err != nil {
		return fmt.Errorf(`migrate mysql ddl: %w`, err)
	}
	return nil
}
