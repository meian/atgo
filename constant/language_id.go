package constant

import "fmt"

type LanguageID int

const (
	LanguageGo_1_20_6 LanguageID = 5002
)

var idNames = map[LanguageID]string{
	LanguageGo_1_20_6: "Go (1.20.6)",
}

func (id LanguageID) StringValue() string {
	return fmt.Sprint(id)
}

func (id LanguageID) Name() string {
	if idName, ok := idNames[id]; ok {
		return idName
	}
	return "Unknown"
}
