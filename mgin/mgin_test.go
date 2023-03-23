package mgin

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMgin(t *testing.T) {
	assertions := require.New(t)

	var testHandler HandlerFunc
	testHandler = func(c *Context) {
		c.JSON(200, "OK")
	}

	r := New()

	r.GET("/test", testHandler)

	assertions.NotNil(r.Routes())
}
