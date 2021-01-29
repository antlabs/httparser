# httparser
[![Go](https://github.com/antlabs/httparser/workflows/Go/badge.svg)](https://github.com/antlabs/httparser/actions)
[![codecov](https://codecov.io/gh/antlabs/httparser/branch/master/graph/badge.svg)](https://codecov.io/gh/antlabs/httparser)  

高性能http 1.1解析器，为你的异步io库插上解析的翅膀[从零实现]

## 出发点
本来想基于异步io库写些好玩的代码，发现没有适用于这些库的http解析库，索性就自己写个，弥补golang生态一小片空白领域。

## 特性
* url解析
* request or response header field解析
* request or response  header value解析
* Content-Length数据包解析
* chunked数据包解析

## parser request
```go
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
```

## response
[response](./_example/response.go)

## request or response
如果你不确定数据包是请求还是响应，可看下面的例子  
[request or response](./_example/request_or_response.go)


## 编译
### 生成 unhex表和tokens表
如果需要修改这两个表，可以到_cmd目录下面修改生成代码的代码
```Makefile
make gen
```

### 编译example
```Makefile
make example
```
### 运行示例
```Makefile
make example.run
```

### 吞吐量
* 测试仓库 https://github.com/junelabs/httparser-benchmark
* Benchmark result: 8192.00 mb | 315.08 mb/s | 637803.27 req/sec | 26.00 s
