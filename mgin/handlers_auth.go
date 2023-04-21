package mgin

import (
	"github.com/kacifer/mc"
	"github.com/kacifer/mc/mjwt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
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
				c.AbortAndWriteInvalidInputDetails(map[string]any{
					"username": "username not exist",
				})
				return
			}
			c.AbortAndWriteInternalError(http.StatusInternalServerError, errors.Wrap(err, "find user by username error"))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(data.Password))
		if err != nil {
			c.AbortAndWriteInvalidInputDetails(map[string]any{
				"password": "password not match",
			})
			return
		}
		tokenString, err := jwt.SignedStringForID(user.GetID())
		if err != nil {
			c.AbortAndWriteInternalError(http.StatusInternalServerError, errors.Wrap(err, "sign error"))
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
			c.AbortAndWriteInternalError(http.StatusInternalServerError, errors.Wrap(err, "sign error"))
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
				c.AbortAndWriteInvalidInputDetails(map[string]any{
					"id": "user not found",
				})
				return
			}
			c.AbortAndWriteInternalError(http.StatusInternalServerError, errors.Wrap(err, "find user error"))
			return
		}
		c.JSON(200, user)
	}
}

func CreateAuthSettingGetHandler(settingStore SettingStore, keysWhitelist []string) HandlerFunc {
	return func(c *Context) {
		key := c.Query("key")

		var keys []string
		if key == "" {
			for _, k := range strings.Split(c.Query("keys"), ",") {
				if strings.TrimSpace(k) != "" {
					keys = append(keys, strings.TrimSpace(k))
				}
			}
		} else {
			keys = []string{key}
		}

		for _, k := range keys {
			if len(keysWhitelist) > 0 && !mc.SliceContains(keysWhitelist, k) {
				c.AbortAndWriteInvalidInputDetails(map[string]any{
					"key": "key not allowed",
				})
				return
			}
		}

		id := c.MustIDContext()

		settings := map[string]string{}
		for _, k := range keys {
			setting, err := settingStore.Get(id, k)
			if err != nil {
				c.AbortAndWriteInternalError(http.StatusInternalServerError, errors.Wrap(err, "get setting error"))
				return
			}
			settings[k] = setting
		}

		c.JSON(200, settings)
	}
}

func CreateAuthSettingSetHandler(settingStore SettingStore, keysWhitelist []string) HandlerFunc {
	return func(c *Context) {
		key := c.Query("key")

		type Data struct {
			Value string `json:"value"`
		}
		var data Data
		if !c.MustBindJSON(&data) {
			return
		}

		if len(keysWhitelist) > 0 && !mc.SliceContains(keysWhitelist, key) {
			c.AbortAndWriteInvalidInputDetails(map[string]any{
				"key": "key not allowed",
			})
			return
		}

		id := c.MustIDContext()

		if err := settingStore.Set(id, key, data.Value); err != nil {
			c.AbortAndWriteInternalError(http.StatusInternalServerError, errors.Wrap(err, "get setting error"))
			return
		}

		c.JSON(200, "OK")
	}
}
