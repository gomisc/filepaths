package filepaths

import (
	"context"
	"io"
	"os"

	"git.corout.in/golibs/errors"
	"git.corout.in/golibs/errors/errgroup"
	"git.corout.in/golibs/fields"
	"git.corout.in/golibs/slog"
)

// CloseAll - закрывает переданные дескрипторы c возвратом ошибок
func CloseAll(closers ...io.Closer) error {
	var err error

	for _, c := range closers {
		if e := c.Close(); e != nil {
			err = errors.And(err, e)
		}
	}

	return err
}

// CloseLogAll - закрывает переданные дескрипторы c логированием ошибок
func CloseLogAll(ctx context.Context, closers ...io.Closer) {
	log := slog.MustFromContext(ctx)

	for _, c := range closers {
		if err := c.Close(); err != nil {
			log.Error("close descriptor", errors.Extract(err))
		}
	}
}

// RemoveLogAll - удаляет переданные файлы и директории с логированием ошибок
func RemoveLogAll(ctx context.Context, paths ...string) {
	log := slog.MustFromContext(ctx)

	for _, p := range paths {
		if err := os.RemoveAll(p); err != nil {
			log.Error("close descriptor", errors.Extract(err))
		}
	}
}

// MakeDirs - асинхронно создает директории из списка
func MakeDirs(names ...string) error {
	var eg = errgroup.New()

	for _, dir := range names {
		d := dir

		eg.Go(func() error {
			if err := os.MkdirAll(d, os.ModePerm); err != nil {
				return errors.Builder(fields.Str("path", d)).
					Wrap(err, "create directory")
			}

			return nil
		})
	}

	return eg.Wait()
}
