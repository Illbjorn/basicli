package tag

func Parse(v string) Tag {
	var t Tag
	if len(v) == 0 {
		return t
	}

	scanner := TagScanner{v: v, i: -1}

	for {
		next := scanner.Peek(1)
		if next == '\x00' {
			scanner.Adv()
			scanner.Mark(markerID)
			break
		}

		switch {
		case next == '=':
			// '='
			scanner.Adv()

			// Determine the type of marker we need to set
			buffered := scanner.Buffered()
			var markerKind int
			if buffered == "default" {
				markerKind = markerDefault

			} else if buffered == "required" {
				markerKind = markerRequired

			} else {
				panic(buffered)
			}

			// Manually move the chains
			scanner.Bump()

			// Consume to ',' or EOF
			for {
				next = scanner.Peek(1)
				if next == ',' || next == '\x00' {
					scanner.Adv() // ','
					scanner.Mark(markerKind)
					break
				}
				scanner.Adv()
			}
			continue

		case next == ',':
			// Mark the name/alias
			scanner.Adv()
			scanner.Mark(markerID)

		default:
			scanner.Adv()
		}
	}

	// Imprint the tag and return
	scanner.Imprint(&t)

	return t
}
