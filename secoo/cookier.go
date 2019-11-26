package secoo

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (m *SessionCookieHelper) ParseCookie4SessionStore(req *http.Request, key string) (*SessionCookieStore, error) {
	ck, err := req.Cookie(key)
	if err != nil {
		return nil, err
	}
	if ck == nil {
		return nil, errors.New("target cookie is not found")
	}
	return m.ParseSessionStore(ck.Value)
}

func (m *SessionCookieHelper) NewSessionCookieTokenValue(level RequestLevel, extra string) (*SessionCookieStore, string) {
	// FIX-ME A step further, validate with the bare session cookie as chained cookies.
	store := NewSessionCookieStore(NewSessionId(), level, extra)
	token, err := m.EncodeSessionStore(store)
	if err != nil {
		fmt.Println("[JWT] Token Encoding Failed:", err, store)
	}
	return store, token
}

// The host
func (m *SessionCookieHelper) NewHttpCookie(key string, host, value string) *http.Cookie {
	if !m.SubDomains {
		host = ""
	}
	return &http.Cookie{
		Name: key,

		Value: value,

		Path: "/",
		//RawExpires: "",
		//MaxAge:     0,
		// @see https://stackoverflow.com/questions/3290424/set-a-cookie-to-never-expire
		Expires: time.Unix(2147483647, 0),
		//Secure: false,
		//SameSite:   0,
		HttpOnly: true,
		// @see https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies
		Domain: host,
		//Raw: "",
		//Unparsed: nil,
	}
}
