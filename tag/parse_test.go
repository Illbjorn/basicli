package tag

import (
  "testing"

  "gotest.tools/v3/assert"
)

var check = assert.Check

func TestParse(t *testing.T) {
  tag := Parse("")
  check(t, tag.Name == "")
  check(t, len(tag.Aliases) == 0)
  check(t, tag.Flags == 0)
  check(t, tag.Default == "")

  tag = Parse("s")
  assert.Check(t, tag.Name == "s")
  assert.Check(t, len(tag.Aliases) == 0)
  assert.Check(t, tag.Flags == 0)

  tag = Parse("s,a")
  assert.Check(t, tag.Name == "s")
  assert.Check(t, len(tag.Aliases) == 1)
  assert.Check(t, tag.Aliases[0] == "a")
  assert.Check(t, tag.Flags == 0)

  tag = Parse("required=true")
  assert.Check(t, tag.Flags == flagRequired)

  tag = Parse("default=hello world")
  assert.Check(t, tag.Default == "hello world")

  tag = Parse("silent,default=hello,s,required=true")
  assert.Check(t, tag.Name == "silent")
  assert.Check(t, len(tag.Aliases) == 1)
  assert.Check(t, tag.Aliases[0] == "s")
  assert.Check(t, tag.Flags.HasDefault())
  assert.Check(t, tag.Default == "hello")
  assert.Check(t, tag.Flags.Required())
}

func BenchmarkParseTags(b *testing.B) {
  const tagStr = "silent,s,default=hello,required=true"
  for b.Loop() {
    Parse(tagStr)
  }
}
