package basicli

import (
  "fmt"
  "os"
  "testing"

  "gotest.tools/v3/assert"
)

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
  assert.NilError(t, Dispatch(&md))

  // Call the MockDispatch.Inner.Hello() method
  os.Args = []string{"", "cmd", "hello"}
  assert.NilError(t, Dispatch(&md))

  // Call the MockDispatch.Bad() method, expect an error since the function
  // signature is illegal
  os.Args = []string{"", "bad"}
  assert.Error(t, Dispatch(&md), "expected a single error return value, found ['[]reflect.Value']: []")

  // Call the MockDispatch.GoodButAlsoBad() method, expect an error returned by
  // the function body naturally
  os.Args = []string{"", "goodbutalsobad"}
  assert.Error(t, Dispatch(&md), "oh no!")
}
