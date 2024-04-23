package util

import (
	"fmt"

	"gopkg.in/guregu/null.v3"
)

func NullIntString(ni null.Int, defValue ...string) string {
	if len(defValue) == 0 {
		defValue = []string{"-"}
	}
	if ni.Valid {
		return fmt.Sprint(ni.Int64)
	}
	return defValue[0]
}
