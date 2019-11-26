package secoo

import (
	"fmt"
	"net/http"

	"github.com/zhanbei/static-server/conf"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	KeyBareSessionCookieId = "_bs0d"

	KeyValidatedSessionCookieId = "_vs1d"
)

func NewSessionId() string {
	return primitive.NewObjectID().Hex()
}

func HandlerSetSessionCookie(next http.Handler, ops *conf.OptionsSessionCookie) http.HandlerFunc {
	helper := NewSessionCookieHelper(ops)
	return func(w http.ResponseWriter, req *http.Request) {
		ck, store := helper.HandleSessionCookie(req)
		if ck != nil {
			http.SetCookie(w, ck)
			fmt.Println("Setting cookie:", ck.String(), ck)
		}
		if store != nil {
			// Basically the store is never nil.
			// Set for recorders/loggers.
			store.SerializeToRequest(req)
		}
		next.ServeHTTP(w, req)
	}
}

type SessionCookieHelper struct {
	Secret []byte `json:"-"`

	SubDomains bool `json:"subDomains"`
}

// Initially set a bare session cookie; set another cookie as the validated one, if the bare session cookie is sending back,
// because [crawlers often support no cookies](https://webmasters.stackexchange.com/questions/59652/what-happens-if-i-try-to-set-a-cookie-on-a-bot).
func (m *SessionCookieHelper) HandleSessionCookie(req *http.Request) (*http.Cookie, *SessionCookieStore) {
	// Check and set cookies.
	//cks := req.Cookies()
	store, err := m.ParseCookie4SessionStore(req, KeyValidatedSessionCookieId)
	if err == nil && store != nil && store.Level == 2 {
		// C. The following requests with a validated cookie.
		// Found validated session store.
		store.Level = LevelFollowingTimeRequest
		return nil, store
	}
	store, err = m.ParseCookie4SessionStore(req, KeyBareSessionCookieId)
	if err == nil && store != nil {
		// B. The second request with a bare cookie.
		// The remote device is validated as a real browser instead of a crawler.
		// Hence generate a validated cookie based on the bare cookie.
		nextStore, nextToken := m.NewSessionCookieTokenValue(LevelSecondTimeRequest, store.SessionId)
		return m.NewHttpCookie(KeyValidatedSessionCookieId, req.Host, nextToken), nextStore
	}
	// A. The first request without a single cookie.
	store, token := m.NewSessionCookieTokenValue(LevelFirstTimeRequest, "")
	return m.NewHttpCookie(KeyBareSessionCookieId, req.Host, token), store
}
