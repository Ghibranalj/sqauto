package sqauto_test

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/ghibranalj/sqauto"
)



func TestJoin(t *testing.T){
	expected := trim(`
		SELECT *,
		spouse.id AS "spouse.id",
		spouse.maiden_name AS "spouse.maiden_name",
		spouse.birth_date AS "spouse.birth_date",
		spouse.car_id AS "spouse.car_id"
		FROM client
		JOIN spouse ON client.spouse_id = spouse.id`)

	query, _, _ := sqauto.Join(sq.StatementBuilder,
		sqauto.Table{Name: "client", Object: Client{}},
		sqauto.JoinTable{Name: "spouse", Object: Spouse{}}).ToSql()

	if query != expected {
		t.Errorf("\nExpected: %s\nGot: %s", expected, query)
		return
	}

	logSuccess(t, query, expected)
}
