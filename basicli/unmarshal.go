package basicli

import (
  "fmt"
  "os"
  "reflect"
  "slices"
  "strconv"
  "strings"

  "github.com/illbjorn/basicli/basicli/argv"
  "github.com/illbjorn/basicli/basicli/tag"
  "github.com/illbjorn/echo"
)

func Unmarshal[P *T, T any](v P) error {
  if v == nil {
    return fmt.Errorf("received uninitialized ['%T']", v)
  }
  // Parse args and flags
  args, flags := argv.Parse(os.Args[1:])
  return unmarshal(Concrete(reflect.ValueOf(v)), args, flags)
}

func unmarshal(rv reflect.Value, args []string, flags map[string][]string) error {
  var arg string
  if len(args) > 0 {
    arg = args[0]
    args = args[1:]
  }

  for i := range rv.NumField() {
    field := Concrete(rv.Field(i))
    fieldType := rv.Type().Field(i)

    if !field.CanAddr() {
      return fmt.Errorf("found unaddressable field ['%s']", field.Type().Name())
    }

    // Parse the struct tag (if any)
    tagStr, tagFound := fieldType.Tag.Lookup(structTag)
    tag := tag.Parse(tagStr)

    ////////////////////////////////////////////////////////////////////////////
    // Flags

    fieldName := fieldType.Name
    for k, v := range flags {
      // First check for a folded match on the map key against the current flag
      if strings.EqualFold(k, fieldName) {
        fieldSet(field, v)
        continue
      }
      // If we found a struct tag, check there
      if tagFound {
        // Check the primary tag name first
        if strings.EqualFold(k, tag.Name) {
          fieldSet(field, v)
          continue
        }
        // Check aliases second
        if slices.ContainsFunc(tag.Aliases, containsStrFold(k)) {
          fieldSet(field, v)
        }
      }
    }

    ////////////////////////////////////////////////////////////////////////////
    // Subcommands

    // No non-structs beyond this point
    if len(arg) == 0 || field.Kind() != reflect.Struct {
      continue
    }

    // Evaluate a subcommand (arg) match on this field's name or tag values
    //
    // If we find one, recurse and populate that member as well
    if strings.EqualFold(fieldName, arg) {
      err := unmarshal(field, args, flags)
      if err != nil {
        return err
      }
    }
  }

  return nil
}

func fieldSet(rv reflect.Value, vs []string) error {
  if len(vs) == 0 {
    return nil
  }

  rv = Concrete(rv)
  rt := rv.Type()

  if !rv.CanAddr() {
    return fmt.Errorf("found unaddressable field ['%s']", rt.Name())
  }

  switch rv.Kind() {
  case reflect.Bool:
    if strings.EqualFold(vs[0], "true") {
      rv.SetBool(true)
    }

  case reflect.String:
    rv.SetString(vs[0])

  case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
    i, err := strconv.ParseInt(vs[0], 10, 64)
    if err != nil {
      return err
    }
    rv.SetInt(i)

  case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
    ui, err := strconv.ParseUint(vs[0], 10, 64)
    if err != nil {
      return err
    }
    rv.SetUint(ui)

  default:
    echo.Fatalf("Found unexpected field kind ['%s'].", rv.Kind())
    panic("")
  }

  return nil
}
