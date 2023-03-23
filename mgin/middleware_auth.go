package mgin

import (
	"github.com/notpm/mc/mjwt"
	"github.com/surfinggo/mc"
)

func CreateAuthMiddleware(jwt mjwt.Engine, skipAuthPaths []string) HandlerFunc {
	return func(c *Context) {
		if !mc.SliceContains(skipAuthPaths, c.Request.URL.Path) {
			_, claims, err := jwt.ValidateHeader(c.Request.Header.Get("Authorization"))
			if err != nil {
				c.AbortWithStatusJSON(401, &E{
					Code:    401,
					Message: err.Error(),
				})
				return
			}

			c.Set(IDKey, claims[mjwt.IDKey])
		}

		c.Next()
	}
}
