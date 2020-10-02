package rewrited_slice

import (
	"sync"
	"testing"
)

var syncPoolEntity = sync.Pool{
	New: func() interface{} {
		return &Entity{}
	},
}

type Entity struct {
	Id     uint64
	Field1 int64
	Field2 int64
}

type rewritedSlice struct {
	capacity uint64
	limit    uint64
	storage  []*Entity
}

func newRewritedSlice(limit, capacity uint64) *rewritedSlice {
	return &rewritedSlice{
		limit:    limit,
		capacity: capacity,
		storage:  make([]*Entity, 0, capacity),
	}
}

func (s *rewritedSlice) shrink1() {
	s.storage = s.storage[s.capacity-s.limit:]
}

func (s *rewritedSlice) shrink2() {
	copy(s.storage, s.storage[s.capacity-s.limit:])
	s.storage = s.storage[:s.limit]
}

func (s *rewritedSlice) shrink3() {
	tmp := s.storage
	copy(s.storage, s.storage[s.capacity-s.limit:])
	s.storage = s.storage[:s.limit]

	for _, e := range tmp[s.limit:] {
		e.Id = 0
		e.Field1 = 0
		e.Field2 = 0
		syncPoolEntity.Put(e)
	}
}

var (
	rs       *rewritedSlice
	limit    uint64 = 60000
	capacity uint64 = 100000
)

//BenchmarkRewritedSlice_Shrink1-16       100000000              145 ns/op              71 B/op          1 allocs/op (shrink1 was called 2498 times)
func BenchmarkRewritedSlice_Shrink1(b *testing.B) {

	rs = newRewritedSlice(limit, capacity)

	var count int
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rs.storage = append(rs.storage, &Entity{Id: uint64(i)})
		if len(rs.storage) >= int(rs.capacity) {
			rs.shrink1()
			count++
		}
	}

	b.Log(count)
}

//BenchmarkRewritedSlice_Shrink2-16       169291350               77.5 ns/op            32 B/op          1 allocs/op (shrink2 was called 4230 times)
func BenchmarkRewritedSlice_Shrink2(b *testing.B) {

	rs = newRewritedSlice(limit, capacity)

	var count int
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rs.storage = append(rs.storage, &Entity{Id: uint64(i)})
		if len(rs.storage) >= int(rs.capacity) {
			rs.shrink2()
			count++
		}
	}

	b.Log(count)
}

//BenchmarkRewritedSlice_Shrink3-16       292590810               43.8 ns/op             0 B/op          0 allocs/op (shrink3 was called 7313 times)
func BenchmarkRewritedSlice_Shrink3(b *testing.B) {

	rs = newRewritedSlice(limit, capacity)

	var count int
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e := syncPoolEntity.Get().(*Entity)
		e.Id = uint64(i)
		rs.storage = append(rs.storage, e)

		if len(rs.storage) >= int(rs.capacity) {
			rs.shrink3()
			count++
		}
	}

	b.Log(count)
}
