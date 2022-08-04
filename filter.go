package filepaths

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"git.corout.in/golibs/errors"
)

const (
	errIgnoreNotExist = errors.Const("ignore file not exist")
)

// FilesFilter - фильтр файлов
type FilesFilter interface {
	// Name - возвращает имя фильтра
	Name() string
	// Exclude - исключать файлы основанные на параметрах
	// абсолютного пути, начальной директории и системной информации о файле
	Exclude(abspath, base string, fi os.FileInfo) (bool, error)
}

type matchPattern struct {
	pattern *regexp.Regexp
	src     string
	negate  bool
}

// MatchStrings - матчит переданные строки паттерну, если паттерн негативный,
// то при совпадении возвращает ложь
func (mp *matchPattern) MatchStrings(ss ...string) bool {
	for _, s := range ss {
		if s == mp.src {
			return true
		}

		if mp.pattern.MatchString(s) {
			return !(mp.negate)
		}
	}

	return false
}

type matchFilter struct {
	patterns []*matchPattern
}

// MatchFilterFromLines - возвращает фильтр работающий по списку паттернов
func MatchFilterFromLines(lines ...string) (FilesFilter, error) {
	filter := &matchFilter{}

	for _, line := range lines {
		pattern, negate := makePatternFromLine(line)
		if pattern != nil {
			filter.patterns = append(filter.patterns, &matchPattern{
				pattern: pattern, negate: negate, src: line,
			})
		}
	}

	return filter, nil
}

// MatchFilterFomFile - возвращает matchFilter по содержимому файла
func MatchFilterFomFile(path string) (FilesFilter, error) {
	if !FileExists(path) {
		return nil, errIgnoreNotExist
	}

	fd, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "open .gitignore file")
	}

	reader := bufio.NewReader(fd)

	var (
		line  string
		lines []string
	)

	for {
		line, err = reader.ReadString('\n')

		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, errors.Wrap(err, "read file")
			}
		}

		// пропускаем строки комментариев
		if !strings.HasPrefix(line, `#`) {
			// очищаем разделитель строк
			lines = append(lines, strings.TrimSuffix(line, "\n"))
		}

		if err != nil {
			break
		}
	}

	return MatchFilterFromLines(lines...)
}

// Name -
func (f *matchFilter) Name() string {
	return "match-filter"
}

// Exclude -
func (f *matchFilter) Exclude(abs, base string, fi os.FileInfo) (bool, error) {
	rel, err := filepath.Rel(base, abs)
	if err != nil {
		return false, errors.Wrap(err, "get relative path")
	}

	for _, pattern := range f.patterns {
		if pattern.MatchStrings(fi.Name(), rel, abs) {
			return filterResult(true, fi)
		}
	}

	return false, nil
}

func makePatternFromLine(line string) (*regexp.Regexp, bool) {
	// очищаем строку от мусора
	line = strings.TrimSuffix(line, "\r")
	line = strings.Trim(line, " ")

	// если после очистки пустая строка
	if line == "" {
		return nil, false
	}

	negate := false
	if line[0] == '!' {
		negate = true
		line = line[1:]
	}

	// если `#` или `!` экранированы  `\`
	if regexp.MustCompile(`^(\#|\!)`).MatchString(line) {
		line = line[1:]
	}

	// для foo/*.blah в директории, добавляем префикс /
	if regexp.MustCompile(`([^\/+])/.*\*\.`).MatchString(line) && line[0] != '/' {
		line = "/" + line
	}

	// обрабатываем экранированную "."
	line = regexp.MustCompile(`\.`).ReplaceAllString(line, `\.`)

	magicStar := "#$~"

	// обрабатываем "/**/"
	if strings.HasPrefix(line, "/**/") {
		line = line[1:]
	}

	line = regexp.MustCompile(`/\*\*/`).ReplaceAllString(line, `(/|/.+/)`)
	line = regexp.MustCompile(`\*\*/`).ReplaceAllString(line, `(|.`+magicStar+`/)`)
	line = regexp.MustCompile(`/\*\*`).ReplaceAllString(line, `(|/.`+magicStar+`)`)

	// обрабатываем "*"
	line = regexp.MustCompile(`\\\*`).ReplaceAllString(line, `\`+magicStar)
	line = regexp.MustCompile(`\*`).ReplaceAllString(line, `([^/]*)`)

	// обрабатываем экранированный "?"
	line = strings.ReplaceAll(line, "?", `\?`)
	line = strings.ReplaceAll(line, magicStar, "*")

	// собираем выражение
	var expr string
	if strings.HasSuffix(line, "/") {
		expr = line + "(|.*)$"
	} else {
		expr = line + "(|/.*)$"
	}

	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}

	pattern, _ := regexp.Compile(expr)

	return pattern, negate
}

func filterResult(ok bool, fi os.FileInfo) (bool, error) {
	if !ok {
		return false, nil
	}

	if fi.IsDir() {
		return false, filepath.SkipDir
	}

	return true, nil
}
