package mjwt

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMjwt(t *testing.T) {
	assertions := require.New(t)

	jwt := NewImpl([]byte("secret"), 30*24*time.Hour)
	jwt.NowFunc = func() time.Time {
		return time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	signed, err := jwt.SignedStringForID(1)
	assertions.Nil(err)
	assertions.Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDE4LTAxLTMxVDAwOjAwOjAwWiIsImlhdCI6IjIwMTgtMDEtMDFUMDA6MDA6MDBaIiwiaWQiOjF9.azFMfNqKYtbqO8tbcRSNRtm6HlhvXXzfyV68ozpSXZA", signed)

	token, err := jwt.Parse(signed)
	assertions.Nil(err)

	id, err := jwt.ExtractID(token)
	assertions.Nil(err)
	assertions.Equal(uint(1), id)

	id, err = jwt.ExtractIDFromSignedString(signed)
	assertions.Nil(err)
	assertions.Equal(uint(1), id)

	id, err = jwt.ExtractIDFromHeader("Bearer " + signed)
	assertions.Nil(err)
	assertions.Equal(uint(1), id)

	signed, err = jwt.SignedStringForName("test")
	assertions.Nil(err)
	assertions.Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDE4LTAxLTMxVDAwOjAwOjAwWiIsImlhdCI6IjIwMTgtMDEtMDFUMDA6MDA6MDBaIiwibmFtZSI6InRlc3QifQ.NpN2HsbulQtMrTLBH3w3ZManSE5oM5FIeMm0lQASJfA", signed)

	token, err = jwt.Parse(signed)
	assertions.Nil(err)

	name, err := jwt.ExtractName(token)
	assertions.Nil(err)
	assertions.Equal("test", name)

	name, err = jwt.ExtractNameFromSignedString(signed)
	assertions.Nil(err)
	assertions.Equal("test", name)

	name, err = jwt.ExtractNameFromHeader("Bearer " + signed)
	assertions.Nil(err)
	assertions.Equal("test", name)
}
