package env

import (
	"fmt"
	"os"
	"strings"
)

// Get a value from environment and run it through the parse function. Return
// the result if there was one.
//
// If the parsing fails or if the variable was not set an error will be
// returned.
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

	parsedValue, err := parse(rawValue)
	if err != nil {
		var noResult V
		return noResult, fmt.Errorf("Parsing %s value: %w", environmentVariableName, err)
	}

	return parsedValue, nil
}

// Get a value from environment and run it through the parse function. Return
// the result if there was one.
//
// If the parsing fails or if the variable was not set then the fallback value
// will be returned.
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
// the result if there was one.
//
// If the parsing fails or if the variable was not set then this function will
// panic.
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
		for index, part := range separatedString {
			parsedValue, err := parse(part)
			if err != nil {
				return nil, fmt.Errorf("Element %d: %w", index+1, err)
			}

			result = append(result, parsedValue)
		}

		return result, nil
	}
}

// Helper function for parsing floats and similar from environment variables.
//
// # Example Usage
//
//	number, err := env.Get("FLOAT", env.WithBitSize(strconv.ParseFloat, 64))
func WithBitSize[V any](parse func(string, int) (V, error), bitSize int) func(string) (V, error) {
	return func(raw string) (V, error) {
		return parse(raw, bitSize)
	}
}

// Helper function for parsing ints of different bases from environment
// variables.
//
// Pro tip: Passing base 0 with [strconv.ParseInt] and [strconv.ParseUint]
// will make them try to figure out the base by themselves.
//
// # Example Usage
//
//	number, err := env.Get("HEX", env.WithBaseAndBitSize(strconv.ParseUint, 0, 64))
//
// [strconv.ParseInt]: https://pkg.go.dev/strconv#ParseInt
// [strconv.ParseUint]: https://pkg.go.dev/strconv#ParseUint
func WithBaseAndBitSize[V any](parse func(string, int, int) (V, error), base, bitSize int) func(string) (V, error) {
	return func(raw string) (V, error) {
		return parse(raw, base, bitSize)
	}
}

// Helper function for reading strings from the environment.
//
// # Example Usage
//
//	userName, err := env.Get("USERNAME", env.Plain)
func Plain(input string) (string, error) {
	return input, nil
}
