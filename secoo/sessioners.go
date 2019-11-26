package secoo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zhanbei/static-server/utils"
)

type RequestLevel int

const (
	LevelFirstTimeRequest RequestLevel = 1

	LevelSecondTimeRequest RequestLevel = 2

	LevelFollowingTimeRequest RequestLevel = 12
)

// Store values if the related meta data is lost?
// Like: SessionID, Version, Landing URI, CreatedTime
type SessionCookieStore struct {
	SessionId string `json:"sid"`
	// The time requested.
	// It could the first time, the second time, or the following time.
	// [ 1 | 2 | 12 ]
	Level RequestLevel `json:"level"`

	Extra string `json:"extra"`

	jwt.StandardClaims
}

func NewSessionCookieStore(sessionId string, level RequestLevel, extra string) *SessionCookieStore {
	return &SessionCookieStore{sessionId, level, extra, jwt.StandardClaims{
		Id: sessionId,

		IssuedAt: time.Now().Unix(),
	}}
}

const KeySessionCookieStore = "_SESSION_COOKIE_STORE"
const KeySessionCookieId = "_SESSION_COOKIE_ID"

func (m *SessionCookieStore) SerializeToRequest(req *http.Request) {
	req.Header.Set(KeySessionCookieId, m.SessionId)
	bts, err := json.Marshal(m)
	if err != nil {
		// FIX-ME Handle hidden errors globally as warnings(no runtime panics :).
		fmt.Println("Failed to serialize session cookie store:", err)
	}
	req.Header.Set(KeySessionCookieStore, string(bts))
}

func GetRequestSessionId(req *http.Request) string {
	return req.Header.Get(KeySessionCookieId)
}

func RestoreFromRequest(req *http.Request) *SessionCookieStore {
	bts := req.Header.Get(KeySessionCookieStore)
	if !utils.NotEmpty(bts) {
		fmt.Println("failed to fetch session cookie store")
	}
	store := new(SessionCookieStore)
	err := json.Unmarshal([]byte(bts), store)
	if err != nil {
		return nil
	}
	return store
}
