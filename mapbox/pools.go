package mapbox

import (
	"strings"
	"sync"
)

type stringsBufferPool struct {
	p sync.Pool
}

func newStringsBufferPool() *stringsBufferPool {
	return &stringsBufferPool{p: sync.Pool{New: func() interface{} {
		return &strings.Builder{}
	}}}
}

func (pool *stringsBufferPool) acquireStringsBuilder() *strings.Builder {
	return pool.p.Get().(*strings.Builder)
}

func (pool *stringsBufferPool) releaseStringsBuilder(b *strings.Builder) {
	b.Reset()
	pool.p.Put(b)
}
