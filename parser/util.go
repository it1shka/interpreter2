package parser

func Includes(element string, elements []string) bool {
	for _, e := range elements {
		if e == element {
			return true
		}
	}
	return false
}
