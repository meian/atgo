package flags

var (
	Version string
)

func init() {
	if Version == "" {
		Version = "no version, please set with ldflags"
	}
}
