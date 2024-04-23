package database

import (
	"reflect"

	"gorm.io/gorm"
)

func TableName(db *gorm.DB, model any) string {
	return db.NamingStrategy.TableName(reflect.TypeOf(model).Name())
}
