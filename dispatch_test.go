package basicli

import (
	"fmt"
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

var check = assert.Check

type MockDispatch struct {
	CMD MockDispatchInner `basicli:"cmd"`
}

func (MockDispatch) GoodButAlsoBad() error { return fmt.Errorf("oh no!") }
func (MockDispatch) Bad()                  {}
func (MockDispatch) Exec() error           { return nil }

type MockDispatchInner struct{}

func (MockDispatchInner) Hello() error { return nil }

func TestDispatch(t *testing.T) {
	var md MockDispatch

	// Call the default `Exec` method
	os.Args = []string{""}
	err := Dispatch(&md)
	check(t, err == nil)

	// Call the MockDispatch.Inner.Hello() method
	os.Args = []string{"", "cmd", "hello"}
	assert.NilError(t, Dispatch(&md))
	check(t, err == nil)

	// Call the MockDispatch.Bad() method, expect an error since the function
	// signature is illegal
	os.Args = []string{"", "bad"}
	err = Dispatch(&md)
	check(t, err != nil)

	// Call the MockDispatch.GoodButAlsoBad() method, expect an error returned by
	// the function body naturally
	os.Args = []string{"", "goodbutalsobad"}
	err = Dispatch(&md)
	check(t, err != nil)
}
