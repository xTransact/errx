# errx

`errx` is a more modern and convenient error handling library, providing more user-friendly stack information:

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

## Install

```shell
go get github.com/xTransact/errx
```

## Quick Start

```go
err := errx.New("oops: something wrong :(")
err = errx.Wrap(err, "failed to do something")
fmt.Printf("%+v\n", err)
```
