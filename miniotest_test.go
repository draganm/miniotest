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
