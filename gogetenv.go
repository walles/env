package gg

import "os"

func GetOr[R any](environmentVariableName string, parse func(string) (R, error), fallback R) R {
	rawValue := os.Getenv(environmentVariableName)
	parsedValue, err := parse(rawValue)
	if err != nil {
		return fallback
	}

	return parsedValue
}
