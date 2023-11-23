package errx

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

// /
// / Inspired by palantir/stacktrace repo
// / -> https://github.com/palantir/stacktrace/blob/master/stacktrace.go
// / -> Apache 2.0 LICENSE
// /

type fake struct{}

var (
	StackTraceMaxDepth  int = 15
	packageName             = reflect.TypeOf(fake{}).PkgPath()
	packageNameExamples     = packageName + "/examples/"
)

type stacktraceFrame struct {
	pc       uintptr
	file     string
	function string
	line     int
}

func (frame *stacktraceFrame) String() string {
	currentFrame := fmt.Sprintf("%v:%v", frame.file, frame.line)
	if frame.function != "" {
		currentFrame = fmt.Sprintf("%v:%v %v()", frame.file, frame.line, frame.function)
	}

	return currentFrame
}

type stacktrace struct {
	span   string
	frames []stacktraceFrame
}

func (st *stacktrace) Error() string {
	return st.String("")
}

func (st *stacktrace) String(deepestFrame string) string {
	var str string

	newline := func() {
		if str != "" && !strings.HasSuffix(str, "\n") {
			str += "\n"
		}
	}

	for _, frame := range st.frames {
		if frame.file != "" {
			currentFrame := frame.String()
			if currentFrame == deepestFrame {
				break
			}

			newline()
			str += "  --- at " + currentFrame
		}
	}

	return str
}

func (st *stacktrace) Source() (string, []string) {
	if len(st.frames) == 0 {
		return "", []string{}
	}

	firstFrame := st.frames[0]

	header := firstFrame.String()
	body := getSourceFromFrame(firstFrame)

	return header, body
}

func newStacktrace(span string) *stacktrace {
	frames := []stacktraceFrame{}

	// We loop until we have StackTraceMaxDepth frames or we run out of frames.
	// Frames from this package are skipped.
	for i := 0; len(frames) < StackTraceMaxDepth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		file = removeGoPath(file)

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		function := shortFuncName(f)

		isGoPkg := len(runtime.GOROOT()) > 0 && strings.Contains(file, runtime.GOROOT()) // skip frames in GOROOT if it's set
		isErrxPkg := isErrxPkg(file)                                                     // skip frames in this package
		isExamplePkg := strings.Contains(file, packageNameExamples)                      // do not skip frames in this package examples
		isTestPkg := strings.Contains(file, "_test.go")                                  // do not skip frames in tests

		if !isGoPkg && (!isErrxPkg || isExamplePkg || isTestPkg) {
			frames = append(frames, stacktraceFrame{
				pc:       pc,
				file:     file,
				function: function,
				line:     line,
			})
		}
	}

	return &stacktrace{
		span:   span,
		frames: frames,
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

func isErrxPkg(file string) bool {
	return strings.Contains(file, packageName) || strings.Contains(file, "github.com/x!transact/errx")
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
