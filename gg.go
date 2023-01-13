package gg

import (
	"fmt"
	"os"
)

// Get a value from environment and run it through the parse function. Return
// the result if there was one, or an error if the parsing fails or if the
// variable was not set.
func Get[R any](environmentVariableName string, parse func(string) (R, error)) (R, error) {
	rawValue, found := os.LookupEnv(environmentVariableName)
	if !found {
		var noResult R
		return noResult, fmt.Errorf("Environment variable not set: %s", environmentVariableName)
	}

	return parse(rawValue)
}

// Get a value from environment and run it through the parse function. Return
// the result if there was one, fallback if the parsing fails or if the variable
// was not set.
func GetOr[R any](environmentVariableName string, parse func(string) (R, error), fallback R) R {
	rawValue, found := os.LookupEnv(environmentVariableName)
	if !found {
		return fallback
	}

	parsedValue, err := parse(rawValue)
	if err != nil {
		return fallback
	}

	return parsedValue
}
