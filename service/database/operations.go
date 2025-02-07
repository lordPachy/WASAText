package database

import (
	"database/sql"
	"fmt"
)

// It inserts data into the database.
//
// Values should be in the exact number of columns of the table, in the format ('str1', 'str2', n, k , Null...)
func (db *appdbimpl) Insert(table string, values string) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s VALUES %s", table, values)
	res, err := db.c.Exec(query)
	return res, err
}

// It updates data into a table.
//
// Update should be in the format col1 = 'val1', col2 = n,... colk = 'valk'.
// Condition should be in the same format.
func (db *appdbimpl) Update(table string, update string, conditions string) (sql.Result, error) {
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, update, conditions)
	res, err := db.c.Exec(query)
	return res, err
}

// It retrieves data from a table.
//
// Columns should be in the format col1, col2, col3... coln.
// Condition should be in the format col1 = val1, col2 = val2,... colk = valk.
func (db *appdbimpl) Select(columns string, table string, conditions string) (*sql.Rows, error) {
	// Actual query
	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", columns, table, conditions)
	res, err := db.c.Query(query)
	if err != nil {
		return nil, err
	}

	return res, err
}

// It retrieves data from a table.
//
// Columns should be in the format col1, col2, col3... coln.
// Filter should be in the format col1 = val1, col2 = val2,... colk = valk.
func (db *appdbimpl) Filter(columns string, table string, group_by string, conditions string) (*sql.Rows, error) {
	// Actual query
	query := fmt.Sprintf("SELECT %s FROM %s GROUP BY %s HAVING %s", columns, table, group_by, conditions)
	res, err := db.c.Query(query)
	if err != nil {
		return nil, err
	}

	return res, err
}

// It deletes data from a table.
//
// Conditions should be in the format col1 = val1, col2 = val2,... colk = valk.
func (db *appdbimpl) Delete(table string, conditions string) (*sql.Rows, error) {
	// Actual query
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, conditions)
	res, err := db.c.Query(query)
	if err != nil {
		return nil, err
	}

	return res, err
}
