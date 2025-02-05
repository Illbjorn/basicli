package basicli

import (
	"fmt"
	"reflect"

	"github.com/illbjorn/conv"
)

type basic interface {
	string | int | bool
}

func convert[P *T, T basic](v string, ptr P) error {
	var rv = reflect.ValueOf(ptr).Elem()
	var rt = rv.Type()

	switch rt.Kind() {
	case reflect.String:
		rv.SetString(v)

	case reflect.Int:
		var i64 int64
		var err error
		if i64, err = conv.Cint(v); err != nil {
			return err
		}
		rv.SetInt(i64)

	case reflect.Bool:
		var b = conv.Cbool(v)
		rv.SetBool(b)

	default:
		return fmt.Errorf("found unexpected conversion type ['%s']", rt.Kind())
	}

	return nil
}
