package mgin

import (
	"github.com/notpm/mc/mjwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func CreateAuthLoginHandler(jwt mjwt.Engine, userStore UserStore) HandlerFunc {
	return func(c *Context) {
		type Data struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var data Data
		if !c.MustBindJSON(&data) {
			return
		}

		user, err := userStore.FindByUsername(data.Username)
		if err != nil {
			if err == ErrUsernameNotFound {
				c.AbortAndWriteInvalidInputError(map[string]any{
					"username": "username not exist",
				})
				return
			}
			c.AbortWithInternalError(500, errors.Wrap(err, "find user by username error"))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(data.Password))
		if err != nil {
			c.AbortAndWriteInvalidInputError(map[string]any{
				"password": "password not match",
			})
			return
		}
		tokenString, err := jwt.SignedStringForID(user.GetID())
		if err != nil {
			c.AbortWithInternalError(500, errors.Wrap(err, "sign error"))
			return
		}
		c.Header("Authorization", tokenString)
		c.JSON(200, "OK")
	}
}

func CreateAuthRefreshHandler(jwt mjwt.Engine) HandlerFunc {
	return func(c *Context) {
		userID := c.MustIDContext()

		tokenString, err := jwt.SignedStringForID(userID)
		if err != nil {
			c.AbortWithInternalError(500, errors.Wrap(err, "sign error"))
			return
		}
		c.Header("Authorization", tokenString)
		c.JSON(200, "OK")
	}
}

func CreateAuthUserHandler(userStore UserStore) HandlerFunc {
	return func(c *Context) {
		id := c.MustIDContext()

		user, err := userStore.Find(id)
		if err != nil {
			if err == ErrUserIDNotFound {
				c.AbortAndWriteInvalidInputError(map[string]any{
					"id": "user not found",
				})
				return
			}
			c.AbortWithInternalError(500, errors.Wrap(err, "find user error"))
			return
		}
		c.JSON(200, user)
	}
}

func CreateAuthSettingGetHandler(settingStore SettingStore) HandlerFunc {
	return func(c *Context) {
		id := c.MustIDContext()

		setting, err := settingStore.Get(id, c.Query("key"))
		if err != nil {
			c.AbortWithInternalError(500, errors.Wrap(err, "get setting error"))
			return
		}

		c.JSON(200, setting)
	}
}

func CreateAuthSettingSetHandler(settingStore SettingStore) HandlerFunc {
	return func(c *Context) {
		type Data struct {
			Value string `json:"value"`
		}
		var data Data
		if !c.MustBindJSON(&data) {
			return
		}

		id := c.MustIDContext()

		if err := settingStore.Set(id, c.Query("key"), data.Value); err != nil {
			c.AbortWithInternalError(500, errors.Wrap(err, "get setting error"))
			return
		}

		c.JSON(200, "OK")
	}
}
