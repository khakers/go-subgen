package projectpath

// Thanks to stackoverflow for this code
// https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../..")

	// TestOut output folder of this project
	TestOut = filepath.Join(Root, "test", "out")

	// TestData folder of this project
	TestData = filepath.Join(Root, "test", "testdata")

	// TestSamples Folder containing sample files
	TestSamples = filepath.Join(Root, "test", "samples")

	TestModels = filepath.Join(Root, "test", "models")

	Test = filepath.Join(Root, "test")
)
