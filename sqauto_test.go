package sqauto_test

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/ghibranalj/sqauto"
	sq "github.com/Masterminds/squirrel"
)

// To emulate guregu/null
type nullString struct {
	sql.NullString
}

type Client struct {
	ID       int64
	Name     string
	Email    string
	SpouseID int64
	Spouse   Spouse
}

type Spouse struct {
	ID         int64
	MaidenName nullString
	BirthDate  time.Time
	CarID      int64
	Cars	   []Car
}

type Car struct {
	ID           int64
	Make         string
	LicencePlate nullString
	SpouseID     int64
}

// remove newlines and tabs
// remove leading and trailing spaces
func trim(s string) string {
	s = strings.Replace(s, "\n", " ", -1)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.TrimSpace(s)
	return s
}

func logSuccess(t *testing.T, query, expected string) {
	t.Logf(`
Query: %s
Matches expected: %s`, query, expected)
}

func TestUpdate(t *testing.T) {
	val := Client{
		Name:     "John Doe",
		Email:    "testing@test.com",
		SpouseID: 2,
	}

	expected := trim(`
		UPDATE client
		SET name = ?, email = ?, spouse_id = ?`)

	query, _, _ := sqauto.Update(sq.StatementBuilder,
		sqauto.Table{Name: "client", Object: val}).ToSql()

	if query != expected {
		t.Errorf("\nExpected: %s\nGot: %s", expected, query)
		return
	}

	logSuccess(t, query, expected)

	spouse := Spouse{
		BirthDate: time.Now(),
	}
	expected = trim(`
		UPDATE spouse
		SET birth_date = ?`)
	query, _, _ = sqauto.Update(sq.StatementBuilder,
		sqauto.Table{Name: "spouse", Object: spouse}).ToSql()

	if query != expected {
		t.Errorf("\nExpected: %s\nGot: %s", expected, query)
		return
	}

	logSuccess(t, query, expected)
}


func TestInsert(t *testing.T) {

	val := Client{
		ID:       1,
		Name:     "John Doe",
		Email:    "testing@test.com",
		SpouseID: 2,
	}

	expected := trim(`
		INSERT INTO client (id,name,email,spouse_id)
		VALUES (?,?,?,?)`)

	query, _, _ := sqauto.Insert(sq.StatementBuilder,
		sqauto.Table{Name: "client", Object: val}).ToSql()

	if query != expected {
		t.Errorf("\nExpected: %s\nGot: %s", expected, query)
		return
	}

	logSuccess(t, query, expected)
}
