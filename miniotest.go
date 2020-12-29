package miniotest

import (
	"context"
	"fmt"
	"net"
	"testing"

	// "github.com/minio/minio-go/pkg/credentials"
	mclient "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/pkg/madmin"
	"github.com/stretchr/testify/require"
)

func StartEmbedded(t *testing.T) (string, func()) {
	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	addr := l.Addr().String()
	err = l.Close()
	require.NoError(t, err)

	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"

	madm, err := madmin.New(addr, accessKeyID, secretAccessKey, false)
	require.NoError(t, err)

	go minio.Main([]string{"minio", "server", "--address", addr, t.TempDir()})

	mc, err := mclient.New(addr, &mclient.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	require.NoError(t, err)

	defer func() {
		if err != nil {
			err = madm.ServiceStop(context.Background())
			require.NoError(t, err)
			fmt.Println("stopped")
		}
	}()

	err = mc.MakeBucket(context.Background(), "test", mclient.MakeBucketOptions{})
	require.NoError(t, err)

	return addr, func() {
		err = madm.ServiceStop(context.Background())
		require.NoError(t, err)
	}
}
