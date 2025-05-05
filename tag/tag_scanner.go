package tag

type TagScanner struct {
	v       string
	i, j    int
	markers [][3]int
}

func (self *TagScanner) Peek(i int) byte {
	pos := self.i + i
	if pos < 0 {
		return '\x00'
	}
	if pos >= len(self.v) {
		return '\x00'
	}
	return self.v[pos]
}

func (self *TagScanner) Adv() {
	self.i++
}

const (
	markerID = 1 + iota
	markerRequired
	markerDefault
)

func (self *TagScanner) Mark(kind int) {
	if self.i <= self.j {
		return
	}
	self.markers = append(self.markers, [3]int{kind, self.j, self.i})
	self.Bump()
}

func (self *TagScanner) Bump() {
	self.j = self.i + 1
}

func (self *TagScanner) Buffered() string {
	if self.j > self.i {
		return ""
	}
	if self.j < 0 {
		return ""
	}
	if self.i >= len(self.v) {
		return ""
	}
	return self.v[self.j:self.i]
}

func (self *TagScanner) Imprint(tag *Tag) {
	if tag == nil {
		return
	}

	for _, marker := range self.markers {
		kind, low, high := marker[0], marker[1], marker[2]
		v := self.v[low:high]

		switch kind {
		case markerID:
			if len(tag.Name) == 0 {
				tag.Name = v
			} else {
				tag.Aliases = append(tag.Aliases, v)
			}

		case markerRequired:
			if v == "true" {
				tag.Flags |= flagRequired
			}

		case markerDefault:
			tag.Default = v
			tag.Flags |= flagHasDefault
		}
	}
}
