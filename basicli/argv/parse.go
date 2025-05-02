package argv

// parse_args.go handles parsing of command-line inputs to positional args or
// flags

import "strings"

func Parse(inputs []string) ([]string, map[string][]string) {
  if len(inputs) == 0 {
    return nil, nil
  }
  return parse(inputs, make([]string, 0, 4), make(map[string][]string))
}

func parse(inputs []string, args []string, flags map[string][]string) ([]string, map[string][]string) {
  if len(inputs) == 0 {
    return args, flags
  }

  more := len(inputs) > 1
  cur, isFlag := stripPrefix(inputs[0], "-")
  var nextIsFlag bool
  if more && strings.HasPrefix(inputs[1], "-") {
    nextIsFlag = true
  }

  switch {
  case isFlag && more && !nextIsFlag: // Flag with value
    flags[cur] = append(flags[cur], inputs[1])
    inputs = inputs[2:]

  case isFlag: // Boolean flag
    flags[cur] = append(flags[cur], "true")
    inputs = inputs[1:]

  default: // Positional arg
    args = append(args, cur)
    inputs = inputs[1:]
  }

  return parse(inputs, args, flags)
}

func stripPrefix(v string, cutset string) (string, bool) {
  if len(v) == 0 || len(cutset) == 0 {
    return "", false
  }

  var found bool
  for len(v) > 0 {
    if strings.IndexByte(cutset, v[0]) != 0 {
      break
    }
    v = v[1:]
    found = true
  }
  return v, found
}
