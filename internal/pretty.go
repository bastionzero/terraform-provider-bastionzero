package internal

import "fmt"

// PrettyOneOf returns "(one of options...)"
func PrettyOneOf[T ~string](options []T) string {
	// (one of options...)
	result := "(one of "
	for i, option := range options {
		if i != len(options)-1 {
			result += fmt.Sprintf("`%s`, ", option)
		} else {
			result += fmt.Sprintf("or `%s`)", option)
		}
	}
	return result
}

// PrettyRFC3339Timestamp returns "formatted as a UTC timestamp string in [RFC
// 3339](https://datatracker.ietf.org/doc/html/rfc3339) format"
func PrettyRFC3339Timestamp() string {
	return "formatted as a UTC timestamp string in [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) format"
}
