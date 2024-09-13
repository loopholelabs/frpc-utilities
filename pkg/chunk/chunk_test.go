// SPDX-License-Identifier: Apache-2.0

package chunk

import (
	"bytes"
	"context"
	"crypto/rand"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T) (*minio.Client, []byte, string, string) {
	client, err := minio.New("play.min.io", &minio.Options{
		Creds:  credentials.NewStaticV4("Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", ""),
		Secure: true,
	})
	require.NoError(t, err)

	bucketName := uuid.New().String()
	err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
		Region: "us-east-1",
	})
	require.NoError(t, err)

	objectName := uuid.New().String()
	objectContent := make([]byte, 1024*1024)
	_, err = rand.Read(objectContent)
	require.NoError(t, err)

	_, err = client.PutObject(context.Background(), bucketName, objectName, bytes.NewReader(objectContent), int64(len(objectContent)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		err = client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
		require.NoError(t, err)
		err = client.RemoveBucket(context.Background(), bucketName)
		require.NoError(t, err)
	})

	return client, objectContent, bucketName, objectName
}

func TestChunk(t *testing.T) {
	client, data, bucket, obj := testClient(t)
	var offset int64 = 0
	const chunkSize = 512

	chunk, err := GetChunk(client, context.Background(), offset, chunkSize, bucket, obj)
	require.NoError(t, err)

	downloadedData, err := chunk.Wait()
	require.NoError(t, err)
	require.Equal(t, chunkSize, len(downloadedData))
	require.Equal(t, data[:chunkSize], downloadedData)

	chunk.Return()

	offset += chunkSize * 2
	chunk, err = GetChunk(client, context.Background(), offset, chunkSize, bucket, obj)
	require.NoError(t, err)

	downloadedData, err = chunk.Wait()
	require.NoError(t, err)
	require.Equal(t, chunkSize, len(downloadedData))
	require.Equal(t, data[offset:offset+chunkSize], downloadedData)

	chunk.Return()
}

func TestConcurrentChunk(t *testing.T) {
	client, data, bucket, obj := testClient(t)
	const offset = 32
	const chunkSize = 512

	chunk, err := GetChunk(client, context.Background(), offset, chunkSize, bucket, obj)
	require.NoError(t, err)

	start := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			<-start
			downloadedData, err := chunk.Wait()
			require.NoError(t, err)
			require.Equal(t, chunkSize, len(downloadedData))
			require.Equal(t, data[offset:offset+chunkSize], downloadedData)
			wg.Done()
		}()
	}

	close(start)
	wg.Wait()

	chunk.Return()
}

func TestInvalidChunkOffset(t *testing.T) {
	client, data, bucket, obj := testClient(t)
	var offset = len(data) + 1
	const chunkSize = 512

	chunk, err := GetChunk(client, context.Background(), int64(offset), chunkSize, bucket, obj)
	require.NoError(t, err)

	_, err = chunk.Wait()
	require.Error(t, err)
	chunk.Return()
}

func TestInvalidChunkSize(t *testing.T) {
	client, data, bucket, obj := testClient(t)
	const offset = 512
	var chunkSize = len(data) + 1

	chunk, err := GetChunk(client, context.Background(), offset, int64(chunkSize), bucket, obj)
	require.NoError(t, err)

	downloadedData, err := chunk.Wait()
	require.NoError(t, err)
	require.Equal(t, len(data)-offset, len(downloadedData))
	require.Equal(t, data[offset:], downloadedData)
	chunk.Return()
}
