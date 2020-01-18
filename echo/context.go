package echo

import (
	"github.com/ezaurum/session"
	"github.com/labstack/echo"
)

const (
	DefaultSessionContextKey = "default session context key"
)

func GetSession(c echo.Context) session.Session {
	return c.Get(DefaultSessionContextKey).(session.Session)
}

func SetSession(c echo.Context, s session.Session) {
	c.Set(DefaultSessionContextKey, s)
}
