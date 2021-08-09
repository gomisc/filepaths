// Package filepaths - функционал для работы с файлами
package filepaths

import (
	"os"
	"path"
)

const homeDirVariable = "HOME"

// FileExists - проверяет существует ли файл/директория
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

// HomePath - возвращает путь к домашней директории
func HomePath(args ...string) string {
	return path.Join(append([]string{os.Getenv(homeDirVariable)}, args...)...)
}

// GoPath - возвращает GOPATH
func GoPath(args ...string) string {
	var p = os.Getenv("GOPATH")
	if p == "" {
		p = HomePath("go")
	}

	parts := append([]string{p}, args...)

	return path.Join(parts...)
}
