package files

import "path/filepath"

func MainFile(dir string) string {
	return filepath.Join(dir, "main.go")
}

func TestFile(dir string) string {
	return filepath.Join(dir, "main_test.go")
}

func ModFile(dir string) string {
	return filepath.Join(dir, "go.mod")
}

func SumFile(dir string) string {
	return filepath.Join(dir, "go.sum")
}

func TaskInfoFile(dir string) string {
	return filepath.Join(dir, "task-info.yaml")
}

func TaskLocalFile(dir string) string {
	return filepath.Join(dir, "task-local.yaml")
}
