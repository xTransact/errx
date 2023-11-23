package errx

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var SourceFragmentsHidden = true

var _ error = (*xerr)(nil)

type xerr struct {
	err        error
	msg        string
	stacktrace *stacktrace
}

// Unwrap returns the underlying error.
func (e xerr) Unwrap() error {
	return e.err
}

func (e xerr) Is(err error) bool {
	return errors.Is(err, e.err)
}

// Error returns the error message, without context.
func (e xerr) Error() string {
	if e.err != nil {
		if e.msg == "" {
			return e.err.Error()
		}

		return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
	}

	return e.msg
}

// Stacktrace returns a pretty printed stacktrace of the error.
func (e xerr) Stacktrace() string {
	blocks := []string{}
	topFrame := ""

	recursive(e, func(e2 xerr) {
		if e2.stacktrace != nil && len(e2.stacktrace.frames) > 0 {
			err := ternaryF(e2.err != nil, func() string { return e2.err.Error() }, func() string { return "" })
			block := coalesceOrEmpty(e2.msg, err, "Error")
			topFrameStacktrace := e2.stacktrace.String(topFrame)
			if strings.TrimSpace(topFrameStacktrace) != "" {
				block += "\n" + topFrameStacktrace
			}

			blocks = append([]string{block}, blocks...)

			topFrame = e2.stacktrace.frames[0].String()
		}
	})

	if len(blocks) == 0 {
		return ""
	}

	return "Error: " + strings.Join(blocks, "\nThrown: ")
}

// Sources returns the source fragments of the error.
func (e xerr) Sources() string {
	blocks := [][]string{}

	recursive(e, func(e xerr) {
		if e.stacktrace != nil && len(e.stacktrace.frames) > 0 {
			header, body := e.stacktrace.Source()

			if e.msg != "" {
				header = fmt.Sprintf("%s\n%s", e.msg, header)
			}

			if header != "" && len(body) > 0 {
				blocks = append(
					[][]string{append([]string{header}, body...)},
					blocks...,
				)
			}
		}
	})

	if len(blocks) == 0 {
		return ""
	}

	return strings.Join(
		manipulateMap(blocks, func(items []string, _ int) string {
			return strings.Join(items, "\n")
		}),
		"\n\nThrown: ",
	)
}

// ToMap returns a map representation of the error.
func (e xerr) ToMap() map[string]any {
	result := make(map[string]any)
	result["error"] = e.Error()
	result["stacktrace"] = e.Stacktrace()
	if !SourceFragmentsHidden {
		result["sources"] = e.Sources()
	}
	return result
}

// MarshalJSON implements json.Marshaler.
func (e xerr) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.ToMap())
}

// Format implements fmt.Formatter.
// If the format is "%+v", then the details of the error are included.
// Otherwise, using "%v", just the summary is included.
func (e xerr) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		fmt.Fprint(s, e.formatVerbose())
	} else {
		fmt.Fprint(s, e.formatSummary())
	}
}

func (e *xerr) formatVerbose() string {
	output := e.Error()

	if st := e.Stacktrace(); st != "" {
		lines := strings.Split(st, "\n")
		st = "  " + strings.Join(lines, "\n  ")
		output += fmt.Sprintf("\nStacktrace:\n%s\n", st)
	}

	if sources := e.Sources(); sources != "" && !SourceFragmentsHidden {
		output += fmt.Sprintf("\nSources:\n%s\n", sources)
	}

	return output
}

func (e *xerr) formatSummary() string {
	return e.Error()
}

func newError() xerr {
	return xerr{
		stacktrace: newStacktrace(""),
	}
}

func (e xerr) copy() xerr {
	return xerr{
		err:        e.err,
		msg:        e.msg,
		stacktrace: e.stacktrace,
	}
}

func (e xerr) wrapErr(err error) error {
	if err == nil {
		return nil
	}

	if e.err == nil {
		e.err = err
	} else {
		e.err = fmt.Errorf("%w: %w", err, e.err)
	}

	return e
}

func (e xerr) wrapMsg(msg string) string {
	if msg == "" {
		return ""
	}

	if e.msg == "" {
		return msg
	}

	return fmt.Sprintf("%s: %s", msg, e.msg)
}

func (e xerr) wrapMsgf(format string, args ...any) string {
	if format == "" {
		return ""
	}

	msg := fmt.Sprintf(format, args...)
	if e.msg == "" {
		return msg
	}

	return fmt.Sprintf("%s: %s", msg, e.msg)
}

func recursive(err xerr, tap func(xerr)) {
	tap(err)

	if err.err == nil {
		return
	}

	if child, ok := AsError(err.err); ok {
		recursive(child, tap)
	}
}
