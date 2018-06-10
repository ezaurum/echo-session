package session

import (
	"github.com/labstack/echo"
)

const (
	DefaultSessionContextKey = "default session context key"
)

func GetSession(c echo.Context) Session {
	return c.Get(DefaultSessionContextKey).(Session)
}

func SetSession(c echo.Context, s Session) {
	c.Set(DefaultSessionContextKey, s)
}
