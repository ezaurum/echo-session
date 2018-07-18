package session

import (
	"time"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/labstack/echo"
	"net/http"
	"bitbucket.org/congkong-revivals/congkong/cookie"
)

const (
	IDCookieName          = "session-id-cookie-name-echo-session"
	DefaultSessionExpires = 60 * 15
)

type Manager struct {
	store               Store
	MaxAge              int
	sessionIDCookieName string
}

func (ca *Manager) Handler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			se, needSession := ca.FindSession(c)

			if needSession {
				se = ca.CreateSession(c)
			}

			ca.ActivateSession(c, se)
			return next(c)
		}
	}
}

func (ca Manager) CreateSession(c echo.Context) Session {
	// created
	session := ca.store.GetNew()

	return session
}

func (ca *Manager) ActivateSession(c echo.Context, s Session) {
	//refresh session expires
	SetSession(c, s)
	ca.SetSessionIDCookie(c, s)
}

func (ca Manager) FindSession(c echo.Context) (Session, bool) {
	sessionIDCookie, e := c.Cookie(ca.sessionIDCookieName)

	if nil != e {
		return nil, true
	}

	//TODO secure
	s, sessionExist := ca.store.Get(sessionIDCookie.Value)

	if !sessionExist {
		// 세션 유효하지 않은 경우, 만료되었거나, 값 조작이거나
		// 해당 쿠키 삭제
		cookie.ClearCookie(c, ca.sessionIDCookieName)
		return nil, true
	}
	return s, false
}

func (ca Manager) SetSessionIDCookie(c echo.Context, session Session) {
	cookie := http.Cookie{
		Name:     ca.sessionIDCookieName,
		Value:    session.Key(),
		MaxAge:   ca.MaxAge,
		Domain:   "",
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}
	c.SetCookie(&cookie)
}

func Default() *Manager {
	return NewMem(0, DefaultSessionExpires)
}

func NewMem(node int64, expiresInSeconds int) *Manager {
	duration := time.Duration(expiresInSeconds) * time.Second
	k := snowflake.New(node)
	manager := &Manager{
		store:               NewStore(k, duration, duration*2),
		MaxAge:              expiresInSeconds,
		sessionIDCookieName: IDCookieName,
	}
	return manager
}
