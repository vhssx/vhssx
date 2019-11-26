package secoo

import (
	"github.com/dgrijalva/jwt-go"
)

// Store values if the related meta data is lost?
// Like: SessionID, Version, Landing URI, CreatedTime
type SessionCookieStore struct {
	SessionId string `json:"sid"`

	Extra string `json:"extra"`

	jwt.StandardClaims
}

func NewSessionCookieStore(sessionId, extra string) *SessionCookieStore {
	return &SessionCookieStore{sessionId, extra, jwt.StandardClaims{}}
}
