package utils

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func Delete(db *sql.DB, table string, id uuid.UUID) error {
	query := `
	DELETE FROM $1
	WHERE id = $2;
	`
	_, err := db.Exec(query, table, id)
	if err != nil {
		return err
	}
	return nil
}
