package basicli

import (
	"testing"

	"github.com/illbjorn/echo"
)

type Test struct {
	Silent string `basicli:"silent,s,default=not-silent"`
	Debug  string `basicli:"debug,d,default=ok"`
}

func (t Test) About() string {
	return `
Hello, world!
`
}

func TestUsage(t *testing.T) {
	var x Test

	var text = Usage(&x)

	echo.Info(text)
}
