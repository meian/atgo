package constant

//go:generate go run github.com/dmarkham/enumer@latest -type=ContestCategory -trimprefix=ContestCategory -transform=kebab
type ContestCategory int

const (
	CategoryTypical       ContestCategory = 6
	CategoryPAST          ContestCategory = 50
	CategoryDailyTraining ContestCategory = 60
	CategoryUnrated       ContestCategory = 101
	CategoryJOI           ContestCategory = 200
	CategoryEPFinal       ContestCategory = 1000
	CategoryEPOpenRated   ContestCategory = 1001
	CategoryEPOpenUnrated ContestCategory = 1002
	CategoryEPABC         ContestCategory = 1005
	CategoryEPARC         ContestCategory = 1004
	CategoryHeuristic     ContestCategory = 1200
	CategoryEPHeuristic   ContestCategory = 1250
)
