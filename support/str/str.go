package str

// Contains checks if string is present in given Slice
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Addslashes returns a string with backslashes added before characters that need to be escaped.
func Addslashes(str string) string {
	var tmpRune []rune
	strRune := []rune(str)

	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		}
	}

	return string(tmpRune)
}
