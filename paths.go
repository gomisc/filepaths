package filepaths

import (
	"os"
	"path"
)

const (
	homeDirVariable = "HOME"
	configDirName   = ".config"
	cacheDirName    = ".cache"
)

// FileExists - проверяет существовать ли файл/директория
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

// GoPath - возвращает GOPATH
func GoPath(args ...string) string {
	p := os.Getenv("GOPATH")

	if p == "" {
		p = HomePath("go")
	}

	parts := append([]string{p}, args...)

	return path.Join(parts...)
}

// HomePath - возвращает путь к домашней директории пользователя
func HomePath(args ...string) string {
	return path.Join(append([]string{os.Getenv(homeDirVariable)}, args...)...)
}

// ConfigPath - возвращает путь к директории конфигов пользователя
func ConfigPath(args ...string) string {
	return HomePath(append([]string{configDirName}, args...)...)
}

// CachePath - возвращает путь к директории кэша пользователя
func CachePath(args ...string) string {
	return HomePath(append([]string{cacheDirName}, args...)...)
}
