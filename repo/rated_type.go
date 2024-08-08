package repo

import (
	"github.com/meian/atgo/database"
	"github.com/meian/atgo/models"
	"github.com/meian/atgo/models/ids"
	"gorm.io/gorm"
)

type RatedType struct {
	*repository[models.RatedType, ids.RatedType]
}

func NewRateType(db *gorm.DB) *RatedType {
	return NewRatedTypeWithDBConn(database.NewDBConn(db))
}

func NewRatedTypeWithDBConn(dbConn *database.DBConn) *RatedType {
	return &RatedType{newRepositoryWithDBConn[models.RatedType, ids.RatedType](dbConn)}
}
