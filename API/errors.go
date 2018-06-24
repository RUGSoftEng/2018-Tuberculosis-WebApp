package main

import (
	"database/sql"
	"github.com/pkg/errors"
)

// Error Messages
const (
	// *** Client errors ***
	ErrLang            = "Invalid Language, languages must be one of ['EN', 'NL', 'DE', 'RO']"
	ErrDecodeJSON      = "Failed to decode JSON body"
	ErrObjectNotFound  = "A requested object does not exist in the database"
	ErrInvalidVariable = "Invalid variable specified in url"
	ErrEmptyVariable   = "Variable is empty"
	ErrDateFormat      = "Error wrong date format, expected YYYY-MM-DD"
	ErrUsernameTaken   = "Username is already taken"

	// *** Server Errors ***
	ErrDBTransactionStartFaillure = "Failed to start database transaction"
	ErrDBInsert                   = "Failed to execute the insert statement"
	ErrDBUpdate                   = "Failed to execute the update statement"
	ErrDBDelete                   = "Failed to execute the delete statement"
	ErrDBCommit                   = "Failed to commit transaction to the database"
	ErrDBSelect                   = "Failed to execute select statement"
	ErrDBScan                     = "Failed to scan queried rows"
	ErrDBAfter                    = "Error after iterating over rows"
	ErrHash                       = "Failed to execute hash"
)

func errorWithRollback(err error, tx *sql.Tx) error {
	err2 := tx.Rollback()
	if err2 != nil {
		err = errors.New(err.Error() + "\n Rollback failed:" + err2.Error())
	}
	return err
}

func selectErrorHandle(err error, stdStatus int, stdMessage string) (status int, errMessage string) {
	switch err {
	case sql.ErrNoRows:
		status = StatusObjectNotFound
		errMessage = ErrObjectNotFound
	default:
		status = stdStatus
		errMessage = stdMessage
	}
	return
}
