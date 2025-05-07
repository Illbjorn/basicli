package builder

import (
	"fmt"
	"io"
	"strings"
	"unsafe"
)

func New(sz int) *Builder {
	return &Builder{
		b:      make([]byte, 0, sz),
		indent: "  ",
	}
}

type Builder struct {
	b      []byte
	depth  int
	indent string
}

/*------------------------------------------------------------------------------
 * Indentation
 *----------------------------------------------------------------------------*/

// Set the string value used to represent an indent.
func (self *Builder) SetIndentString(s string) {
	self.indent = s
}

// Move the indent level in one.
func (self *Builder) In() *Builder {
	self.depth++
	return self
}

// Move the indent level out one.
func (self *Builder) Out() *Builder {
	if self.depth <= 0 {
		return self
	}
	self.depth--
	return self
}

// Apply an indent.
func (self *Builder) Indent() *Builder {
	self.b = append(self.b, strings.Repeat(self.indent, self.depth)...)
	return self
}

/*------------------------------------------------------------------------------
 * Misc QoL
 *----------------------------------------------------------------------------*/

// Apply a newline ('\n').
func (self *Builder) Newline() *Builder {
	self.b = append(self.b, '\n')
	return self
}

/*------------------------------------------------------------------------------
 * Write Fns
 *----------------------------------------------------------------------------*/

// Write a string to the underlying buffer.
func (self *Builder) Write(s string) *Builder {
	self.b = append(self.b, s...)
	return self
}

// Write a given string `n` times to the underlying buffer.
func (self *Builder) Writen(s string, n int) *Builder {
	if n <= 0 {
		return self
	}
	self.b = append(self.b, strings.Repeat(s, n)...)
	return self
}

// Write a string terminated by a newline to the underlying buffer.
func (self *Builder) Writeln(s string) *Builder {
	self.b = append(self.b, s...)
	self.b = append(self.b, '\n')
	return self
}

// Write a formatted string to the underlying buffer.
func (self *Builder) Writef(s string, vs ...any) *Builder {
	self.b = fmt.Appendf(self.b, s, vs...)
	return self
}

// Write a formatted string terminated by a newline to the underlying buffer.
func (self *Builder) Writefln(s string, vs ...any) *Builder {
	self.b = fmt.Appendf(self.b, s, vs...)
	self.b = append(self.b, '\n')
	return self
}

// Write an optionally formatted string to the underlying buffer if `cond` is
// `true`.
func (self *Builder) Writeif(cond bool, v string, vs ...any) *Builder {
	if cond {
		self.b = fmt.Appendf(self.b, v, vs...)
	}
	return self
}

/*------------------------------------------------------------------------------
 * Control Fns
 *----------------------------------------------------------------------------*/

// Reset the underlying buffer length, retaining its capacity.
func (self *Builder) Reset() {
	self.b = self.b[:0]
}

/*------------------------------------------------------------------------------
 * fmt.Stringer
 *----------------------------------------------------------------------------*/

func (self *Builder) String() string {
	if len(self.b) == 0 {
		return ""
	}
	return string(self.b)
}

func (self *Builder) UString() string {
	if len(self.b) == 0 {
		return ""
	}
	sd := unsafe.SliceData(self.b)
	return unsafe.String(sd, len(self.b))
}

/*------------------------------------------------------------------------------
 * io.Reader
 *----------------------------------------------------------------------------*/

var _ io.Reader = (*Builder)(nil)

func (self *Builder) Read(dst []byte) (n int, err error) {
	if self.depth >= len(self.b) {
		return 0, io.EOF
	}

	n = copy(dst, self.b[self.depth:])
	self.depth += n
	if n < len(dst) {
		err = io.EOF
	}

	return n, nil
}

/*------------------------------------------------------------------------------
 * io.Seeker
 *----------------------------------------------------------------------------*/

var _ io.Seeker = (*Builder)(nil)

func (self *Builder) Seek(offset int64, whence int) (n int64, err error) {
	switch whence {
	case io.SeekStart:
		if offset < 0 {
			err = io.EOF
		}
		self.depth = int(offset)
		n = offset

	case io.SeekCurrent:
		self.depth += int(offset)
		n = int64(self.depth)

	case io.SeekEnd:
		if offset >= int64(len(self.b)) {
			err = io.EOF
		}
		self.depth = len(self.b) - int(offset)
		n = int64(self.depth)

	default:
		panic("What?")
	}

	return
}
