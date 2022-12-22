/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package download

import (
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestDownload(t *testing.T) {
	client, data, bucket, obj := testClient(t)
	var offset int64 = 0
	const chunkSize = 512

	download := NewDownload(client, context.Background(), chunkSize, bucket, obj)

	downloadedData, err := download.Get(offset)
	require.NoError(t, err)
	require.Equal(t, chunkSize, len(downloadedData))
	require.Equal(t, data[:chunkSize], downloadedData)

	offset += chunkSize * 2
	downloadedData, err = download.Get(offset)
	require.NoError(t, err)
	require.Equal(t, chunkSize, len(downloadedData))
	require.Equal(t, data[offset:offset+chunkSize], downloadedData)
}

func TestConcurrentDownload(t *testing.T) {
	client, data, bucket, obj := testClient(t)
	const offset = 32
	const chunkSize = 512

	download := NewDownload(client, context.Background(), chunkSize, bucket, obj)

	start := make(chan struct{})
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			<-start
			downloadedData, err := download.Get(offset)
			require.NoError(t, err)
			require.Equal(t, chunkSize, len(downloadedData))
			require.Equal(t, data[offset:offset+chunkSize], downloadedData)
			wg.Done()
		}()
	}

	close(start)
	wg.Wait()
}
