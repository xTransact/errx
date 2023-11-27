package errcode

import (
	"fmt"
	"net/http"
)

type Code int

const (
	DefaultCode Code = 0

	BadRequest                 Code = 400
	Unauthorized               Code = 401
	Forbidden                  Code = 403
	NotFound                   Code = 404
	MethodNotAllowed           Code = 405
	NotAcceptable              Code = 406
	RequestTimeout             Code = 408
	Conflict                   Code = 409
	LengthRequired             Code = 411
	RequestEntityTooLarge      Code = 413
	RequestURITooLong          Code = 414
	UnsupportedMediaType       Code = 415
	RequestRangeNotSatisfied   Code = 416
	MisdirectedRequest         Code = 421
	TooManyRequests            Code = 429
	UnavailableForLegalReasons Code = 451

	InternalServerError           Code = 500
	NotImplemented                Code = 501
	BadGateway                    Code = 502
	ServiceUnavailable            Code = 503
	GatewayTimeout                Code = 504
	InsufficientStorage           Code = 507
	NetworkAuthenticationRequired Code = 511
)

func New(code int) Code {
	return Code(code)
}

func IsDefaultCode(code Code) bool {
	return code == DefaultCode
}

func IsDefaultCodeValue(code int) bool {
	return code == DefaultCode.Int()
}

func IsNotDefaultCode(code Code) bool {
	return !IsDefaultCode(code)
}

func (c Code) Int() int {
	return int(c)
}

func (c Code) String() string {
	code := c.Int()
	str := http.StatusText(code)
	if str == "" {
		str = "Unknown Error"
	}
	return fmt.Sprintf("%d: %s", code, str)
}
