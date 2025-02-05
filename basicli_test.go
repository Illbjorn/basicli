package basicli

import (
	"flag"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	noFlags          = []string{"./some/path/bin"}
	flagWithoutValue = []string{"./some/path/bin", "-some"}
	flagWithValue    = []string{"./some/path/bin", "-some", "value"}
)

type BasicFlags struct {
	Key string `basicli:"some,s"`
}

func TestUnmarshal_BasicFlags(t *testing.T) {
	var x = unmarshalPass[BasicFlags](t, noFlags)
	assert.Empty(t, x.Key)
	x = unmarshalFail[BasicFlags](t, flagWithoutValue)
	assert.Empty(t, x.Key)
	x = unmarshalPass[BasicFlags](t, flagWithValue)
	assert.Equal(t, "value", x.Key)
}

type RequiredFlags struct {
	Key string `basicli:"some,s,required=true"`
}

func TestUnmarshal_RequiredFlags(t *testing.T) {
	var x = unmarshalFail[RequiredFlags](t, noFlags)
	assert.Empty(t, x.Key)
	x = unmarshalFail[RequiredFlags](t, flagWithoutValue)
	assert.Empty(t, x.Key)
	x = unmarshalPass[RequiredFlags](t, flagWithValue)
	assert.Equal(t, "value", x.Key)
}

type FlagsWithDefaults struct {
	Key string `basicli:"some,s,default=Hello, world!"`
}

func TestUnmarshal_FlagsWithDefaults(t *testing.T) {
	var x = unmarshalPass[FlagsWithDefaults](t, noFlags)
	assert.Equal(t, "Hello, world!", x.Key)
	x = unmarshalFail[FlagsWithDefaults](t, flagWithoutValue)
	assert.Equal(t, "Hello, world!", x.Key)
	x = unmarshalPass[FlagsWithDefaults](t, flagWithValue)
	assert.Equal(t, "value", x.Key)
}

type FlagsWithIncompatibleTags struct {
	Key string `basicli:"some,s,required=true,default=Hello, world!"`
}

func TestUnmarshal_FlagsWithIncompatibleTags(t *testing.T) {
	unmarshalFail[FlagsWithIncompatibleTags](t, noFlags)
	unmarshalFail[FlagsWithIncompatibleTags](t, flagWithoutValue)
	unmarshalFail[FlagsWithIncompatibleTags](t, flagWithValue)
}

var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func rndStr() string {
	var ret string
	for range 6 {
		var charPosition = rand.IntN(len(chars) - 1)
		var nextChar = chars[charPosition]
		ret += string(nextChar)
	}
	return ret
}

func unmarshalPass[T any](t *testing.T, args []string) T {
	t.Helper()
	var rnd = rndStr()
	set = flag.NewFlagSet(rnd, flag.ContinueOnError)
	set.Usage = func() {}
	os.Args = args
	var v T
	var err = Unmarshal(&v)
	assert.NoError(t, err)
	return v
}

func unmarshalFail[T any](t *testing.T, args []string) T {
	t.Helper()
	var rnd = rndStr()
	set = flag.NewFlagSet(rnd, flag.ContinueOnError)
	set.Usage = func() {}
	os.Args = args
	var v T
	var err = Unmarshal(&v)
	assert.Error(t, err)
	return v
}
