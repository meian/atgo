package repo

import (
	"github.com/meian/atgo/database"
	"github.com/meian/atgo/models"
	"gorm.io/gorm"
)

type RatedType struct {
	*repository[models.RatedType]
}

func NewRateType(db *gorm.DB) *RatedType {
	return NewRatedTypeWithDBConn(database.NewDBConn(db))
}

func NewRatedTypeWithDBConn(dbConn *database.DBConn) *RatedType {
	return &RatedType{newRepositoryWithDBConn[models.RatedType](dbConn)}
}
