package basicli

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/illbjorn/basicli/argv"
	"github.com/illbjorn/basicli/tag"
)

// Unmarshal `os.Args` input to provided `P` instance `v`.
func Unmarshal[P *T, T any](v P) error {
	// Must be a non-nil pointer
	if v == nil {
		var v T
		return fmt.Errorf("received uninitialized [%T]", v)
	}

	// The pointer must be to a struct
	rv := Concrete(reflect.ValueOf(v))
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, found [%s]", rv.Kind())
	}

	// Parse args and flags
	args, flags := argv.Parse(os.Args[1:])

	// Unmarshal and return
	return unmarshal(rv, args, flags, []string{})
}

// unmarshal recursively consumes input `args`, locating a nested struct on `rv`
// with a case-insensitively-matching name. When found, `unmarshal` recurses to
// that nested struct until `args` is empty.
//
// Once `args` is empty, the leaf struct fields are iterated and any flag values
// contained in `flags` which match either the struct field name or struct field
// tag value(s) exactly have the respective provided flag value assigned. This
// assignment includes conversion of the string input to the data type of the
// field.
func unmarshal(rv reflect.Value, args []string, flags map[string][]string, found []string) error {
	// If we have positional args, locate the nested struct and recurse
	if len(args) > 0 {
		// Slice off the first arg
		arg := args[0]
		args = args[1:]
		// Iterate struct fields
		for i := range rv.NumField() {
			ft := rv.Type().Field(i)
			// Look for a struct tag
			t, ok := ft.Tag.Lookup(structTag)
			if ok {
				// Look for a match against the next arg
				parsed := tag.Parse(t)
				if slices.Contains(append(parsed.Aliases, parsed.Name), arg) {
					// Recurse
					return unmarshal(rv.Field(i), args, flags, found)
				}
			}
			// Look for a case-insensitive match against the field name itself
			if len(ft.Name) != len(arg) {
				continue
			}
			if strings.EqualFold(ft.Name, arg) {
				// Recurse
				return unmarshal(rv.Field(i), args, flags, found)
			}
		}
		// We failed to locate a nested member for the referenced subcommand
		return fmt.Errorf("failed to locate subcommand [%s]", arg)
	}

	// Otherwise, we're just here to populate field values
next:
	for i := range rv.NumField() {
		ft := rv.Type().Field(i)
		found = append(found, ft.Name)
		// Check for a direct match on field name
		//
		// NOTE: Flags are case-sensitive, while subcommands are not
		flag, ok := flags[ft.Name]
		if ok {
			// Set the field value
			fieldSet(rv.Field(i), flag)
		}
		// Also look for a struct tag
		t, ok := ft.Tag.Lookup(structTag)
		if ok {
			tag := tag.Parse(t)
			// Register "found" flags
			found = append(found, tag.Name)
			found = append(found, tag.Aliases...)
			// Look for a provided flag value which matches
			for _, name := range append(tag.Aliases, tag.Name) {
				flag, ok := flags[name]
				if ok {
					// Set the field value
					fieldSet(rv.Field(i), flag)
					continue next
				}
			}
			// If we made it here and the tag is required, we have a problem
			if tag.Flags.Required() {
				fmt.Printf("%#v\n", flags)
				return fmt.Errorf("flag [%s] is required but was not provided", tag.Name)
			}
		}
	}

	// Confirm we didn't encounter any flags which were not defined on the struct
	for k := range flags {
		if !slices.Contains(found, k) {
			return fmt.Errorf("received unexpected flag [%s]", k)
		}
	}

	return nil
}

// fieldSet evaluates the type of the struct field contained in `rv`, converting
// `vs` to values of that type then assigning them to the field.
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
		return fmt.Errorf("found unexpected struct field kind [%s]", rv.Kind())
	}

	return nil
}
