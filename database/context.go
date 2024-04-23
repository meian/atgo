package database

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type key string

const dbKey key = "db"

func FromContext(ctx context.Context) *gorm.DB {
	v := ctx.Value(dbKey)
	if v == nil {
		panic(errors.New("DB is not initialized"))
	}
	return v.(*gorm.DB)
}

func ContextWith(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}
