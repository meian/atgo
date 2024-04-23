package constant

//go:generate go run github.com/dmarkham/enumer@latest -type=RatedType -trimprefix=RatedType -transform=kebab
type RatedType int

const (
	RatedTypeAll RatedType = 0
	RatedTypeABC RatedType = 1
	RatedTypeARC RatedType = 2
	RatedTypeAGC RatedType = 3
	RatedTypeAHC RatedType = 4
)
