package models

import "github.com/meian/atgo/models/ids"

type RatedType struct {
	Type ids.RatedType `gorm:"primary_key"`

	Contests []Contest `gorm:"foreignKey:RatedType;constraint:OnDelete:CASCADE"`
}

func init() {
	addMigrateTarget(RatedType{})
}
