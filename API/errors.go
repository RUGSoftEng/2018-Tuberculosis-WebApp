package main

import (
	"database/sql"
	"github.com/pkg/errors"
	http "net/http"
)

const (
	// Error Messages :
	// ErrLang : Error when language specified is invalid
	ErrLang = "Invalid Language, languages must be one of ['EN', 'NL', 'DE', 'RO']"
	// ErrDecodeJSON : Error message when decoding failed
	ErrDecodeJSON                 = "Failed to decode JSON body"
	ErrDBTransactionStartFaillure = "Failed to start database transaction"
	ErrDBInsert                   = "Failed to execute the insert statement"
	ErrDBUpdate                   = "Failed to execute the update statement"
	ErrDBDelete                   = "Failed to execute the delete statement"
	ErrDBCommit                   = "Failed to commit transaction to the database"
	ErrDBSelect                   = "Failed to execute select statement"
	ErrDBScan                     = "Failed to scan queried rows"
	ErrDBAfter                    = "Error after iterating over rows"
)

func errorWithRollback(err error, tx *sql.Tx) error {
	err2 := tx.Rollback()
	if err2 != nil {
		err = errors.New(err.Error() + "\n Rollback failed:" + err2.Error())
	}
	return err
}
