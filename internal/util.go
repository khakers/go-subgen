package internal

import "path/filepath"

func GetFileName(file string) string {
	basepath := filepath.Base(file)
	return basepath[:len(basepath)-len(filepath.Ext(basepath))]
}
