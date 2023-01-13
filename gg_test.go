package gg

import (
	"os"
	"strconv"
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
