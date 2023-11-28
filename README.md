# errx

`errx` is a more modern and convenient error handling library, providing more user-friendly stacktrace:


## Quick Start

```go
err := errx.New("oops: something wrong :(")
err = errx.Wrap(err, "failed to do something")
fmt.Printf("%+v\n", err)
```

## v3

> v3 supports `Error Code`

### Install

```shell
go get github.com/xTransact/errx/v3
```

### Example

```go
err := errx.NewInternalServerError().New("failed to do something")
err = Wrap(err, "oops: :(")
err = Wrapf(err, "xxxxxxx")
fmt.Printf("%+v\n", err)
```

Result:
```text
500: Internal Server Error: xxxxxxx: oops: :(: failed to do something
  Thrown: failed to do something
    --- at /home/xbank/go/pkg/errx/errx_test.go:17 TestCode1()
  Thrown: oops: :(
    --- at /home/xbank/go/pkg/errx/errx_test.go:18 TestCode1()
  Thrown: xxxxxxx
    --- at /home/xbank/go/pkg/errx/errx_test.go:19 TestCode1()
```


## v2

> v2 provides cleaner stacktrace.

### Install

```shell
go get github.com/xTransact/errx/v2
```

### Example

```go
func a() error {
	return errx.Wrap(b(), "a()")
}

func b() error {
	return c()
}

func c() error {
	return d()
}

func d() error {
	return e()
}

func e() error {
	return f()
}

func f() error {
	return errx.Wrap(g(), "f()")
}

func g() error {
	return fmt.Errorf("nil pointer dereference")
}

func main() {
    err := a()
    err = Wrapf(err, "internal server error: %s", "xBank")
    fmt.Printf("%+v\n", err)
}
```

Result:
```text
internal server error: xBank: a(): f(): nil pointer dereference
  Thrown: f()
    --- at /home/xbank/go/pkg/errx/errx_test.go:47 f()
  Thrown: a()
    --- at /home/xbank/go/pkg/errx/errx_test.go:27 a()
  Thrown: internal server error: xBank
    --- at /home/xbank/go/pkg/errx/errx_test.go:56 main()
```


## v1


## Install

```shell
go get github.com/xTransact/errx
```

## Example

```text
TestErrxWrap Failed: oops user: xBank: failed to check: xxx: something wrong.
Stacktrace:
  Error: something wrong.
    --- at /home/xbank/go/pkg/errx/errx_test.go:29 f()
    --- at /home/xbank/go/pkg/errx/errx_test.go:25 e()
    --- at /home/xbank/go/pkg/errx/errx_test.go:21 d()
    --- at /home/xbank/go/pkg/errx/errx_test.go:17 c()
    --- at /home/xbank/go/pkg/errx/errx_test.go:13 b()
  Thrown: xxx
    --- at /home/xbank/go/pkg/errx/errx_test.go:9 a()
    --- at /home/xbank/go/pkg/errx/errx_test.go:33 TestErrxWrap()
  Thrown: failed to check
    --- at /home/xbank/go/pkg/errx/errx_test.go:34 TestErrxWrap()
  Thrown: oops user: xBank
    --- at /home/xbank/go/pkg/errx/errx_test.go:35 TestErrxWrap()
```
