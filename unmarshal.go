package basicli

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// We use a custom FlagSet to allow for easier unit testing. A flag name can
// only be registered a single time to a flag set. This means a series of unit
// tests will clobber each other with multiple definitions of the same flag. In
// this scenario `flag.Parse()` returns an error.
//
// NOTE: Important to remember when parsing args using a custom flag set - the
// parse will "fail" (PRODUCE NO ERROR, BUT ALSO PARSE NO FLAGS) if the provided
// args contain a file path or subcommand.
var set = flag.NewFlagSet("basicli", flag.ContinueOnError)

// # Overview
//
// Unmarshals os.Args to provided struct `v`.
//
// # Example
//
//	type Some {
//	  // Defines flag "--field" with alias "--f"
//	  // The struct field type indicates the type of flag value we expect (string)
//	  // The word `required` appearing in the tag means `Unmarshal()` will
//	  // produce an error if this flag isn't provided
//	  Field string `basicli:"field,f,required"`
//	}
//
//	func main() {
//	  var x Some
//	  os.Args = []string{"/some/bin/path", "-f", "some_value"}
//	  basicli.Unmarshal(&x)
//	  fmt.Println(x.Field) // "some_value"
//	}
func Unmarshal[P *T, T any](v P) error {
	if v == nil {
		return fmt.Errorf("received nil ['%T']", v)
	}

	// Safe reflect.Value dereference here as `Unmarshal` is constrained by P *T
	var rvalue = reflect.ValueOf(v).Elem()
	var rtype = rvalue.Type()

	// We accumulate a list of closures which validate any flags marked "required"
	// were provided
	//
	// We invoke these closures further down
	var requiredChecks []func() error

	for i := range rvalue.NumField() {
		var sfv = rvalue.Field(i) // Struct field value
		var sft = rtype.Field(i)  // Struct field type

		if !sfv.CanAddr() {
			return fmt.Errorf("found unaddressable field ['%s']", sfv.Type().Name())
		}

		var tag string
		var ok bool
		if tag, ok = sft.Tag.Lookup("basicli"); !ok {
			continue
		}

		var check func() error
		var err error
		switch sft.Type.Kind() {
		case reflect.String:
			check, err = registerFlag(tag, sfv, sft, set.StringVar)

		case reflect.Int:
			check, err = registerFlag(tag, sfv, sft, set.IntVar)

		case reflect.Bool:
			check, err = registerFlag(tag, sfv, sft, set.BoolVar)

		default:
			return fmt.Errorf("found unexpected struct field type ['%s']", sft.Type.Kind())
		}

		if err != nil {
			return fmt.Errorf("failed to register flag: %s", err)
		} else if check != nil {
			requiredChecks = append(requiredChecks, check)
		}
	}

	var err error
	if err = set.Parse(os.Args[1:]); err != nil {
		return err
	}

	// Process any required flag checks
	for _, check := range requiredChecks {
		var err error
		if err = check(); err != nil {
			return err
		}
	}

	return nil
}

type FlagInitializer[P *T, T basic] func(P, string, T, string)

func registerFlag[P *T, T basic](
	tag string,
	v reflect.Value,
	t reflect.StructField,
	initializer FlagInitializer[P, T],
) (func() error, error) {
	var ptr = (P)(v.Addr().UnsafePointer())

	// Init the flag
	var f Flag[P, T]
	var err error
	if f, err = NewFlag[P, T](tag, ptr); err != nil {
		return nil, fmt.Errorf(
			"failed to initialize flag ['%s']: %s",
			t.Name,
			err)
	}

	// Register the flag
	for _, name := range f.Names {
		initializer(ptr, name, f.Default, "")
	}

	if f.Required {
		return requiredFlagValidator(f.Names), nil
	}

	return nil, nil
}

func requiredFlagValidator(names []string) func() error {
	return func() error {
		for _, name := range names {
			if isFlagSet(name) {
				return nil
			}
		}

		return fmt.Errorf("required flag ['%s'] was not provided", strings.Join(names, ", "))
	}
}

func isFlagSet(name string) bool {
	var isSet bool

	set.Visit(func(f *flag.Flag) {
		if f.Name == name {
			isSet = true
			return
		}
	})

	return isSet
}
