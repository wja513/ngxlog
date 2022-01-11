# Usage:
```go
package main

import (
	"fmt"
	"github.com/wja513/ngxlog"
	"os"
)

func main() {
	format := `$remote_addr $http_host $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $request_length $request_time $upstream_response_time "$http_x_forwarded_for" $http_x_real_ip`
	f, err := os.Open("example/demo.example.com-access.log")
	if err != nil {
		panic("demo.example.com-access.log does not exist")
	}
	defer f.Close()
	s := ngxlog.NewScanner(format, f)
	for s.Scan() {
		rec := s.Record()
		if rec.Mismatch() {
			ngxlog.LogMismatch(rec) //deal mismatch line
			continue
		}
		// do your business logic
		fmt.Println(rec.FieldTime("time_local"))
		//fmt.Println(s.Col(3))
	}
}
```

# Benchmark:
ngxlog:
```text
goos: darwin
goarch: amd64
pkg: github.com/wja513/ngxlog
cpu: Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz
BenchmarkParseLine
BenchmarkParseLine-12    	 2944748	       408.6 ns/op
```
nxgo:
```text
goos: darwin
goarch: amd64
pkg: daily_practice/nginxlog/gonx
cpu: Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz
BenchmarkParseLogRecord2
BenchmarkParseLogRecord2-12    	  169730	      6799 ns/op
```

# Reference
- https://github.com/satyrius/gonx
- https://pkg.go.dev/bufio#NewScanner