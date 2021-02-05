package main

import (
	"fmt"
	"github.com/antlabs/httparser"
)

func main() {
	var data = []byte(
		"HTTP/1.1 200 OK\r\n" +
			"Content-Type: text/plain\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"7\r\n" +
			"Mozilla\r\n" +
			"9\r\n" +
			"Developer\r\n" +
			"7\r\n" +
			"Network\r\n" +
			"0\r\n" +
			"\r\n")

	var setting = httparser.Setting{
		MessageBegin: func(*httparser.Parser) {
			//解析器开始工作
			fmt.Printf("begin\n")
		},
		URL: func(_ *httparser.Parser, buf []byte) {
			//url数据
			fmt.Printf("url->%s\n", buf)
		},
		Status: func(*httparser.Parser, []byte) {
			// 响应包才需要用到
		},
		HeaderField: func(_ *httparser.Parser, buf []byte) {
			// http header field
			fmt.Printf("header field:%s\n", buf)
		},
		HeaderValue: func(_ *httparser.Parser, buf []byte) {
			// http header value
			fmt.Printf("header value:%s\n", buf)
		},
		HeadersComplete: func(_ *httparser.Parser) {
			// http header解析结束
			fmt.Printf("header complete\n")
		},
		Body: func(_ *httparser.Parser, buf []byte) {
			fmt.Printf("%s", buf)
			// Content-Length 或者chunked数据包
		},
		MessageComplete: func(_ *httparser.Parser) {
			// 消息解析结束
			fmt.Printf("\n")
		},
	}

	p := httparser.New(httparser.RESPONSE)
	success, err := p.Execute(&setting, data)

	fmt.Printf("success:%d, err:%v\n", success, err)
}
