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

func (b Builder) Build() error {
	err := xerr(b)
	return &err
}
