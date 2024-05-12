package flags

import (
	"reflect"
)

var (
	Version   string
	CommitSHA string
)

func Package() string {
	type tip int
	return reflect.TypeOf(tip(0)).PkgPath()
}
