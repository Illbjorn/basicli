package basicli

import (
  "reflect"
  "testing"

  "gotest.tools/v3/assert"
)

func TestConcrete(t *testing.T) {
  // Create some contrived example
  var (
    a any = struct {
      x string
    }{"Hello!"}

    b = struct {
      y *any
    }{&a}
  )

  // Carry out the manual "resolution" path, as a sanity check alongside the
  // Concrete function
  rv := reflect.ValueOf(b.y)
  assert.Check(t, rv.Kind() == reflect.Pointer, "(%s)", rv.Kind())
  rv = rv.Elem()
  assert.Check(t, rv.Kind() == reflect.Interface, "(%s)", rv.Kind())
  rv = rv.Elem()
  assert.Check(t, rv.Kind() == reflect.Struct, "(%s)", rv.Kind())

  // Immediately produces the value we're after
  rv = Concrete(reflect.ValueOf(b.y))
  assert.Check(t, rv.Kind() == reflect.Struct)

  // For bonus points, look up the field name and confirm its type
  rv = rv.FieldByName("x")
  rv = Concrete(rv)
  assert.Check(t, rv.Kind() == reflect.String)
}
