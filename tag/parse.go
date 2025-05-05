package tag

func Parse(v string) tag {
  var t tag
  if len(v) == 0 {
    return t
  }

  scanner := tagScanner{v: v, i: -1}

  for {
    next := scanner.peek(1)
    if next == '\x00' {
      scanner.adv()
      scanner.mark(markerID)
      break
    }

    switch {
    case next == '=':
      // '='
      scanner.adv()

      // Determine the type of marker we need to set
      buffered := scanner.buffered()
      var markerKind int
      if buffered == "default" {
        markerKind = markerDefault

      } else if buffered == "required" {
        markerKind = markerRequired

      } else {
        panic(buffered)
      }

      // Manually move the chains
      scanner.bump()

      // Consume to ',' or EOF
      for {
        next = scanner.peek(1)
        if next == ',' || next == '\x00' {
          scanner.adv() // ','
          scanner.mark(markerKind)
          break
        }
        scanner.adv()
      }
      continue

    case next == ',':
      // Mark the name/alias
      scanner.adv()
      scanner.mark(markerID)

    default:
      scanner.adv()
    }
  }

  // Imprint the tag and return
  scanner.imprint(&t)

  return t
}
