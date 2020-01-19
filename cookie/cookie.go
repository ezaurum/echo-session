package cookie

import (
	"net/http"
	"time"
)

// 기본 쿠키 관리자 인터페이스
type Jar interface {
	Set(cookie *http.Cookie)
	Remove(cookieName string)
	Get(cookieName string) *http.Cookie
	Write(responseWriter http.ResponseWriter)
}

// 쿠키 관리자 구현
type jar struct {
	request  map[string]*http.Cookie
	response map[string]*http.Cookie
}

func (j *jar) Set(cookie *http.Cookie) {
	j.response[cookie.Name] = cookie
}

func (j *jar) Remove(cookieName string) {
	if oldCookie, b := j.request[cookieName]; b {
		n := *oldCookie
		if n.MaxAge != 0 || len(n.RawExpires) > 0 {
			n.MaxAge = 0
			n.Expires = time.Now().Add(-time.Hour)
		}
		n.Value = ""
		j.response[cookieName] = &n
	} else {
		delete(j.response, cookieName)
	}
}

func (j *jar) Get(cookieName string) *http.Cookie {
	if ck, b := j.response[cookieName]; b {
		return ck
	}
	if ck, b := j.request[cookieName]; b {
		return ck
	}
	return nil
}

// response에 쓰기
func (j *jar) Write(responseWriter http.ResponseWriter) {
	for _, ck := range j.response {
		http.SetCookie(responseWriter, ck)
	}
}

func New(request *http.Request) Jar {
	jar := jar{
		request:  make(map[string]*http.Cookie),
		response: make(map[string]*http.Cookie),
	}
	for _, c := range request.Cookies() {
		jar.request[c.Name] = c
	}
	return &jar
}
