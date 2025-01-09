package main

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cristalhq/builq"
)

func TestSQLInjectionWithSqlmock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	qf := func(db *sql.DB, tableName, username string, unsafeQuery bool) (*sql.Rows, error) {
		if unsafeQuery {
			query := fmt.Sprintf("SELECT * FROM %s WHERE username = '%s'", tableName, username)
			return db.Query(query)
		}

		return db.Query("SELECT * FROM users WHERE username = ?", username)
	}

	tests := []struct {
		name    string
		user    string
		table   string
		unsafe  bool
		prepare func(mock sqlmock.Sqlmock, user string)
	}{
		{
			name: "Tautologies",
			user: "anything' OR 'x'='x",
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username = ?")).
					WithArgs(user).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
		},
		{
			name: "Illegal/Logically Incorrect Queries",
			user: "admin' AND 1=2 UNION SELECT * FROM users --",
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username = ?")).
					WithArgs(user).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
		},
		{
			name: "Union Query",
			user: "admin' UNION SELECT * FROM users --",
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username = ?")).
					WithArgs(user).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
		},
		{
			name: "Piggy-Backed Queries",
			user: "admin'; DROP TABLE users; --",
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username = ?")).
					WithArgs(user).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:   "FmtSprintf Injection",
			user:   "'; DROP TABLE users; --",
			unsafe: true,
			table:  "users",
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery("SELECT \\* FROM users WHERE username = '.*; DROP TABLE users; --'").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
		},
		{
			name:   "TableName Injection",
			user:   "admin",
			table:  "users; DROP TABLE sensitive_data; --",
			unsafe: true,
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery(".+ DROP TABLE sensitive_data;.+").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(mock, tc.user)

			_, err := qf(db, tc.table, tc.user, tc.unsafe)

			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				t.Errorf("Unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSQLInjectionPreventionUsingBuilq(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name    string
		user    string
		bb      func(user string) *builq.Builder
		prepare func(mock sqlmock.Sqlmock, user string)
	}{
		{
			name: "Tautologies",
			user: "anything' OR 'x'='x",
			bb: func(user string) *builq.Builder {
				return builq.New()("SELECT %s FROM %s WHERE username = %$", builq.Columns{"username"}, "users", user)
			},
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery("SELECT username FROM users WHERE username = \\$1").
					WithArgs(user).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
		},
		{
			name: "Illegal/Logically Incorrect Queries",
			user: "admin' AND 1=2 UNION SELECT * FROM users --",
			bb: func(user string) *builq.Builder {
				return builq.New()("SELECT %s FROM %s WHERE username = %$", builq.Columns{"username"}, "users", user)
			},
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT username FROM users WHERE username = $1")).
					WithArgs(user).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
		},
		{
			name: "Union Query",
			user: "admin' UNION SELECT * FROM users --",
			bb: func(user string) *builq.Builder {
				return builq.New()("SELECT %s FROM %s WHERE username = %$", builq.Columns{"username"}, "users", user)
			},
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT username FROM users WHERE username = $1")).
					WithArgs(user).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
		},
		{
			name: "Piggy-Backed Queries",
			user: "admin'; DROP TABLE users; --",
			bb: func(user string) *builq.Builder {
				return builq.New()("SELECT %s FROM %s WHERE username = %$", builq.Columns{"username"}, "users", user)
			},
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT username FROM users WHERE username = $1")).
					WithArgs(user).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "FmtSprintf Injection",
			user: "'; DROP TABLE users; --",
			bb: func(user string) *builq.Builder {
				return builq.New()("SELECT * FROM users WHERE username = '%s'", user)
			},
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery("SELECT \\* FROM users WHERE username = '.*; DROP TABLE users; --'").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
		},
		{
			name: "TableName Injection",
			user: "admin",
			bb: func(user string) *builq.Builder {
				tableName := "users; DROP TABLE sensitive_data; --"
				return builq.New()("SELECT * FROM %s WHERE username = '%s'", tableName, user)
			},
			prepare: func(mock sqlmock.Sqlmock, user string) {
				mock.ExpectQuery(".+ DROP TABLE sensitive_data;.+").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.prepare(mock, tc.user)

			query, args, err := tc.bb(tc.user).Build()
			if err != nil {
				t.Fatalf("could not build query: %v", err)
			}

			_, err = db.Query(query, args...)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				t.Errorf("Unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}
		})
	}
}
