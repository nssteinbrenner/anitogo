package anitogo

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	indexTooLargeErr    = "index is too large for array or string"
	indexTooSmallErr    = "index is less than 0"
	endIndexTooSmallErr = "end index is smaller than begin index"
	tokensEmptyErr      = "tokens are empty"
	tokenNotFoundErr    = "could not find token in tokens"
	emptyFilenameErr    = "can not parse empty filename"
)

func traceError(s string) error {
	_, file, line, _ := runtime.Caller(2)
	prefix := filepath.Join(build.Default.GOPATH, "src") + string(os.PathSeparator)
	file = strings.TrimPrefix(file, prefix)
	return fmt.Errorf("%s:%d %s", file, line, s)
}
