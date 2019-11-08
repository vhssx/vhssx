package db

import (
	"context"
	"time"
)

func NewTimoutContext(second int64) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(second*int64(time.Second)))
	return ctx
}

func newCrudContext() context.Context {
	return NewTimoutContext(5)
}
