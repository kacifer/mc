package mgin

import (
	"github.com/kacifer/mc"
	"github.com/kacifer/mc/mjwt"
	"net/http"
)

func CreateAuthMiddleware(jwt mjwt.Engine, skipAuthPaths []string) HandlerFunc {
	return func(c *Context) {
		if !mc.SliceContains(skipAuthPaths, c.Request.URL.Path) {
			_, claims, err := jwt.ValidateHeader(c.Request.Header.Get("Authorization"))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, &E{
					Code:    http.StatusUnauthorized,
					Message: err.Error(),
				})
				return
			}

			c.Set(IDKey, claims[mjwt.IDKey])
		}

		c.Next()
	}
}
