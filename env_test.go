package env

import (
	"os"
	"strconv"
	"testing"
	"time"

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
	is.Equal(err.Error(), `Parsing TEST value: strconv.Atoi: parsing "kalaspuffar": invalid syntax`)

	is.NoErr(os.Unsetenv("TEST"))
	_, err = Get("TEST", strconv.Atoi)
	is.True(err != nil)
	is.Equal(err.Error(), "Environment variable not set: TEST")
}

func TestGetListOfInts(t *testing.T) {
	is := is.New(t)

	is.NoErr(os.Setenv("TEST", "1,2,3,4"))
	value, err := Get("TEST", ListOf(strconv.Atoi, ","))
	is.NoErr(err)
	is.Equal([]int{1, 2, 3, 4}, value)

	is.NoErr(os.Setenv("TEST", "kalaspuffar"))
	_, err = Get("TEST", ListOf(strconv.Atoi, ","))
	is.True(err != nil)
	is.Equal(err.Error(), `Parsing TEST value: Element 1: strconv.Atoi: parsing "kalaspuffar": invalid syntax`)

	is.NoErr(os.Unsetenv("TEST"))
	_, err = Get("TEST", ListOf(strconv.Atoi, ","))
	is.True(err != nil)
	is.Equal(err.Error(), "Environment variable not set: TEST")
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

func TestGetOrString(t *testing.T) {
	is := is.New(t)

	is.NoErr(os.Setenv("TEST", "Hello"))
	is.Equal("Hello", GetOr("TEST", String, "??"))

	is.NoErr(os.Unsetenv("TEST"))
	is.Equal("??", GetOr("TEST", String, "??"))
}

func TestGetOrFloat64(t *testing.T) {
	is := is.New(t)

	is.NoErr(os.Setenv("TEST", "1.23"))
	is.Equal(1.23, GetOr("TEST", WithBitSize(strconv.ParseFloat, 64), 13))

	is.NoErr(os.Setenv("TEST", "kalaspuffar"))
	is.Equal(13.0, GetOr("TEST", WithBitSize(strconv.ParseFloat, 64), 13))

	is.NoErr(os.Unsetenv("TEST"))
	is.Equal(13.0, GetOr("TEST", WithBitSize(strconv.ParseFloat, 64), 13))
}

func TestGetOrHexInt(t *testing.T) {
	is := is.New(t)

	is.NoErr(os.Setenv("TEST", "C0de"))
	is.Equal(int64(0xc0de), GetOr("TEST", WithBaseAndBitSize(strconv.ParseInt, 16, 64), 13))

	is.NoErr(os.Setenv("TEST", "0xC0de"))
	is.Equal(int64(0xc0de), GetOr("TEST", WithBaseAndBitSize(strconv.ParseInt, 0, 64), 13))

	is.NoErr(os.Setenv("TEST", "kalaspuffar"))
	is.Equal(int64(13), GetOr("TEST", WithBaseAndBitSize(strconv.ParseInt, 16, 64), 13))

	is.NoErr(os.Unsetenv("TEST"))
	is.Equal(int64(13), GetOr("TEST", WithBaseAndBitSize(strconv.ParseInt, 16, 64), 13))
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
		is.Equal(err.Error(), `Parsing TEST value: strconv.Atoi: parsing "kalaspuffar": invalid syntax`)
	}()
	MustGet("TEST", strconv.Atoi)

	is.Fail() // MustGet with unparsable variable value should have panicked
}

func TestGetTimestamp(t *testing.T) {
	is := is.New(t)

	timeString := "2022-03-04T05:06:07+08:00"
	is.NoErr(os.Setenv("TEST", timeString))

	parsed, err := Get("TEST", WithTimeSpec(time.Parse, time.RFC3339))
	is.NoErr(err)

	location := time.FixedZone("UTC+8", 8*60*60)
	is.True(parsed.Equal(time.Date(2022, 3, 4, 5, 6, 7, 0, location)))
}

func TestGetMap(t *testing.T) {
	is := is.New(t)

	mapString := "a:11,b:22"
	is.NoErr(os.Setenv("TEST", mapString))

	parsed, err := Get("TEST", Map(String, ":", strconv.Atoi, ","))
	is.NoErr(err)

	is.Equal(parsed, map[string]int{
		"a": 11,
		"b": 22,
	})

	// Test a broken entry
	is.NoErr(os.Setenv("TEST", "1:apa:true"))
	_, err = Get("TEST", Map(strconv.Atoi, ":", String, ","))
	is.Equal(err.Error(), `Parsing TEST value: Element 1 doesn't have exactly one separator (":"): 1:apa:true`)

	// Test a malformed key
	is.NoErr(os.Setenv("TEST", "apa:fisk"))
	_, err = Get("TEST", Map(strconv.Atoi, ":", String, ","))
	is.Equal(err.Error(), `Parsing TEST value: Key 1: strconv.Atoi: parsing "apa": invalid syntax`)

	// Test a malformed value
	is.NoErr(os.Setenv("TEST", "apa:fisk"))
	_, err = Get("TEST", Map(String, ":", strconv.Atoi, ","))
	is.Equal(err.Error(), `Parsing TEST value: Value 1: strconv.Atoi: parsing "fisk": invalid syntax`)
}
