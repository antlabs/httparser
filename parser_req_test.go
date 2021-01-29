package httparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试请求行
func Test_ParserResponse_RequestLine(t *testing.T) {
	p := New(REQUEST)

	messageBegin := false
	rsp := []byte("GET /somedir/page.html HTTP/1.1\r\n")

	url := []byte("/somedir/page.html")
	setting := &Setting{
		URL: func(buf []byte) {
			assert.Equal(t, url, buf)
		}, MessageBegin: func() {
			messageBegin = true
		},
	}

	_, err := p.Execute(setting, rsp)

	assert.NoError(t, err)
	assert.True(t, messageBegin)
}

// 测试请求body
func Test_ParserResponse_RequestBody(t *testing.T) {
	p := New(REQUEST)

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

	messageBegin := false
	messageComplete := false
	headersComplete := false

	var body []byte
	var url []byte
	var field []byte
	var value []byte

	body2 := "hello world"
	var value2 = "github.com" +
		"1" +
		"gzip, deflate, sdch" +
		"ru-RU,ru;q=0.8,en-US;q=0.6,en;q=0.4" +
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/39.0.2171.65 Safari/537.36" +
		"text/html,application/xhtml+xml,application/xml;q=0.9," +
		"image/webp,*/*;q=0.8" +
		"https://github.com/joyent/http-parser" +
		"keep-alive" +
		"chunked" +
		"max-age=0"
	var field2 = []byte(
		"Host" +
			"DNT" +
			"Accept-Encoding" +
			"Accept-Language" +
			"User-Agent" +
			"Accept" +
			"Referer" +
			"Connection" +
			"Transfer-Encoding" +
			"Cache-Control")

	var setting = Setting{
		MessageBegin: func() {
			messageBegin = true
		},
		URL: func(buf []byte) {
			url = append(url, buf...)
		},
		Status: func([]byte) {
			// 响应包才需要用到
		},
		HeaderField: func(buf []byte) {
			field = append(field, buf...)
		},
		HeaderValue: func(buf []byte) {
			value = append(value, buf...)
		},
		HeadersComplete: func() {
			headersComplete = true
		},
		Body: func(buf []byte) {
			body = append(body, buf...)
		},
		MessageComplete: func() {
			messageComplete = true
		},
	}

	i, err := p.Execute(&setting, data)
	assert.NoError(t, err)
	assert.Equal(t, url, []byte("/joyent/http-parser")) //url
	assert.Equal(t, i, len(data))                       //总数据长度
	assert.Equal(t, field, field2)                      //header field
	assert.Equal(t, string(value), value2)              //header field
	assert.Equal(t, string(body), body2)                //chunked body
	assert.True(t, messageBegin)
	assert.True(t, messageComplete)
	assert.True(t, headersComplete)
	assert.True(t, p.Eof())

	//fmt.Printf("##:%s", stateTab[p.currState])
}

// 测试请求body2, chunked是两位数的
func Test_ParserResponse_RequestBody2(t *testing.T) {
	p := New(REQUEST)

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
			"Cache-Control: max-age=0\r\n\r\n10\r\nhello world12345\r\n0\r\n")

	messageBegin := false
	messageComplete := false
	headersComplete := false

	var body []byte
	var url []byte
	var field []byte
	var value []byte

	body2 := "hello world12345"
	var value2 = "github.com" +
		"1" +
		"gzip, deflate, sdch" +
		"ru-RU,ru;q=0.8,en-US;q=0.6,en;q=0.4" +
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/39.0.2171.65 Safari/537.36" +
		"text/html,application/xhtml+xml,application/xml;q=0.9," +
		"image/webp,*/*;q=0.8" +
		"https://github.com/joyent/http-parser" +
		"keep-alive" +
		"chunked" +
		"max-age=0"
	var field2 = []byte(
		"Host" +
			"DNT" +
			"Accept-Encoding" +
			"Accept-Language" +
			"User-Agent" +
			"Accept" +
			"Referer" +
			"Connection" +
			"Transfer-Encoding" +
			"Cache-Control")

	var setting = Setting{
		MessageBegin: func() {
			messageBegin = true
		},
		URL: func(buf []byte) {
			url = append(url, buf...)
		},
		Status: func([]byte) {
			// 响应包才需要用到
		},
		HeaderField: func(buf []byte) {
			field = append(field, buf...)
		},
		HeaderValue: func(buf []byte) {
			value = append(value, buf...)
		},
		HeadersComplete: func() {
			headersComplete = true
		},
		Body: func(buf []byte) {
			body = append(body, buf...)
		},
		MessageComplete: func() {
			messageComplete = true
		},
	}

	i, err := p.Execute(&setting, data)
	assert.NoError(t, err)
	assert.Equal(t, url, []byte("/joyent/http-parser")) //url
	assert.Equal(t, i, len(data))                       //总数据长度
	assert.Equal(t, field, field2)                      //header field
	assert.Equal(t, string(value), value2)              //header field
	assert.Equal(t, string(body), body2)                //chunked body
	assert.True(t, messageBegin)
	assert.True(t, messageComplete)
	assert.True(t, headersComplete)
	assert.True(t, p.Eof())

	//fmt.Printf("##:%s", stateTab[p.currState])
}
