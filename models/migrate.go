package models

var MigrateTargets []any

func addMigrateTarget(m any) {
	MigrateTargets = append(MigrateTargets, m)
}
