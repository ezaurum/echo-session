package jarmiddle

import (
	"github.com/ezaurum/remember/cookie"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"net/http"
)

const cookieJarContextKey = "cookie.jar.context-key"

var (
	JarIsNotPresentErr = errors.New("jar is not present error")
)

func Remove(c echo.Context, cookieName string) error {
	jar := c.Get(cookieJarContextKey).(cookie.Jar)
	if nil == jar {
		return JarIsNotPresentErr
	}
	jar.Remove(cookieName)
	return nil
}

func Set(c echo.Context, setCookie *http.Cookie) error {
	jar := c.Get(cookieJarContextKey).(cookie.Jar)
	if nil == jar {
		return JarIsNotPresentErr
	}
	jar.Set(setCookie)
	return nil
}

func Get(c echo.Context, cookieName string) (*http.Cookie, error) {
	jar := c.Get(cookieJarContextKey).(cookie.Jar)
	if nil == jar {
		return nil, JarIsNotPresentErr
	}
	return jar.Get(cookieName), nil
}

func Middleware(config Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}
			jar := cookie.New(c.Request())
			c.Set(cookieJarContextKey, jar)
			if err := next(c); nil != err {
				return err
			}
			jar.Write(c.Response())
			return nil
		}
	}
}

type Config struct {
	Skipper middleware.Skipper
}
