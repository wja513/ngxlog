package ngxlog

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestNewScanner(t *testing.T) {
	// 4 lines
	logLines := bytes.NewBufferString(`127.0.0.1 demo.example.com - [22/Nov/2021:09:14:08 +0800] "GET /open/serviceCode?siteId=121915&tid= HTTP/1.1" 200 73 "-" "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36" 1433 0.032 0.032 "101.80.148.128, 101.80.148.128" 101.80.148.128

demo.example.com - [22/Nov/2021:09:14:08 +0800] "GET /open/serviceCode?siteId=121915&tid= HTTP/1.1" 200 73 "-" "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36" 1433 0.032 0.032 "101.80.148.128, 101.80.148.128" 101.80.148.128
`)
	largeLine := `127.0.0.1 demo.example.com - [22/Nov/2021:09:14:08 +0800] "GET /open/serviceCode?siteId=121915&tid=%s HTTP/1.1" 200 73 "-" "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36" 1433 0.032 0.032 "101.80.148.128, 101.80.148.128" 101.80.148.128`
	logLines.WriteString(fmt.Sprintf(largeLine, strings.Repeat("x", 1024*64))) // default max line size is bufio.MaxScanTokenSize = 64kb

	format := `$remote_addr $http_host $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $request_length $request_time $upstream_response_time "$http_x_forwarded_for" $http_x_real_ip`
	s := NewScanner(format, logLines)

	// normal line
	s.Scan()
	if rec := s.Record(); rec.Field("time_local") != "22/Nov/2021:09:14:08 +0800" {
		t.Error("normal line test failed")
	}

	// empty line
	s.Scan()
	if rec := s.Record(); !rec.Mismatch() {
		t.Error("empty line test failed")
	}

	// field mismatch line
	s.Scan()
	if rec := s.Record(); !rec.Mismatch() {
		t.Error("empty line test failed")
	}

	// super large line exceed 64kb will be ignored
	s.Scan()
	if rec := s.Record(); !rec.Mismatch() {
		t.Error("super large line test failed")
	}
}

func BenchmarkParseLine(b *testing.B) {
	format := s2b(`$remote_addr $http_host $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" $request_length $request_time $upstream_response_time "$http_x_forwarded_for" $http_x_real_ip`)
	line := s2b(`127.0.0.1 demo.example.com - [22/Nov/2021:09:14:08 +0800] "GET /open/serviceCode?siteId=121915&tid= HTTP/1.1" 200 73 "-" "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36" 1433 0.032 0.032 "101.80.148.128, 101.80.148.128" 101.80.148.128`)
	n := len(parseLine(format, 0))
	for i := 0; i < b.N; i++ {
		parseLine(line, n)
	}
}
