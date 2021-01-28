package main

import (
	"fmt"
	"github.com/antlabs/httparser"
)

func main() {
	var data = []byte(
		"POST /joyent/http-parser HTTP/1.1\r\n" +
			"Host: github.com\r\n" +
			"DNT: 1\r\n" +
			"Accept-Encoding: gzip, deflate, sdch\r\n" +
			"Accept-Language: ru-RU,ru;q=0.8,en-US;q=0.6,en;q=0.4\r\n" +
			"User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) " +
			"AppleWebKit/537.36 (KHTML, like Gecko) " +
			"Chrome/39.0.2171.65 Safari/537.36\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9," +
			"image/webp,*/*;q=0.8\r\n" +
			"Referer: https://github.com/joyent/http-parser\r\n" +
			"Connection: keep-alive\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"Cache-Control: max-age=0\r\n\r\nb\r\nhello world\r\n0\r\n")

	var setting = httparser.Setting{
		MessageBegin: func() {
			//解析器开始工作
			fmt.Printf("begin\n")
		},
		URL: func(buf []byte) {
			//url数据
			fmt.Printf("url->%s\n", buf)
		},
		Status: func([]byte) {
			// 响应包才需要用到
		},
		HeaderField: func(buf []byte) {
			// http header field
			fmt.Printf("header field:%s\n", buf)
		},
		HeaderValue: func(buf []byte) {
			// http header value
			fmt.Printf("header value:%s\n", buf)
		},
		HeadersComplete: func() {
			// http header解析结束
			fmt.Printf("header complete\n")
		},
		Body: func(buf []byte) {
			fmt.Printf("%s", buf)
			// Content-Length 或者chunked数据包
		},
		MessageComplete: func() {
			// 消息解析结束
			fmt.Printf("\n")
		},
	}

	p := httparser.New(httparser.REQUEST)
	success, err := p.Execute(&setting, data)

	fmt.Printf("success:%d, err:%v\n", success, err)
}
