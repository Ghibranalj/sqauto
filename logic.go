package sqauto

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"time"
)

var (
	scannerType = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	valuerType  = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
)

func implements(t, iface reflect.Type) bool {
	return t.Implements(iface) || reflect.PtrTo(t).Implements(iface)
}

func includedIface(t, iface reflect.Type) bool {
	if t == reflect.TypeOf(time.Time{}) {
		return true
	}
	switch t.Kind() {
	case reflect.Struct:
		return implements(t, iface)
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map:
		return includedIface(t, iface)
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Complex64, reflect.Complex128:
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer, reflect.Invalid:
		return false
	}
	return true
}

func selectIncluded(t reflect.Type) bool {
	return includedIface(t, scannerType)
}

func insertIncluded(t reflect.Type) bool {
	return includedIface(t, valuerType)
}
