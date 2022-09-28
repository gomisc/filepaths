package filepaths

import (
	"io/fs"
	"os"
	"path/filepath"

	"git.corout.in/golibs/errors"
)

// MakeFilesMap - создает карту директории, с возможностью
// последовательного применения фильтров файлов
func MakeFilesMap(base string, withRelPath bool, filters ...FilesFilter) (map[string]os.FileInfo, error) {
	fm := make(map[string]os.FileInfo)

	err := filepath.Walk(
		base, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path == base {
				return nil
			}

			for _, filter := range filters {
				var ok bool

				ok, err = filter.Exclude(path, base, info)
				if err != nil {
					if !errors.Is(err, filepath.SkipDir) {
						return errors.Wrapf(err, "apply filter %s", filter.Name())
					}

					return nil
				}

				if ok {
					return nil
				}
			}

			fileKey := path
			if withRelPath {
				fileKey, err = filepath.Rel(base, path)
				if err != nil {
					return errors.Wrapf(err, "get relative path of %s", path)
				}
			}

			fm[fileKey] = info

			return nil
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "walk directory")
	}

	return fm, nil
}
