package mgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/notpm/mc"
	"github.com/notpm/mc/mlog"
	"github.com/pkg/errors"
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
	v := c.MustGet(key)
	switch v.(type) {
	case int:
		return uint(v.(int))
	case int8:
		return uint(v.(int8))
	case int16:
		return uint(v.(int16))
	case int32:
		return uint(v.(int32))
	case int64:
		return uint(v.(int64))
	case uint:
		return v.(uint)
	case uint8:
		return uint(v.(uint8))
	case uint16:
		return uint(v.(uint16))
	case uint32:
		return uint(v.(uint32))
	case uint64:
		return uint(v.(uint64))
	case float32:
		return uint(v.(float32))
	case float64:
		return uint(v.(float64))
	default:
		return mc.StringToUint(fmt.Sprintf("%v", v))
	}
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

// AbortAndWriteError aborts the context and write standard error response
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

// AbortAndWriteInternalError aborts the context, push error to error stack and write standard error response,
// if gin.IsDebugging() is false, it also hides the error message.
func (c *Context) AbortAndWriteInternalError(code int, err any) {
	mlog.Error("request abort with code %v, error: %v", code, err)
	if gin.IsDebugging() {
		c.AbortAndWriteError(code, err)
		return
	} else {
		// hide error message if not debugging
		c.AbortAndWriteError(code, "server error")
	}
}

func (c *Context) AbortAndWriteInternalServerError(err any) {
	c.AbortAndWriteInternalError(http.StatusInternalServerError, err)
}

// AbortAndWriteInvalidInputError aborts the context and write standard error response with invalid input error
func (c *Context) AbortAndWriteInvalidInputError(e *E) {
	c.AbortAndWriteError(http.StatusUnprocessableEntity, e)
}

// AbortAndWriteInvalidInputDetails aborts the context and write standard error response with invalid input details
func (c *Context) AbortAndWriteInvalidInputDetails(details map[string]any) {
	message := "invalid input"
	if len(details) > 0 {
		for _, v := range details {
			message = fmt.Sprintf("%v", v)
			break
		}
	}
	c.AbortAndWriteError(http.StatusUnprocessableEntity, &E{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Details: details,
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
