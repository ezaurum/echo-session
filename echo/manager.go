package session

import (
	"time"
	"github.com/ezaurum/cthulthu/generators/snowflake"
	"github.com/labstack/echo"
	"net/http"
	"bitbucket.org/congkong-revivals/congkong/cookie"
	"github.com/ezaurum/echo-session"
	"github.com/ezaurum/echo-session/memstore"
	"github.com/labstack/gommon/random"
)

const (
	IDCookieName          = "session-id-cookie-name-echo-session"
	DefaultSessionExpires = 60 * 15
)

type Manager struct {
	store               session.Store
	random              *random.Random
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
			err := next(c)
			if nil != err {
				return err
			}
			// auto extend
			ca.store.Set(se)
			return nil
		}
	}
}

func (ca Manager) CreateSession(c echo.Context) session.Session {
	// created
	session := ca.store.GetNew(c.RealIP(), c.Request().UserAgent())

	return session
}

func (ca *Manager) ActivateSession(c echo.Context, s session.Session) {
	//refresh session expires
	SetSession(c, s)
	ca.SetSessionIDCookie(c, s)
}

func (ca Manager) FindSession(c echo.Context) (session.Session, bool) {
	sessionIDCookie, e := c.Cookie(ca.sessionIDCookieName)

	if nil != e {
		return nil, true
	}

	s, sessionExist := ca.store.Get(sessionIDCookie.Value)
	if !sessionExist {
		// 세션 유효하지 않은 경우, 만료되었거나, 값 조작이거나
		// 해당 쿠키 삭제
		cookie.ClearCookie(c, ca.sessionIDCookieName)
		return nil, true
	}
	return s, false
}

func (ca Manager) SetSessionIDCookie(c echo.Context, session session.Session) {
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
	manager := &Manager{
		store:               memstore.NewStore(snowflake.New(node), duration, duration*2),
		MaxAge:              expiresInSeconds,
		sessionIDCookieName: IDCookieName,
		random:random.New(),
	}
	return manager
}
