package secoo

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/zhanbei/static-server/conf"
)

func NewSessionCookieHelper(ops *conf.OptionsSessionCookie) *SessionCookieHelper {
	return &SessionCookieHelper{[]byte(ops.Secret), ops.AllSubDomains}
}

func (m *SessionCookieHelper) EncodeSessionStore(claims *SessionCookieStore) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(m.Secret)
	return ss, err
}

// sample token is expired.  override time so it parses as valid
func (m *SessionCookieHelper) ParseSessionStore(token string) (*SessionCookieStore, error) {
	store := new(SessionCookieStore)
	result, err := jwt.ParseWithClaims(token, store, func(token *jwt.Token) (interface{}, error) {
		return m.Secret, nil
	})
	if err != nil {
		fmt.Println("[JWT] Token Parsing Failed:[", err, "][", token, "]")
		return nil, err
	}

	if result == nil || !result.Valid {
		fmt.Println("[JWT] Token Validating Failed:", result)
		return nil, errors.New("validating jwt failed")
	}
	return store, nil
}
