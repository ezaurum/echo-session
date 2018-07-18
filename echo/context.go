package session

import (
	"github.com/labstack/echo"
	"github.com/ezaurum/echo-session"
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
