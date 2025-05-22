package basicli

import "reflect"

// Concrete retrieves the value underlying any pointers or interfaces
// recursively.
//
// For example, if Concrete is called on an interface which itself contains a
// pointer, the returned reflect.Value will be the dereferenced value the
// interface contained.
func Concrete(v reflect.Value) reflect.Value {
  if v.Kind() == reflect.Pointer {
    if v.IsNil() {
      return reflect.Value{}
    }
    return Concrete(v.Elem())
  }

  if v.Kind() == reflect.Interface {
    if v.IsNil() {
      return reflect.Value{}
    }
    return Concrete(v.Elem())
  }

  return v
}
