package models

type RatedType struct {
	Type string `gorm:"primary_key"`
}

func init() {
	addMigrateTarget(RatedType{})
}
