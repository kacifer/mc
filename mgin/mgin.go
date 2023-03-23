package mgin

import (
	"github.com/gin-gonic/gin"
	"github.com/notpm/mc/mjwt"
	"github.com/notpm/mc/mlog"
	"net/http"
)

// HandlerFunc defines the handler used by gin middleware as return value.
type HandlerFunc func(*Context)

// Engine is the framework's instance, it contains the muxer, middleware and configuration settings.
// Create an instance of Engine, by using New() or Default()
type Engine struct {
	*gin.Engine
}

func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.Engine.NoRoute(AdaptHandlers(handlers)...)
}

func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	engine.Engine.NoMethod(AdaptHandlers(handlers)...)
}

func (engine *Engine) Use(middleware ...HandlerFunc) {
	engine.Engine.Use(AdaptHandlers(middleware)...)
}

func (engine *Engine) Group(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.Group(relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) {
	engine.Engine.Handle(httpMethod, relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) POST(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.POST(relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) GET(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.GET(relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) DELETE(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.DELETE(relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) PATCH(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.PATCH(relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) PUT(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.PUT(relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.OPTIONS(relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) HEAD(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.HEAD(relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) Any(relativePath string, handlers ...HandlerFunc) {
	engine.Engine.Any(relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) Match(methods []string, relativePath string, handlers ...HandlerFunc) {
	engine.Engine.Match(methods, relativePath, AdaptHandlers(handlers)...)
}

func (engine *Engine) Static(relativePath, root string) {
	engine.Engine.Static(relativePath, root)
}

func (engine *Engine) StaticFS(relativePath string, fs http.FileSystem) {
	engine.Engine.StaticFS(relativePath, fs)
}

func (engine *Engine) StaticFile(relativePath, filepath string) {
	engine.Engine.StaticFile(relativePath, filepath)
}

func (engine *Engine) StaticFileFS(relativePath, filepath string, fs http.FileSystem) {
	engine.Engine.StaticFileFS(relativePath, filepath, fs)
}

func New() *Engine {
	return &Engine{gin.New()}
}

func Default() *Engine {
	return &Engine{gin.Default()}
}

var HealthCheckPaths = []string{"/healthz", "/api/v1/healthz"}

const AuthLoginPath = "/api/v1/auth/login"
const AuthRefreshPath = "/api/v1/auth/refresh"
const AuthUserPath = "/api/v1/auth/user"
const SettingGetPath = "/api/v1/settings"
const SettingSetPath = "/api/v1/settings"

type CustomAuthConfig struct {
	Jwt           mjwt.Engine
	SkipAuthPaths []string
	UserStore     UserStore
	SettingStore  SettingStore
}

type CustomConfig struct {
	Version string
	Auth    *CustomAuthConfig
}

func Custom(config CustomConfig) *Engine {
	base := gin.New()

	engine := &Engine{base}

	var middlewares []HandlerFunc
	if config.Auth != nil {
		if config.Auth.Jwt != nil {
			skipAuthPaths := config.Auth.SkipAuthPaths
			for _, path := range HealthCheckPaths {
				skipAuthPaths = append(skipAuthPaths, path)
			}
			skipAuthPaths = append(skipAuthPaths, AuthLoginPath)
			middlewares = append(middlewares, CreateAuthMiddleware(config.Auth.Jwt, skipAuthPaths))
		}
	}
	middlewares = append(middlewares, WrapHandler(gin.LoggerWithConfig(gin.LoggerConfig{
		Output:    mlog.DefaultLogger.Out,
		SkipPaths: HealthCheckPaths,
	})))
	middlewares = append(middlewares, WrapHandler(gin.Recovery()))

	engine.Use(middlewares...)

	healthzHandler := CreateHealthzHandler(config.Version)
	for _, path := range HealthCheckPaths {
		engine.Any(path, healthzHandler)
	}

	if config.Auth != nil {
		if config.Auth.Jwt != nil && config.Auth.UserStore != nil {
			engine.POST(AuthLoginPath, CreateAuthLoginHandler(config.Auth.Jwt, config.Auth.UserStore))
			engine.GET(AuthRefreshPath, CreateAuthRefreshHandler(config.Auth.Jwt))
			engine.GET(AuthUserPath, CreateAuthUserHandler(config.Auth.UserStore))
		}

		if config.Auth.SettingStore != nil {
			engine.GET(SettingGetPath, CreateAuthSettingGetHandler(config.Auth.SettingStore))
			engine.PUT(SettingSetPath, CreateAuthSettingSetHandler(config.Auth.SettingStore))
		}
	}

	engine.HandleMethodNotAllowed = true

	engine.NoRoute(func(c *Context) {
		c.AbortAndWriteError(http.StatusNotFound, "page not found")
	})
	engine.NoMethod(func(c *Context) {
		c.AbortAndWriteError(http.StatusMethodNotAllowed, "method not allowed")
	})

	return engine
}
