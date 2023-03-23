package mgin

import "fmt"

func CreateHealthzHandler(version string) HandlerFunc {
	return func(c *Context) {
		c.JSON(200, fmt.Sprintf("current running version: %s", version))
	}
}
