package filepaths

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"git.corout.in/golibs/errors"
)

// OpenSec - безопасно открывает дескриптор
func OpenSec(name string) (*os.File, error) {
	fd, err := os.Open(filepath.Clean(name))
	if err != nil {
		return nil, errors.Wrap(err, "os open")
	}

	return fd, nil
}

// OpenFileSec - безопасно открывает файл
// nolint: gosec
func OpenFileSec(name string, flag int, perm os.FileMode) (*os.File, error) {
	fd, err := os.OpenFile(filepath.Clean(name), flag, perm)
	if err != nil {
		return nil, errors.Wrap(err, "os open file")
	}

	return fd, nil
}

// ReadFileSec - безопасно считывает содержимое файла по пути
func ReadFileSec(path string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Clean(path))
}
