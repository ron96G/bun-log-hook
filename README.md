# Bun Log Hook
Log Hook for [bun](https://github.com/uptrace/bun).

## Benchmark results

```
$ go test -run=. -bench=. -benchtime=5s -count 5 -benchmem
goos: windows
goarch: amd64
pkg: github.com/ron96G/bun-log-hook
cpu: AMD Ryzen 9 5900X 12-Core Processor
```
Benchmark Name|Iterations|ns/op|B/op|allocs/op
----|----|----|----|----
BenchmarkLoghook-24 | 14692111 | 387.8 ns/op | 269 B/op | 1 allocs/op
BenchmarkLoghook-24 | 18036872 | 366.9 ns/op | 240 B/op | 1 allocs/op
BenchmarkLoghook-24 | 18299829 | 335.7 ns/op | 47 B/op | 1 allocs/op
BenchmarkLoghook-24 | 18169380 | 335.2 ns/op | 47 B/op | 1 allocs/op
BenchmarkLoghook-24 | 18313262 | 335.2 ns/op | 47 B/op | 1 allocs/op
BenchmarkLoghookFailed-24 | 17416347 | 352.5 ns/op | 47 B/op | 1 allocs/op
BenchmarkLoghookFailed-24 | 17238441 | 347.8 ns/op | 47 B/op | 1 allocs/op
BenchmarkLoghookFailed-24 | 17280080 | 348.6 ns/op | 47 B/op | 1 allocs/op
BenchmarkLoghookFailed-24 | 17256240 | 347.7 ns/op | 47 B/op | 1 allocs/op
BenchmarkLoghookFailed-24 | 17189740 | 356.4 ns/op | 47 B/op | 2 allocs/op