package httparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Isuse1(t *testing.T) {
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
			"Cache-Control: max-age=0\r\n\r\nb\r\nhello world\r\n0\r\n\r\n" +

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
			"Cache-Control: max-age=0\r\n\r\nb\r\nhello world\r\n0\r\n\r\n")

	body := []byte{}
	var setting = Setting{
		MessageBegin: func(*Parser) {
			fmt.Println("---- begin")
		},
		URL: func(p *Parser, buf []byte) {
		},
		Status: func(*Parser, []byte) {
			// 响应包才需要用到
		},
		HeaderField: func(p *Parser, buf []byte) {
		},
		HeaderValue: func(p *Parser, buf []byte) {
		},
		HeadersComplete: func(p *Parser) {
		},
		Body: func(p *Parser, buf []byte) {
			body = append(body, buf...)
		},
		MessageComplete: func(p *Parser) {
			fmt.Printf("%#v\n", p)
			p.Reset()
		},
	}

	p := New(REQUEST)
	fmt.Printf("req_len=%d\n", len(data)/2)
	// 一个POST 518，一共两个POST，第一次解析600字节，第二次解析剩余的
	data1, data2 := data[:600], data[600:]
	n, err := p.Execute(&setting, data1)
	if err != nil {
		panic(err.Error())
	}

	_, err = p.Execute(&setting, append(data1[n:], data2...))
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, body, []byte("hello worldhello world"))
	p.Reset()

}

func Test_Issue2(t *testing.T) {

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

	var body []byte
	var setting = Setting{
		MessageBegin: func(*Parser) {
		},
		URL: func(p *Parser, buf []byte) {
		},
		Status: func(*Parser, []byte) {
		},
		HeaderField: func(p *Parser, buf []byte) {
		},
		HeaderValue: func(p *Parser, buf []byte) {
		},
		HeadersComplete: func(p *Parser) {
		},
		Body: func(p *Parser, buf []byte) {
			body = append(body, buf...)
		},
		MessageComplete: func(p *Parser) {
			p.Reset()
		},
	}

	p := New(REQUEST)
	fmt.Printf("req_len=%d\n", len(data))
	// 一个POST 518，一共两个POST，第一次解析600字节，第二次解析剩余的
	data1, data2 := data[:300], data[300:]
	sucess, err := p.Execute(&setting, data1)
	if err != nil {
		panic(err.Error())
	}

	sucess, err = p.Execute(&setting, append(data1[sucess:], data2...))
	if err != nil {
		panic(err.Error())
	}

	p.Reset()

}
