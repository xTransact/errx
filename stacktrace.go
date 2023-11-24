package errx

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type stacktrace struct {
	file     string
	function string
	line     int
}

func newStacktrace() *stacktrace {
	// Caller of newStacktrace is newError, New or Wrap, so user's code is 3 up.
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return nil
	}
	file = removeGoPath(file)

	st := &stacktrace{
		file: file,
		line: line,
	}

	f := runtime.FuncForPC(pc)
	if f == nil {
		return nil
	}
	st.function = shortFuncName(f)

	return st
}

func (s *stacktrace) String() string {
	if s.function != "" {
		return fmt.Sprintf("    --- at %v:%v %v()", s.file, s.line, s.function)
	} else {
		return fmt.Sprintf("    --- at %v:%v", s.file, s.line)
	}
}

func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/palantir/shield/package.FuncName"
	// - "github.com/palantir/shield/package.Receiver.MethodName"
	// - "github.com/palantir/shield/package.(*PtrReceiver).MethodName"
	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}

/*
RemoveGoPath makes a path relative to one of the src directories in the $GOPATH
environment variable. If $GOPATH is empty or the input path is not contained
within any of the src directories in $GOPATH, the original path is returned. If
the input path is contained within multiple of the src directories in $GOPATH,
it is made relative to the longest one of them.
*/
func removeGoPath(path string) string {
	dirs := filepath.SplitList(os.Getenv("GOPATH"))
	// Sort in decreasing order by length so the longest matching prefix is removed
	sort.Stable(longestFirst(dirs))
	for _, dir := range dirs {
		srcdir := filepath.Join(dir, "src")
		rel, err := filepath.Rel(srcdir, path)
		// filepath.Rel can traverse parent directories, don't want those
		if err == nil && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return rel
		}
	}
	return path
}

type longestFirst []string

func (strs longestFirst) Len() int           { return len(strs) }
func (strs longestFirst) Less(i, j int) bool { return len(strs[i]) > len(strs[j]) }
func (strs longestFirst) Swap(i, j int)      { strs[i], strs[j] = strs[j], strs[i] }
