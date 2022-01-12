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
	//s.Buffer(make([]byte, 8*1024*1024), 8*1024*1024) // max line size 8MB,default 64KB
	for s.Scan() {
		rec := s.Record()
		// do your business logic
		//fmt.Println(rec.FieldTime("time_local"))
		//fmt.Println(s.Col(3))
		fmt.Println(rec.Field("http_x_real_ip"))
	}
}
