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
	"github.com/minio/minio-go/v7"
	"sync"
)

const (
	DefaultChunkSize = 1024 * 1024 * 4 // 4MB
)

type Download struct {
	client     *minio.Client
	ctx        context.Context
	chunkSize  int64
	bucket     string
	key        string
	inflight   map[int64]*Chunk
	inflightMu sync.RWMutex
}

func NewDownload(client *minio.Client, ctx context.Context, chunkSize int64, bucket string, key string) *Download {
	if chunkSize == 0 {
		chunkSize = DefaultChunkSize
	}
	return &Download{
		client:    client,
		ctx:       ctx,
		chunkSize: chunkSize,
		bucket:    bucket,
		key:       key,
		inflight:  make(map[int64]*Chunk),
	}
}

// Get downloads a new chunk of data from the remote server
// using a range request. The chunk size is determined by the
// chunkSize field of the Download object. If the chunk size is 0,
// then the default chunk size is DefaultChunkSize (4MB).
func (d *Download) Get(offset int64) ([]byte, error) {
	d.inflightMu.RLock()
	c, ok := d.inflight[offset]
	d.inflightMu.RUnlock()
	if !ok {
		d.inflightMu.Lock()
		c, ok = d.inflight[offset]
		if !ok {
			var err error
			c, err = NewChunk(d.client, d.ctx, d.chunkSize, offset, d.bucket, d.key)
			if err != nil {
				d.inflightMu.Unlock()
				return nil, err
			}
			d.inflight[offset] = c
		}
		d.inflightMu.Unlock()
	}
	return c.Wait()
}
