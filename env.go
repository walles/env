package gg

import (
	"fmt"
	"os"
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
