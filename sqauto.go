package sqauto

import (
	"reflect"

	sq "github.com/Masterminds/squirrel"
)

func Insert(b sq.StatementBuilderType, table string, obj any) (string, []any, error) {
	st := reflect.TypeOf(obj)
	ib := b.Insert(table)
	var values []any

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		name := field.Tag.Get("sq")
		ib = ib.Columns(name)
		values = append(values, reflect.ValueOf(obj).Field(i).Interface())
	}
	ib = ib.Values(values...)

	return ib.ToSql()
}

func Update(b sq.StatementBuilderType, table string, obj any) (string, []any, error) {
	st := reflect.TypeOf(obj)
	ub := b.Update(table)
	var values []any

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		name := field.Tag.Get("sq")
		ub = ub.Set(name, reflect.ValueOf(obj).Field(i).Interface())
		values = append(values, reflect.ValueOf(obj).Field(i).Interface())
	}

	return ub.ToSql()
}

func CoalesceUpdate(b sq.StatementBuilderType, table string, obj any) (string, []any, error) {
	st := reflect.TypeOf(obj)
	ub := b.Update(table)
	var values []any

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		name := field.Tag.Get("sq")
		// if value is not zero or empty, set it
		if reflect.ValueOf(obj).Field(i).IsZero() {
			continue
		}
		ub = ub.Set(name, reflect.ValueOf(obj).Field(i).Interface())
		values = append(values, reflect.ValueOf(obj).Field(i).Interface())
	}

	return ub.ToSql()
}

func Select(b sq.StatementBuilderType, table string, obj any) (string, []any, error) {
	st := reflect.TypeOf(obj)
	sb := b.Select()

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		name := field.Tag.Get("sq")
		sb = sb.Columns(name)
	}

	return sb.From(table).ToSql()
}

func SelectJoin(b sq.StatementBuilderType, table string, obj any, jointable ...string) (string, []any, error) {

	st := reflect.TypeOf(obj)
	selectstr := ""
	otherTable := make(map[string]bool)
	for _, v := range jointable {
		otherTable[v] = true
	}

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		name := field.Tag.Get("sq")

		if otherTable[name] {

			stchild := reflect.TypeOf(reflect.ValueOf(obj).Field(i).Interface())
			// if pointer, get the type of the pointer
			if stchild.Kind() == reflect.Ptr {
				stchild = stchild.Elem()
			}
			for j := 0; j < stchild.NumField(); j++ {
				fieldchild := stchild.Field(j)
				namechild := fieldchild.Tag.Get("sq")
				selectstr += name + "." + namechild + " AS " + "\"" + name + "." + namechild + "\"" + ", "
			}
			continue
		}

		selectstr += table + "." + name + ", "
	}
	// remove last comma
	selectstr = selectstr[:len(selectstr)-2]
	sb := b.Select(selectstr).From(table)

	for _, jtable := range jointable {
		sb = sb.Join(jtable + " ON " + table + "." + jtable + "_id = " + jtable + ".id")
	}

	return sb.ToSql()
}
