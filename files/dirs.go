package files

import "path/filepath"

func TestDataDir(dir string) string {
	return filepath.Join(dir, "testdata")
}

func SamplesDir(dir string) string {
	return filepath.Join(TestDataDir(dir), "samples")
}
