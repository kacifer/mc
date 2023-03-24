package main

import (
	"github.com/gin-gonic/gin"
	"github.com/notpm/mc"
	"github.com/notpm/mc/mgin"
	"github.com/pkg/errors"
	"log"
)

func main() {
	var authMiddleware func(mgin.HandlerFunc) mgin.HandlerFunc
	authMiddleware = func(f mgin.HandlerFunc) mgin.HandlerFunc {
		return func(c *mgin.Context) {
			c.Set(mgin.IDKey, mc.StringToUint(c.GetHeader("X-User-ID")))
			f(c)
		}
	}

	var testHandler mgin.HandlerFunc
	testHandler = func(c *mgin.Context) {
		log.Printf("id param: %v", c.IDParam())
		log.Printf("id query: %v", c.IDQuery())
		log.Printf("user id context: %v", c.MustIDContext())
		c.JSON(200, "OK")
	}

	r := mgin.Default()

	r.GET("/test", authMiddleware(testHandler))
	r.GET("/test/:id", authMiddleware(testHandler))
	r.GET("/error", func(c *mgin.Context) {
		c.AbortWithStatusJSON(400, errors.New("error happened"))
	})

	r.Use(mgin.WrapHandler(gin.Logger()))

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
