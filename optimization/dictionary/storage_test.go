package dictionary

import "testing"

//BenchmarkStorage_Get/map-16         	31081088	        38.0 ns/op	       0 B/op	       0 allocs/op
//BenchmarkStorage_Get/array-16       	189884908	         6.12 ns/op	       0 B/op	       0 allocs/op
func BenchmarkStorage_Get(b *testing.B) {

	maxId := 1000

	testCases := []struct {
		name string
		ds   Storage
	}{
		{
			"map",
			NewMapStorage(uint16(maxId)),
		},
		{
			"array",
			NewArrayStorage(uint16(maxId)),
		},
	}

	for _, tt := range testCases {
		b.Run(tt.name, func(b *testing.B) {
			for i := 1; i <= maxId; i++ {
				err := tt.ds.Add(Entity{
					Id:     uint16(i),
					Field1: 0,
					Field2: "",
				})
				if err != nil {
					b.Fatal(err)
				}
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				e, err := tt.ds.Get(uint16(i%maxId) + 1)
				if err != nil {
					b.Fatal(err)
				}
				_ = e
			}
		})
	}
}
