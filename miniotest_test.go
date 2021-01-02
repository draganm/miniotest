package miniotest_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/draganm/miniotest"
	mclient "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {

	mc, err := mclient.New(addr, &mclient.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	require.NoError(t, err)

	data := []byte("test")

	_, err = mc.PutObject(context.Background(), "test", "foo/var", bytes.NewReader(data), int64(len(data)), mclient.PutObjectOptions{})
	require.NoError(t, err)
}

func Test2(t *testing.T) {

	mc, err := mclient.New(addr, &mclient.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	require.NoError(t, err)

	data := []byte("test")

	_, err = mc.PutObject(context.Background(), "test", "foo/var", bytes.NewReader(data), int64(len(data)), mclient.PutObjectOptions{})
	require.NoError(t, err)
}

var addr string

func TestMain(m *testing.M) {
	var cleanup func() error
	var err error
	addr, cleanup, err = miniotest.StartEmbedded()

	if err != nil {
		fmt.Fprintf(os.Stderr, "while starting embedded server: %s", err)
		os.Exit(1)
	}

	exitCode := m.Run()

	err = cleanup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "while stopping embedded server: %s", err)
	}

	os.Exit(exitCode)
}
