package internal

import "fmt"

func PrettyOneOf[T ~string](options []T) string {
	// (one of options...)
	result := "(one of "
	for i, option := range options {
		if i != len(options)-1 {
			result += fmt.Sprintf("\"%s\", ", option)
		} else {
			result += fmt.Sprintf("or \"%s\")", option)
		}
	}
	return result
}
