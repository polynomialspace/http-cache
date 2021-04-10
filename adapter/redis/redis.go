/*
MIT License

Copyright (c) 2018 Victor Springer

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package redis

import (
	"context"
	"time"

	redisCache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	cache "github.com/polynomialspace/http-cache"
	"github.com/vmihailenco/msgpack"
)

// Adapter is the memory adapter data structure.
type Adapter struct {
	store *redisCache.Cache
}

type RedisOptions = redis.Options

// Get implements the cache Adapter interface Get method.
func (a *Adapter) Get(key uint64) ([]byte, bool) {
	var c []byte
	if err := a.store.Get(context.Background(), cache.KeyAsString(key), &c); err == nil {
		return c, true
	}

	return nil, false
}

// Set implements the cache Adapter interface Set method.
func (a *Adapter) Set(key uint64, response []byte, expiration time.Time) {
	a.store.Set(&redisCache.Item{
		Key:   cache.KeyAsString(key),
		Value: response,
		TTL:   expiration.Sub(time.Now()),
	})
}

// Release implements the cache Adapter interface Release method.
func (a *Adapter) Release(key uint64) {
	a.store.Delete(context.Background(), cache.KeyAsString(key))
}

// NewAdapter initializes Redis adapter.
func NewAdapter(opt *RedisOptions) cache.Adapter {
	return &Adapter{
		redisCache.New(&redisCache.Options{
			Redis: redis.NewClient(opt),
			Marshal: func(v interface{}) ([]byte, error) {
				return msgpack.Marshal(v)

			},
			Unmarshal: func(b []byte, v interface{}) error {
				return msgpack.Unmarshal(b, v)
			},
		}),
	}
}
