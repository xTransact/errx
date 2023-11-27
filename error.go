package errx

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/xTransact/errx/v3/errcode"
)

var _ error = (*xerr)(nil)

type xerr struct {
	err        error
	code       errcode.Code
	msg        string
	stacktrace *stacktrace
}

func newError(err error, code errcode.Code, format string, args ...any) *xerr {
	return &xerr{
		err:        err,
		code:       code,
		msg:        fmt.Sprintf(format, args...),
		stacktrace: newStacktrace(),
	}
}

// Error returns the error message, without context.
func (e *xerr) Error() string {
	if e.err != nil {
		if e.msg == "" {
			return e.err.Error()
		}

		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}

	return e.msg
}

func (e *xerr) Code() errcode.Code {
	return e.code
}

func (e *xerr) Cause() error {
	return e.err
}

func (e *xerr) Unwrap() error {
	return e.err
}

// Format implements fmt.Formatter.
// If the format is "%+v", then the details of the error are included.
// Otherwise, using "%v", just the summary is included.
func (e *xerr) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		fmt.Fprint(s, e.formatVerbose())
	} else {
		fmt.Fprint(s, e.formatSummary())
	}
}

func (e *xerr) formatVerbose() string {
	return e.formatStacktrace()
}

func (e *xerr) formatSummary() string {
	return e.Error()
}

func (e *xerr) formatStacktrace() string {
	if e.stacktrace == nil {
		return ""
	}

	rows := make([]string, 0)

	recursive(e, func(err *xerr) {
		var row string
		newline := func() {
			if row != "" && !strings.HasSuffix(row, "\n") {
				row += "\n"
			}
		}

		if err == nil {
			return
		}

		if err.msg != "" {
			row += "  Thrown: " + err.msg
		}

		if err.stacktrace != nil {
			if st := err.stacktrace.String(); st != "" {
				newline()
				row += err.stacktrace.String()
			}
		}

		if strings.TrimSpace(row) != "" {
			rows = append(rows, row)
		}
	})

	slices.Reverse(rows)

	first := e.Error()
	if errcode.IsNotDefaultCode(e.code) {
		first = e.code.String() + ": " + first
	}

	rows = slices.Insert(rows, 0, first)

	return strings.Join(rows, "\n")
}

func recursive(e *xerr, tap func(*xerr)) {
	tap(e)

	if e.err == nil {
		return
	}

	var err *xerr
	if ok := errors.As(e.err, &err); ok {
		recursive(err, tap)
	}
}
