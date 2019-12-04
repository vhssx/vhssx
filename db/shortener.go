package db

import (
	"github.com/zhanbei/dxb"
)

func LoadShortenerRedirectionRecords(results interface{}) error {
	cur, err := colShortenerRedirections.Find(dxb.NewTimoutContext(5), dxb.M{})
	if err != nil {
		return err
	}
	defer cur.Close(dxb.NewBackgroundContext())
	return cur.All(dxb.DefaultContext(), results)
}
