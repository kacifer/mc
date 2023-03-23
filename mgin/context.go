package mgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/surfinggo/mc"
	"net/http"
)

const IDKey = "id"

const IDQuery = "id"

const IDParam = "id"

type Context struct {
	*gin.Context
}

func (c *Context) MustBindQuery(obj any) (ok bool) {
	if err := c.Context.ShouldBindQuery(obj); err != nil {
		c.AbortAndWriteError(http.StatusBadRequest, errors.Wrap(err, "query decode error"))
		return false
	}
	return true
}

func (c *Context) MustBindJSON(obj any) (ok bool) {
	if err := c.Context.ShouldBindJSON(obj); err != nil {
		c.AbortAndWriteError(http.StatusBadRequest, errors.Wrap(err, "JSON decode error"))
		return false
	}
	return true
}

func (c *Context) MustUintContext(key string) uint {
	return c.MustGet(key).(uint)
}

func (c *Context) MustIDContext() uint {
	return c.MustUintContext(IDKey)
}

func (c *Context) UintContext(key string) (uint, bool) {
	id, exist := c.Get(key)
	if !exist {
		return 0, false
	}
	return id.(uint), true
}

func (c *Context) IDContext() (uint, bool) {
	return c.UintContext(IDKey)
}

func (c *Context) UintQuery(key string) uint {
	return mc.StringToUint(c.Query(key))
}

func (c *Context) IDQuery() uint {
	return c.UintQuery(IDQuery)
}

func (c *Context) UintParam(key string) uint {
	return mc.StringToUint(c.Param(key))
}

func (c *Context) IDParam() uint {
	return c.UintParam(IDParam)
}

func (c *Context) AbortAndWriteError(code int, err any) {
	switch err.(type) {
	case *E:
		c.AbortWithStatusJSON(code, err.(*E))
	case error:
		c.AbortWithStatusJSON(code, &E{
			Code:    code,
			Message: err.(error).Error(),
		})
	case string:
		c.AbortWithStatusJSON(code, &E{
			Code:    code,
			Message: err.(string),
		})
	default:
		c.AbortWithStatusJSON(code, &E{
			Code:    code,
			Message: fmt.Sprintf("%v", err),
		})
	}
}

func (c *Context) AbortWithInternalError(code int, err error) {
	_ = c.Error(err)
	if gin.IsDebugging() {
		c.AbortWithStatusJSON(code, &E{
			Code:    code,
			Message: err.Error(),
		})
		return
	} else {
		// hide error message if not debugging
		c.AbortWithStatusJSON(code, &E{
			Code:    code,
			Message: "server error",
		})
	}
}

func (c *Context) AbortAndWriteInvalidInputError(fields map[string]any) {
	message := "invalid input"
	if len(fields) == 1 {
		for _, v := range fields {
			message = fmt.Sprintf("%v", v)
			break
		}
	}
	c.AbortAndWriteError(http.StatusUnprocessableEntity, &E{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Details: fields,
	})
}

func WrapContext(c *gin.Context) *Context {
	return &Context{c}
}

func WrapHandler(f gin.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		f(c.Context)
	}
}

func WrapHandlers(f []gin.HandlerFunc) []HandlerFunc {
	wrappedHandlers := make([]HandlerFunc, len(f))
	for i, handler := range f {
		wrappedHandlers[i] = WrapHandler(handler)
	}
	return wrappedHandlers
}

func AdaptHandler(f HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(WrapContext(c))
	}
}

func AdaptHandlers(f []HandlerFunc) []gin.HandlerFunc {
	wrappedHandlers := make([]gin.HandlerFunc, len(f))
	for i, handler := range f {
		wrappedHandlers[i] = AdaptHandler(handler)
	}
	return wrappedHandlers
}
