package secoo_test

import (
	"fmt"
	"testing"

	"github.com/zhanbei/static-server/secoo"
)

func TestGenerateSessionCookieValue(t *testing.T) {
	gene := secoo.SessionCookieHelper{[]byte("Hola World"), true}

	store := secoo.NewSessionCookieStore("Hello", secoo.LevelFirstTimeRequest, "World")
	token, err := gene.EncodeSessionStore(store)
	fmt.Println(token, err)

}

func TestSessionCookieHelper_ParseSessionStore(t *testing.T) {
	gene := secoo.SessionCookieHelper{[]byte("Hola World"), true}
	store, err := gene.ParseSessionStore("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWQiOiJIZWxsbyIsImV4dHJhIjoiV29ybGQifQ.apts5KBAT-MPcuP893Eh3xm9NNtTlSTGKd68s9kutP0")
	fmt.Println(store, err)
}
