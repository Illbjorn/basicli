package tag

type tag struct {
  Name    string
  Aliases []string
  Default string
  Flags   tagFlags
}

type tagFlags uint8

const (
  flagRequired tagFlags = 1 << iota
  flagHasDefault
)

func (self tagFlags) Required() bool {
  return self&flagRequired == flagRequired
}

func (self tagFlags) HasDefault() bool {
  return self&flagHasDefault == flagHasDefault
}
