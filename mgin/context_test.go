package mgin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func CreateTestContext() (
	*httptest.ResponseRecorder,
	*Context,
	*Engine,
) {
	gin.SetMode(gin.ReleaseMode)
	recorder := httptest.NewRecorder()
	gc, ge := gin.CreateTestContext(recorder)
	c := &Context{gc}
	e := &Engine{ge}
	return recorder, c, e
}

func TestMustUintContext(t *testing.T) {
	assertions := require.New(t)
	_, c, _ := CreateTestContext()
	c.Set(IDKey, uint(1))
	assertions.Equal(uint(1), c.MustUintContext(IDKey))
	c.Set(IDKey, float64(1))
	assertions.Equal(uint(1), c.MustUintContext(IDKey))
	c.Set(IDKey, "1")
	assertions.Equal(uint(1), c.MustUintContext(IDKey))
}

func TestUintContext(t *testing.T) {
	assertions := require.New(t)
	_, c, _ := CreateTestContext()
	c.Set(IDKey, uint(1))
	id, exist := c.UintContext(IDKey)
	assertions.True(exist)
	assertions.Equal(uint(1), id)
}

func TestMustIDContext(t *testing.T) {
	assertions := require.New(t)
	_, c, _ := CreateTestContext()
	c.Set(IDKey, uint(1))
	assertions.Equal(uint(1), c.MustIDContext())
}

func TestContext_IDContext(t *testing.T) {
	assertions := require.New(t)
	_, c, _ := CreateTestContext()
	c.Set(IDKey, uint(1))
	id, exist := c.IDContext()
	assertions.True(exist)
	assertions.Equal(uint(1), id)
}

func TestContext_UintQuery(t *testing.T) {
	assertions := require.New(t)
	_, c, _ := CreateTestContext()
	c.Request = &http.Request{}
	c.Request.URL, _ = url.Parse("http://localhost:8080/?id=1")
	assertions.Equal(uint(1), c.UintQuery(IDQuery))
}

func TestContext_IDQuery(t *testing.T) {
	assertions := require.New(t)
	_, c, _ := CreateTestContext()
	c.Request = &http.Request{}
	c.Request.URL, _ = url.Parse("http://localhost:8080/?id=1")
	assertions.Equal(uint(1), c.IDQuery())
}

func TestContext_UintParam(t *testing.T) {
	assertions := require.New(t)
	_, c, _ := CreateTestContext()
	c.Params = gin.Params{gin.Param{Key: IDParam, Value: "1"}}
	assertions.Equal(uint(1), c.UintParam(IDParam))
}

func TestContext_IDParam(t *testing.T) {
	assertions := require.New(t)
	_, c, _ := CreateTestContext()
	c.Params = gin.Params{gin.Param{Key: IDParam, Value: "1"}}
	assertions.Equal(uint(1), c.IDParam())
}

func TestContext_AbortAndWriteError(t *testing.T) {
	assertions := require.New(t)
	recorder, c, _ := CreateTestContext()
	c.AbortAndWriteError(http.StatusBadRequest, errors.New("bad request"))
	assertions.Equal(http.StatusBadRequest, c.Writer.Status())
	body, err := io.ReadAll(recorder.Result().Body)
	assertions.Nil(err)
	expected, err := json.Marshal(&E{
		Code:    http.StatusBadRequest,
		Message: "bad request",
	})
	assertions.Nil(err)
	assertions.Equal(string(expected), string(body))
}

func TestContext_AbortAndWriteInternalError(t *testing.T) {
	assertions := require.New(t)
	{
		recorder, c, _ := CreateTestContext()
		c.AbortAndWriteInternalError(http.StatusInternalServerError, errors.New("database query error"))
		assertions.Equal(http.StatusInternalServerError, c.Writer.Status())
		body, err := io.ReadAll(recorder.Result().Body)
		assertions.Nil(err)
		expected, err := json.Marshal(&E{
			Code:    http.StatusInternalServerError,
			Message: "server error",
		})
		assertions.Nil(err)
		assertions.Equal(string(expected), string(body))
	}
	{
		recorder, c, _ := CreateTestContext()
		gin.SetMode(gin.DebugMode)
		c.AbortAndWriteInternalError(http.StatusInternalServerError, errors.New("database query error"))
		assertions.Equal(http.StatusInternalServerError, c.Writer.Status())
		body, err := io.ReadAll(recorder.Result().Body)
		assertions.Nil(err)
		expected, err := json.Marshal(&E{
			Code:    http.StatusInternalServerError,
			Message: "database query error",
		})
		assertions.Nil(err)
		assertions.Equal(string(expected), string(body))
	}
}

func TestContext_AbortAndWriteInvalidInputError(t *testing.T) {
	assertions := require.New(t)
	recorder, c, _ := CreateTestContext()
	c.AbortAndWriteInvalidInputDetails(map[string]any{
		"username": "username is required",
	})
	assertions.Equal(http.StatusUnprocessableEntity, c.Writer.Status())
	body, err := io.ReadAll(recorder.Result().Body)
	assertions.Nil(err)
	expected, err := json.Marshal(&E{
		Code:    http.StatusUnprocessableEntity,
		Message: "username is required",
		Details: map[string]any{
			"username": "username is required",
		},
	})
	assertions.Nil(err)
	assertions.Equal(string(expected), string(body))
}
