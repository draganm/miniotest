# Miniotest

Convenience Golang module enabling you to run embedded [Minio](https://min.io/) server for purpose of integration testing of AWS S3 operations.

## Motivation

There is an [open feature request](https://github.com/minio/minio/issues/5146) on the [Minio GitHub project](https://github.com/minio/minio).
It describes the necessary steps to use Minio in your tests, but still there is no easily (re-)useable code for that purpose.

Notably, following features would we very useful:
- Do no hard-code the port to which Minio server will bind, instead use a free port determined at run time
- Create a bucket after starting the server
- Provide a tear-down function that will shut down the Minio server after the test is done
- Clean up files created by the test after the test is done

All of those features are provided by this module.

## Requirements

- Golang 1.15: uses Golang's test [TempDir](https://golang.org/pkg/testing/#T.TempDir) feature

## Installation

```bash
$ go get github.com/draganm/miniotest

```

## Use
[this project's test](./miniotest_test.go) demonstrates the use of the embedded Minio for testing.

```golang
package miniotest_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/draganm/miniotest"
	mclient "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	addr, cleanup := miniotest.StartEmbedded(t)
	defer cleanup()

	mc, err := mclient.New(addr, &mclient.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	require.NoError(t, err)

	data := []byte("test")

	_, err = mc.PutObject(context.Background(), "test", "foo/var", bytes.NewReader(data), int64(len(data)), mclient.PutObjectOptions{})
	require.NoError(t, err)
}

```

Function `miniotest.StartEmbedded()` starts Minio server and returns a string representing bound `address:port` where Minio is listening.
You can use this address as the endpoint for your S3 client (e.g. minio client).

Second value returned is a function that will shut down the started Minio server when called.
The best use of this function is to be deferred until the test function returns (as shown in the example).

Passed `t *testing.T` parameter is used to create a test scoped temporary directory that will be deleted after the test has been run.

## Default values

### Hostname and port

Miniotest will force embedded Minio server to bind to `localhost` and a free port.
Free port is determined by opening a TCP listener with the port `0`, getting the bound port and closing the listener.
`hostname:port` bind address will be returned as the first return value of `miniotest.StartEmbedded()`

Note: there is a possible race condition with this approach, which can lead to Minio server not starting if some other process binds the same port between closing the listener and starting the Minio server.
Window of opportunity for this to happen is very short and one can assume this won't happen in 99.999% of cases.

### `accessKeyID` and `secretAccessKey`

When not specified otherwise, Minio will assume that both `accessKeyID` and `secretAccessKey` are set to `minioadmin` and you should be using this for your tests.
Please do not set environment variables `MINIO_ACCESS_KEY` and `MINIO_SECRET_KEY` because this will break the setup code for the embedded Minio server.


### Automatically created bucket

Miniotest will create bucket named `test` as a convenience for the tests run.

### Secure mode

Embedded Minio server won't be listening for TLS connections, hence one should not use the S3 client with this option enabled.
If someone really needs this feature, please raise an issue on this project, or even better offer a pull request for it.

## Contributing

Contributions are welcome, send your issues and PRs to this repo.

## License

[MIT](LICENSE) - Copyright Dragan Milic and contributors
