package basicli

import (
  "fmt"
  "reflect"
  "strconv"
  "strings"
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
    i, err := strconv.Atoi(v)
    if err != nil {
      return err
    }
    rv.SetInt(int64(i))

  case reflect.Bool:
    if strings.ToLower(v) == "true" {
      rv.SetBool(true)
    }

  default:
    return fmt.Errorf("found unexpected conversion type ['%s']", rt.Kind())
  }

  return nil
}
