package env

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestGetInt(t *testing.T) {
	is := is.New(t)

	is.NoErr(os.Setenv("TEST", "1234"))
	value, err := Get("TEST", strconv.Atoi)
	is.NoErr(err)
	is.Equal(1234, value)

	is.NoErr(os.Setenv("TEST", "kalaspuffar"))
	_, err = Get("TEST", strconv.Atoi)
	is.True(err != nil)

	is.NoErr(os.Unsetenv("TEST"))
	_, err = Get("TEST", strconv.Atoi)
	is.True(err != nil)
}

func TestGetOrInt(t *testing.T) {
	is := is.New(t)

	is.NoErr(os.Setenv("TEST", "1234"))
	is.Equal(1234, GetOr("TEST", strconv.Atoi, 13))

	is.NoErr(os.Setenv("TEST", "kalaspuffar"))
	is.Equal(13, GetOr("TEST", strconv.Atoi, 13))

	is.NoErr(os.Unsetenv("TEST"))
	is.Equal(13, GetOr("TEST", strconv.Atoi, 13))
}

func TestMustGetIntHappy(t *testing.T) {
	is := is.New(t)

	is.NoErr(os.Setenv("TEST", "1234"))
	value := MustGet("TEST", strconv.Atoi)
	is.Equal(1234, value)
}

func TestMustGetIntUnset(t *testing.T) {
	is := is.New(t)

	is.NoErr(os.Unsetenv("TEST"))

	defer func() {
		r := recover()
		is.True(r != nil)

		err, ok := r.(error)
		is.True(ok)
		is.Equal(err.Error(), "Environment variable not set: TEST")
	}()
	MustGet("TEST", strconv.Atoi)

	is.Fail() // MustGet on unset variable should have panicked
}

func TestMustGetIntParseError(t *testing.T) {
	is := is.New(t)

	is.NoErr(os.Setenv("TEST", "kalaspuffar"))

	defer func() {
		r := recover()
		is.True(r != nil)

		err, ok := r.(error)
		is.True(ok)
		is.True(strings.Contains(err.Error(), "kalaspuffar"))
	}()
	MustGet("TEST", strconv.Atoi)

	is.Fail() // MustGet with unparsable variable value should have panicked
}
