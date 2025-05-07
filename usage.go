package basicli

import (
	"reflect"

	"github.com/illbjorn/basicli/internal/builder"
	"github.com/illbjorn/basicli/tag"
)

func Usage[P *T, T any](v P) string {
	if v == nil {
		return ""
	}

	var desc string
	if v, ok := any(v).(interface{ About() string }); ok {
		desc = v.About() + "\n"
	}

	var rv = Concrete(reflect.ValueOf(v))
	var rt = rv.Type()
	var l = buildLayout(rt)

	var b = builder.New(1024)
	b.Write(desc)
	b.Writeln("FLAGS")
	b.Newline()

	b.In()
	for _, row := range l.Rows {
		b.Indent()
		for i, col := range row {
			var longest = l.Widths[i] + 2
			b.Write(col)
			b.Writen(" ", longest-len(col))
		}
		b.Newline()
	}
	b.Out()

	return b.String()
}

type Layout struct {
	Rows   [][]string
	Widths []int
}

func buildLayout(rt reflect.Type) Layout {
	var l Layout

	// Size on the high side (flag, alias, description)
	l.Rows = make([][]string, 0, rt.NumField())

	for i := range rt.NumField() {
		var ft = rt.Field(i)

		// Ignore unexported fields

		if !ft.IsExported() {
			continue
		}

		// Analyze struct tags

		var tagStr string
		var ok bool
		if tagStr, ok = ft.Tag.Lookup(structTag); !ok {
			continue
		}
		var tag = tag.Parse(tagStr)

		// Assemble the row and column widths

		// name |  aliases* | default | description
		var sz = 1 + len(tag.Aliases) + 1 + 1
		if sz > len(l.Widths) {
			l.Widths = append(l.Widths, make([]int, sz-len(l.Widths))...)
		}

		var row = make([]string, sz)

		// Name (primary flag name)
		var name = "--" + tag.Name
		if len(name) > l.Widths[0] {
			l.Widths[0] = len(name)
		}
		row[0] = name

		// Aliases
		for i, alias := range tag.Aliases {
			var n = i + 1
			alias = "-" + alias
			row[n] = alias
			if len(alias) > l.Widths[n] {
				l.Widths[n] = len(alias)
			}
		}

		// Default
		var n = 1 + len(tag.Aliases)
		row[n] = tag.Default
		if len(tag.Default) > l.Widths[n] {
			l.Widths[n] = len(tag.Default)
		}
		n++

		// Description
		//
		// TODO: Placeholder until I figure out an ergonomic way to get BYO
		// descriptions here.
		row[n] = ""
		if len(row[n]) > l.Widths[n] {
			l.Widths[n] = len(row[n])
		}

		l.Rows = append(l.Rows, row)
	}

	return l
}
