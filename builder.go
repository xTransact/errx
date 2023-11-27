package errx

import (
	"fmt"
)

type Builder xerr

func (b Builder) New(message string) error {
	b.msg = message
	return b.Build()
}

func (b Builder) Errorf(format string, args ...any) error {
	b.msg = fmt.Sprintf(format, args...)
	return b.Build()
}

func (b Builder) Wrap(err error, message string) error {
	b.err = err
	b.msg = message
	return b.Build()
}

func (b Builder) Wrapf(err error, format string, args ...any) error {
	b.err = err
	b.msg = fmt.Sprintf(format, args...)
	return b.Build()
}

func (b Builder) WithStack(err error) error {
	b.err = err
	return b.Build()
}

func (b Builder) Build() error {
	err := xerr(b)
	return &err
}
