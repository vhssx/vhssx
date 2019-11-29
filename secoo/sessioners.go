package secoo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zhanbei/static-server/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestLevel int

const (
	LevelFirstTimeRequest RequestLevel = 1
	// The requests from crawlers will carry no cookies, hence the level is always 1.
	LevelCrawlerRequest RequestLevel = 11

	LevelSecondTimeRequest RequestLevel = 2

	LevelFollowingTimeRequest RequestLevel = 12

	LevelUnknownRequest RequestLevel = -1
)

// Store values if the related meta data is lost?
// Like: SessionID, Version, Landing URI, CreatedTime
type SessionCookieStore struct {
	SessionIdHex string `json:"sid" bson:"sid"`
	// The time requested.
	// It could the first time, the second time, or the following time.
	// [ 1 | 2 | 12 ]
	Level RequestLevel `json:"level" bson:"level"`

	PreviousSessionIdHex string `json:"extra" bson:"extra"`

	jwt.StandardClaims `json:"meta,omitempty" bson:"meta,omitempty"`
}

func NewSessionCookieStore(sessionId string, level RequestLevel, extra string) *SessionCookieStore {
	return &SessionCookieStore{sessionId, level, extra, jwt.StandardClaims{
		Id: sessionId,

		IssuedAt: time.Now().Unix(),
	}}
}

type ObjectId = primitive.ObjectID

func (m *SessionCookieStore) GetSessionIds() (sid ObjectId, xid *ObjectId) {
	oid, err := primitive.ObjectIDFromHex(m.SessionIdHex)
	if err != nil {
		fmt.Println("failed to parse the session ID from hex:[", m.SessionIdHex, "].")
		sid = primitive.NewObjectID()
	} else {
		sid = oid
	}
	if m.PreviousSessionIdHex == "" {
		return
	}
	oid, err = primitive.ObjectIDFromHex(m.PreviousSessionIdHex)
	if err != nil {
		fmt.Println("failed to parse the previous session ID from hex:[", m.PreviousSessionIdHex, "].")
	} else {
		xid = &oid
	}
	return
}

const KeySessionCookieStore = "_SESSION_COOKIE_STORE"
const KeySessionCookieId = "_SESSION_COOKIE_ID"

func (m *SessionCookieStore) SerializeToRequest(req *http.Request) {
	req.Header.Set(KeySessionCookieId, m.SessionIdHex)
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
