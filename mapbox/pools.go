package mapbox

import (
	"bytes"
	"sync"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type stringsBufferPool struct {
	noCopy noCopy
	p      sync.Pool
}

func newStringsBufferPool() *stringsBufferPool {
	return &stringsBufferPool{p: sync.Pool{New: func() interface{} {
		return &bytes.Buffer{}
	}}}
}

func (pool *stringsBufferPool) acquireStringsBuilder() *bytes.Buffer {
	return pool.p.Get().(*bytes.Buffer)
}

func (pool *stringsBufferPool) releaseStringsBuilder(b *bytes.Buffer) {
	b.Reset()
	pool.p.Put(b)
}
