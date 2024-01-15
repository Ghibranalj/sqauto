package sqauto_test

import (
	"database/sql"
	"time"
)

// To emulate guregu/null
type nullString struct {
	sql.NullString
}

type Client struct {
	ID int64
	Name string
	Email string
	SpouseID int64
	Spouse Spouse
}

type Spouse struct {
	ID int64
	MaidenName nullString
	BirthDate time.Time
}

type Car struct {
	ID int64
	Make string
	LicencePlate nullString
}
