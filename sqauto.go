package sqauto

import (
	"reflect"

	sq "github.com/Masterminds/squirrel"
)

type Table struct {
	Name   string
	Object any

	// primary key column name
	// default "id"
	PrimaryKey string
}

// if a is not zero, return a
// if a is zero, return b
func or[T any](a, b T) T {
	if reflect.ValueOf(a).IsZero() {
		return b
	}
	return a
}

func Insert(b sq.StatementBuilderType, table Table) sq.InsertBuilder {
	if reflect.TypeOf(table.Object).Kind() == reflect.Slice ||
		reflect.TypeOf(table.Object).Kind() == reflect.Array {
		return insertMany(b, table)
	}
	return insertOne(b, table)
}

func insertOne(b sq.StatementBuilderType, table Table) sq.InsertBuilder {
	table.PrimaryKey = or(table.PrimaryKey, "id")

	obj := table.Object
	st := reflect.TypeOf(table.Object)
	ib := b.Insert(table.Name)
	var values []any

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		name := or(field.Tag.Get("sq"), toSnakeCase(field.Name))
		if reflect.ValueOf(obj).Field(i).IsZero() {
			continue
		}
		if !insertIncluded(field.Type) {
			continue
		}
		ib = ib.Columns(name)
		values = append(values, reflect.ValueOf(obj).Field(i).Interface())
	}
	ib = ib.Values(values...)

	return ib
}

// UpdateAll will update all fields into the table
// if primary key is zero or empty, it will not update
func UpdateAll(b sq.StatementBuilderType, table Table) sq.UpdateBuilder {
	table.PrimaryKey = or(table.PrimaryKey, "id")

	obj := table.Object
	st := reflect.TypeOf(obj)
	ub := b.Update(table.Name)

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		name := or(field.Tag.Get("sq"), toSnakeCase(field.Name))
		if name == table.PrimaryKey &&
			reflect.ValueOf(table.Object).Field(i).IsZero() {
			continue
		}
		ub = ub.Set(name, reflect.ValueOf(obj).Field(i).Interface())
	}

	return ub
}

// Update will only update fields that are not zero or empty
// if primary key is zero or empty, it will not update
func Update(b sq.StatementBuilderType, table Table) sq.UpdateBuilder {
	table.PrimaryKey = or(table.PrimaryKey, "id")

	obj := table.Object
	st := reflect.TypeOf(obj)
	ub := b.Update(table.Name)

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		name := or(field.Tag.Get("sq"), toSnakeCase(field.Name))
		// if value is not zero or empty, set it
		if reflect.ValueOf(obj).Field(i).IsZero() {
			continue
		}

		ub = ub.Set(name, reflect.ValueOf(obj).Field(i).Interface())
	}

	return ub
}

// insert slice
func insertMany(b sq.StatementBuilderType, table Table) sq.InsertBuilder {
	sliceType := reflect.TypeOf(table.Object)
	if sliceType.Kind() != reflect.Slice || sliceType.Kind() != reflect.Array {
		panic("object must be a slice")
	}
	ib := b.Insert(table.Name)
	skipped := map[int]bool{}

	elem := sliceType.Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)

		if !insertIncluded(field.Type) {
			skipped[i] = true
			continue
		}

		name := or(field.Tag.Get("sq"), toSnakeCase(field.Name))
		ib = ib.Columns(name)
	}

	// iterate over slice
	slice := reflect.ValueOf(table.Object)
	for i := 0; i < slice.Len(); i++ {
		// iterate over struct fields
		for j := 0; j < elem.NumField(); j++ {
			if skipped[j] {
				continue
			}
			val := slice.Index(i).Field(j).Interface()
			if reflect.ValueOf(val).IsZero() {
				val = nil
			}
			ib = ib.Values(val)
		}
	}

	return ib
}
