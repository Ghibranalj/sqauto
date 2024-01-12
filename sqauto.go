package sqauto

import (
	"bytes"
	"fmt"
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

func SelectJoin(b sq.StatementBuilderType, table string, obj any, jointable ...string) (string, []any, error) {
	st := reflect.TypeOf(obj)

	selectbuf := bytes.NewBuffer([]byte{})
	otherTable := make(map[string]bool)
	for _, v := range jointable {
		otherTable[v] = true
	}

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		name := field.Tag.Get("sq")
		if name == "" {
			name = toSnakeCase(field.Name)
		}

		stchild := reflect.TypeOf(reflect.ValueOf(obj).Field(i).Interface())
		if otherTable[name] {
			// if pointer, get the type of the pointer
			if stchild.Kind() == reflect.Ptr {
				stchild = stchild.Elem()
			}
			for j := 0; j < stchild.NumField(); j++ {
				fieldchild := stchild.Field(j)
				namechild := fieldchild.Tag.Get("sq")
				if namechild == "" {
					namechild = toSnakeCase(fieldchild.Name)
				}
				// get grandchild

				fmt.Fprintf(selectbuf, `%s.%s AS "%s.%s",`+"\n", name, namechild, name, namechild)
			}
			continue
		}

		fmt.Fprintf(selectbuf, `%s.%s, `+"\n", table, name)
	}
	str := selectbuf.String()
	// remove last comma
	str = str[:len(str)-3]
	sb := b.Select(str).From(table)

	for _, jtable := range jointable {
		// JOIN example ON table.example_id = example.id
		joinstr := fmt.Sprintf("%s ON %s.%s_id = %s.id", jtable, table, jtable, jtable)
		sb = sb.Join(joinstr)
	}

	return sb.ToSql()
}
