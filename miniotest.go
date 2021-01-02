package miniotest

import (
	"context"
	"io/ioutil"
	"net"
	"os"

	// "github.com/minio/minio-go/pkg/credentials"
	mclient "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/pkg/madmin"
	"github.com/pkg/errors"
)

func StartEmbedded() (string, func() error, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", nil, errors.Wrap(err, "while creating listener")
	}

	addr := l.Addr().String()
	err = l.Close()
	if err != nil {
		return "", nil, errors.Wrap(err, "while closing listener")
	}

	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"

	madm, err := madmin.New(addr, accessKeyID, secretAccessKey, false)
	if err != nil {
		return "", nil, errors.Wrap(err, "while creating madimin")
	}

	td, err := ioutil.TempDir("", "")
	if err != nil {
		return "", nil, errors.Wrap(err, "while creating temp dir")
	}

	go minio.Main([]string{"minio", "server", "--quiet", "--address", addr, td})

	mc, err := mclient.New(addr, &mclient.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})

	err = mc.MakeBucket(context.Background(), "test", mclient.MakeBucketOptions{})
	if err != nil {
		return "", nil, errors.Wrap(err, "while creating bucket")
	}

	return addr, func() error {
		err := madm.ServiceStop(context.Background())
		if err != nil {
			return errors.Wrap(err, "while stopping service")
		}

		err = os.RemoveAll(td)
		if err != nil {
			return errors.Wrap(err, "while deleting temp dir")
		}

		return nil
	}, nil

}
