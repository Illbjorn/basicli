package basicli

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

type Sample struct {
	Silent     bool   `basicli:"silent,s,required=true"`
	Debug      bool   `basicli:"debug,d"`
	Path       string `basicli:"path,p"`
	Subcommand Subcommand
}

type Subcommand struct {
	Name string `basicli:"name,n"`
}

func TestUnmarshal(t *testing.T) {
	var sample Sample

	// (good) Basic case
	os.Args = []string{
		"", "--silent", "--debug", "--path", "hellope/world",
	}
	assert.NilError(t, Unmarshal(&sample))
	assert.Check(t, sample.Debug)
	assert.Check(t, sample.Path == "hellope/world")
	assert.Check(t, sample.Silent)

	// (good) Subcommand reference which exists
	os.Args = []string{
		"", "subcommand", "--name", "cmft",
	}
	assert.NilError(t, Unmarshal(&sample))
	assert.Equal(t, sample.Subcommand.Name, "cmft")

	// (bad) Extra, undefined flag
	os.Args = []string{
		"", "--silent", "--debug", "--path", "hello", "-x",
	}
	assert.Error(t, Unmarshal(&sample), "received unexpected flag [x]")

	// (bad) Subcommand reference which does not exist
	os.Args = []string{
		"", "a", "b", "--silent", "--debug", "--path", "hello",
	}
	assert.Error(t, Unmarshal(&sample), "failed to locate subcommand [a]")
}
