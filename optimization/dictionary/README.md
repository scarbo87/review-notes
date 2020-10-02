### Dictionary

If you need to keep a dictionary (with uint keys) in memory, it'll be better to do it in an array then a map.

```
//BenchmarkStorage_Get/map-16         	31081088	        38.0 ns/op	       0 B/op	       0 allocs/op
//BenchmarkStorage_Get/array-16       	189884908	         6.12 ns/op	       0 B/op	       0 allocs/op
```