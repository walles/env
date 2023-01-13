package gg

import "os"

// Get a value from environment and run it through the parse function. Return
// the result if there was one, or an error if the parsing fails.
func Get[R any](environmentVariableName string, parse func(string) (R, error)) (R, error) {
	rawValue := os.Getenv(environmentVariableName)
	return parse(rawValue)
}

// Get a value from environment and run it through the parse function. Return
// the result if there was one, fallback if the parsing fails.
func GetOr[R any](environmentVariableName string, parse func(string) (R, error), fallback R) R {
	rawValue := os.Getenv(environmentVariableName)
	parsedValue, err := parse(rawValue)
	if err != nil {
		return fallback
	}

	return parsedValue
}
