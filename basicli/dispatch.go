package basicli

import (
  "fmt"
  "os"
  "reflect"
  "slices"
  "strings"

  "github.com/illbjorn/basicli/basicli/argv"
  "github.com/illbjorn/basicli/basicli/tag"
)

func Dispatch[P *T, T any](v P) error {
  args, _ := argv.Parse(os.Args[1:])

  rv := Concrete(reflect.ValueOf(v))
  if rv.Kind() != reflect.Struct {
    return fmt.Errorf("received non-struct type ['%T'] in call to Dispatch", v)
  }

  return dispatch[T](rv, args)
}

func dispatch[T any](rv reflect.Value, args []string) error {
  switch {
  case len(args) == 0:
    // When we run out of args, dispatch the `Exec` method on the current
    // reflected value, which should be a method
    method := rv.MethodByName(methodExec)
    if method.Kind() != reflect.Func {
      return fmt.Errorf(
        "failed to locate ['%s'] method on type ['%s']",
        methodExec, rv.Type().Name(),
      )
    }

    // Call it!
    outputs := method.Call(nil)
    if len(outputs) == 0 {
      return nil
    }
    return outputs[0].Interface().(error)

  case len(args) == 1:
    // On the last (or first) positional arg, what we're after can either be:
    //
    // 1. An immediate method on our current reflected type
    // 2. Another nested struct, under which we call the `Exec` method
    //
    // First, look for an immediate matching method
    sought := args[0]
    args = args[1:]
    for i := range rv.NumMethod() {
      method := rv.Method(i)
      methodType := rv.Type().Method(i)
      name := methodType.Name
      if strings.EqualFold(sought, name) {
        // Call it!
        res := method.Call(nil)
        if len(res) != 1 {
          return fmt.Errorf(
            "expected a single error return value, found ['%[1]T']: %[1]v",
            res,
          )
        }
        err := res[0]
        if err.CanInterface() {
          if err.IsNil() {
            return nil
          }
          return res[0].Interface().(error)
        }
      }
    }

    // Alternatively, look through the methods and attempt to recurse again
    // and catch an `Exec` method
    for i := range rv.NumField() {
      field := rv.Field(i)
      fieldType := rv.Type().Field(i)

      if strings.EqualFold(sought, fieldType.Name) {
        return dispatch[T](field, args)
      }
    }

    return fmt.Errorf("failed to identify leaf subcommand method target ['%s'] on type ['%s']", sought, rv.Type().Name())

  default:
    // With >1 positional args, we continue to descend into the nested structs
    // (subcommands)
    sought := args[0]
    args = args[1:]
    for i := range rv.NumField() {
      field := rv.Field(i)
      fieldType := rv.Type().Field(i)

      // Look for a struct tag first
      tagStr, ok := fieldType.Tag.Lookup(structTag)
      if ok {
        // Parse the tag string
        tag := tag.Parse(tagStr)

        // Compare what we're after against the primary tag name
        if strings.EqualFold(sought, tag.Name) {
          return dispatch[T](field, args)
        }

        // Iterate any aliases and check those as well
        if slices.ContainsFunc(tag.Aliases, containsStrFold(sought)) {
          return dispatch[T](field, args)
        }
      }

      // Check for a direct match on field name
      if strings.EqualFold(sought, fieldType.Name) {
        return dispatch[T](field, args)
      }
    }
    return fmt.Errorf("failed to identify nested composite type for subcommand ['%s']", sought)
  }
}

func containsStrFold(v string) func(sliceValue string) bool {
  return func(sliceValue string) bool {
    return strings.EqualFold(v, sliceValue)
  }
}
