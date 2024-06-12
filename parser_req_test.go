package httparser

import (
	"bytes"
	"fmt"
	"io"
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
		URL: func(p *Parser, buf []byte, _ int) {
			assert.Equal(t, string(url), string(buf))
		}, MessageBegin: func(p *Parser, _ int) {
			messageBegin = true
		},
	}

	_, err := p.Execute(setting, rsp)

	if err != nil {
		t.Fatalf("Execute:%v", err)
	}
	if messageBegin == false {
		t.Fatalf("messageBegin is false\n")
	}
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
			"Cache-Control: max-age=0\r\n\r\nb\r\nhello world\r\n0\r\n\r\n")

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
		MessageBegin: func(*Parser, int) {
			messageBegin = true
		},
		URL: func(p *Parser, buf []byte, _ int) {
			url = append(url, buf...)
		},
		Status: func(*Parser, []byte, int) {
			// 响应包才需要用到
		},
		HeaderField: func(p *Parser, buf []byte, _ int) {
			field = append(field, buf...)
		},
		HeaderValue: func(_ *Parser, buf []byte, _ int) {
			value = append(value, buf...)
		},
		HeadersComplete: func(*Parser, int) {
			headersComplete = true
		},
		Body: func(_ *Parser, buf []byte, _ int) {
			body = append(body, buf...)
		},
		MessageComplete: func(*Parser, int) {
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
	assert.True(t, p.EOF())

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
			"Cache-Control: max-age=0\r\n\r\n10\r\nhello world12345\r\n0\r\n\r\n")

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
		MessageBegin: func(*Parser, int) {
			messageBegin = true
		},
		URL: func(_ *Parser, buf []byte, _ int) {
			url = append(url, buf...)
		},
		Status: func(*Parser, []byte, int) {
			// 响应包才需要用到
		},
		HeaderField: func(_ *Parser, buf []byte, _ int) {
			field = append(field, buf...)
		},
		HeaderValue: func(_ *Parser, buf []byte, _ int) {
			value = append(value, buf...)
		},
		HeadersComplete: func(*Parser, int) {
			headersComplete = true
		},
		Body: func(_ *Parser, buf []byte, _ int) {
			body = append(body, buf...)
		},
		MessageComplete: func(_ *Parser, _ int) {
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
	assert.True(t, p.EOF())

	//fmt.Printf("##:%s", stateTab[p.currState])
}

// https://github.com/antlabs/httparser/issues/1
func Test_ParserRequest_chunked_segment(t *testing.T) {
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

	var body []byte
	var setting = Setting{
		Body: func(_ *Parser, buf []byte, _ int) {
			//fmt.Printf("###:%s\n", buf)
			body = append(body, buf...)
		},
	}

	p := New(REQUEST)

	for size := 120; size < 2*len(data); size++ {

		// 双缓冲buffer
		// 左边放溢出的，右边放本次读入数据, 这么设计可以减少内存拷贝
		tb := NewTwoBuf(size * 2)

		body = []byte{}
		totalSentBuf := []byte{} //存放送入Execute的总数据

		r := bytes.NewReader(data)
		for {

			// 取右边buffer
			buf := tb.Right()

			//模拟从异步io里面填充一块buffer
			n, err := r.Read(buf)
			if err == io.EOF {
				break
			}

			// 把溢出数据包含进来
			// 左边放需要重新解析数据，右边放新塞的buffer
			currSentData := tb.All(n)

			//解析
			success, err := p.Execute(&setting, currSentData)
			if err != nil {
				panic(err.Error() + fmt.Sprintf(" size:%d", size))
			}

			if success != len(currSentData) {
				// 测试用, 把送入解析器的buffer累加起来，最后验证下数据送得对不对
				totalSentBuf = append(totalSentBuf, currSentData[:success]...)

				tb.MoveLeft(currSentData[success:])
			} else {
				// 测试用
				totalSentBuf = append(totalSentBuf, currSentData...)

				tb.Reset()

			}

		}
		tb.Reset()

		b := assert.Equal(t, string(data), string(totalSentBuf))
		if !b {
			return
		}

		b = assert.Equal(t, body, []byte("hello worldhello world"))
		if !b {
			return
		}
	}

	p.Reset()

}
