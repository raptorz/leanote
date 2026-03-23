package postgres

import (
	"database/sql"

	"github.com/leanote/leanote/app/db"
)

type PostgresTransaction struct {
	tx *sql.Tx
}

func NewPostgresTransaction(tx *sql.Tx) db.Transaction {
	return &PostgresTransaction{tx: tx}
}

func (pt *PostgresTransaction) Commit() error {
	return pt.tx.Commit()
}

func (pt *PostgresTransaction) Rollback() error {
	return pt.tx.Rollback()
}

func (pt *PostgresTransaction) Exec(query string, args ...interface{}) (db.Result, error) {
	return pt.tx.Exec(query, args...)
}

func (pt *PostgresTransaction) Query(query string, args ...interface{}) (db.Rows, error) {
	return pt.tx.Query(query, args...)
}

func (pt *PostgresTransaction) QueryRow(query string, args ...interface{}) db.Row {
	return pt.tx.QueryRow(query, args...)
}
