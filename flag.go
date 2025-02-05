package basicli

import (
	"fmt"
	"strings"

	"github.com/illbjorn/conv"
)

type Flag[P *T, T basic] struct {
	Names      []string
	Required   bool
	Ptr        P
	Default    T
	HasDefault bool
}

func NewFlag[P *T, T basic](v string, ptr P) (Flag[P, T], error) {
	var f Flag[P, T]
	f.Ptr = ptr

	if len(v) == 0 {
		return f, nil
	}

	var instructions = strings.Split(v, ",")
	for i, instruction := range instructions {
		switch {
		// Flag is Required
		case strings.HasPrefix(instruction, "required="):
			if f.HasDefault {
				return f, ErrRequiredAndDefault
			}
			instruction = strings.TrimPrefix(instruction, "required=")
			f.Required = conv.Cbool(instruction)

		// Default Value
		case strings.HasPrefix(instruction, "default="):
			if f.Required {
				return f, ErrRequiredAndDefault
			}
			// Indicate the flag has a provided default value
			f.HasDefault = true
			// Strip the `default=` prefix from the raw tag string value
			instructions[i] = strings.TrimPrefix(instructions[i], "default=")
			// Since the default tag value could contain a comma, `default=` must come
			// last in the comma-delimited tag values
			//
			// Thus, when we find `default=`, we strip the prefix and rejoin all
			// instructions on a `comma` and return here
			var defaultStr = strings.Join(instructions[i:], ",")
			var err error
			if err = convert[P](defaultStr, &f.Default); err != nil {
				return f, fmt.Errorf(
					"failed to convert default flag value ['%s'] to type ['%T']: %s",
					defaultStr,
					f.Ptr,
					err,
				)
			}
			return f, nil

		// Flag name or alias
		default:
			f.Names = append(f.Names, instruction)
		}
	}

	return f, nil
}
