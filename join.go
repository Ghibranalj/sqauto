package sqauto

import (
	"fmt"
	"reflect"

	sq "github.com/Masterminds/squirrel"
)


type JoinTable struct {
	Name   string
	Object any
	// optional alias for table
	Alias  string

	// name of primary key column default "id"
	PrimaryKey string
	// name of foreign key column in main table default "table_name_id"
	ForeignKey string
}

func Join(b sq.StatementBuilderType, mainTable Table, joins ...JoinTable) sq.SelectBuilder {
	sb := b.Select("*")
	sb = sb.From(mainTable.Name)

	// now select all columns from all joined tables
	for _, join := range joins {
		alias := or(join.Alias, join.Name)

		st := reflect.TypeOf(join.Object)
		for i := 0; i < st.NumField(); i++ {
			field := st.Field(i)
			name := or(field.Tag.Get("sq"), toSnakeCase(field.Name))
			// SELECT alias.field AS "alias.field"
			sb = sb.Column(fmt.Sprintf(`%s.%s AS "%s.%s"`, alias, name, alias, name))
		}

		primaryKey := or(join.PrimaryKey, "id")
		foreignKey := or(join.ForeignKey, fmt.Sprintf("%s_id", join.Name))

		aliasClause := ""
		if join.Alias != "" {
			aliasClause = fmt.Sprintf(" AS %s", join.Alias)
		}

		str := fmt.Sprintf("%s%s ON %s.%s = %s.%s", join.Name, aliasClause, mainTable.Name, foreignKey, alias, primaryKey)
		sb = sb.Join(str)
	}

	return sb
}
