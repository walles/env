Functions for parsing environment variable values into typed values.

# Examples

Note that the resulting values are all typed.

```go
// Enabled will be of type bool
enabled := env.GetOr("ENABLED", strconv.ParseBool, false)

// Duration will be of type time.Duration
duration, err := env.Get("PORT", time.ParseDuration)

// Username will be of type string. If it's not set in the environment,
// then MustGet will panic.
username := env.MustGet("USERNAME", env.Plain)

// LuckyNumbers will be of type []int
luckyNumbers, err := env.Get("LUCKY_NUMBERS", env.ListOf(strconv.Atoi, ","))

// FluffyNumber will be a 64 bit precision float
fluffyNumber, err := env.Get("FLOAT", env.WithBitSize(strconv.ParseFloat, 64))
```

# Alternatives

If you like bindings based APIs better then this one seems popular:

* <https://github.com/kelseyhightower/envconfig>
