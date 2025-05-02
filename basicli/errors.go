package basicli

import (
	"errors"
)

var ErrRequiredAndDefault = errors.New(
	"flag is marked required, required flags may not have default values",
)
