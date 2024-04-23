package models

type RatedType struct {
	Type string `gorm:"primary_key"`

	Contests []Contest `gorm:"foreignKey:RatedType;constraint:OnDelete:CASCADE"`
}

func init() {
	addMigrateTarget(RatedType{})
}
