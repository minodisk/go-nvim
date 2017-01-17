package window

import "regexp"

var (
	rSpace = regexp.MustCompile(`([[:space:]])`)
)

func Escape(path string) string {
	return rSpace.ReplaceAllString(path, `\\$1`)
}
