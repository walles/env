package env

import (
	"fmt"
	"os"
	"strings"
)

// Get a value from environment and run it through the parse function. Return
// the result if there was one, or an error if the parsing fails or if the
// variable was not set.
//
// # Example Usage
//
//	port, err := env.Get("PORT", strconv.Atoi)
func Get[V any](environmentVariableName string, parse func(string) (V, error)) (V, error) {
	rawValue, found := os.LookupEnv(environmentVariableName)
	if !found {
		var noResult V
		return noResult, fmt.Errorf("Environment variable not set: %s", environmentVariableName)
	}

	return parse(rawValue)
}

// Get a value from environment and run it through the parse function. Return
// the result if there was one, fallback if the parsing fails or if the variable
// was not set.
//
// # Example Usage
//
//	port := env.GetOr("PORT", strconv.Atoi, 8080)
func GetOr[V any](environmentVariableName string, parse func(string) (V, error), fallback V) V {
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

// Get a value from environment and run it through the parse function. Return
// the result if there was one, or panic if the parsing fails or if the variable
// was not set.
//
// # Example Usage
//
//	port := env.MustGet("PORT", strconv.Atoi)
func MustGet[V any](environmentVariableName string, parse func(string) (V, error)) V {
	parsedValue, err := Get(environmentVariableName, parse)
	if err != nil {
		panic(err)
	}

	return parsedValue
}

// Helper function for reading lists from environment variables.
//
// # Example Usage
//
//	numbers, err := env.Get("NUMBERS", env.ListOf(strconv.Atoi, ","))
func ListOf[V any](parse func(string) (V, error), separator string) func(string) ([]V, error) {
	return func(stringWithSeparators string) ([]V, error) {
		separatedString := strings.Split(stringWithSeparators, separator)

		var result []V
		for _, part := range separatedString {
			parsedValue, err := parse(part)
			if err != nil {
				return nil, err
			}

			result = append(result, parsedValue)
		}

		return result, nil
	}
}
