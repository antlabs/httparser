package httparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试请求body
func Test_ParserResponse_RequestBody_BOTH(t *testing.T) {
	p := New(BOTH)

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
		MessageBegin: func(*Parser) {
			messageBegin = true
		},
		URL: func(p *Parser, buf []byte) {
			url = append(url, buf...)
		},
		Status: func(*Parser, []byte) {
			// 响应包才需要用到
		},
		HeaderField: func(p *Parser, buf []byte) {
			field = append(field, buf...)
		},
		HeaderValue: func(p *Parser, buf []byte) {
			value = append(value, buf...)
		},
		HeadersComplete: func(*Parser) {
			headersComplete = true
		},
		Body: func(p *Parser, buf []byte) {
			body = append(body, buf...)
		},
		MessageComplete: func(p *Parser) {
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

func Test_ParserResponse_Chunked_Both(t *testing.T) {
	p := New(BOTH)

	messageBegin := false
	rcvBuf := []byte{}
	setting := &Setting{Status: func(p *Parser, buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func(p *Parser) {
		messageBegin = true
	}, HeaderField: func(p *Parser, buf []byte) {

	}, HeaderValue: func(p *Parser, buf []byte) {
	}, Body: func(p *Parser, buf []byte) {
		rcvBuf = append(rcvBuf, buf...)
	},
	}

	var rsp [3]string
	rsp[0] = "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Transfer-Encoding: chunked\r\n\r\n" +
		"7\r\n" +
		"Mozilla\r\n" +
		"9\r\n" +
		"Developer\r\n" +
		"7\r\n" +
		"Network\r\n"

	rsp[1] = "8\r\n" +
		"new year\r\n"

	rsp[2] = "0\r\n\r\n"
	sentTotal := 0
	parserTotal := 0
	for _, buf := range rsp {
		rv, err := p.Execute(setting, []byte(buf))
		assert.NoError(t, err)
		if err != nil {
			return
		}

		parserTotal += rv
		sentTotal += len(buf)
	}

	assert.Equal(t, rcvBuf, []byte("MozillaDeveloperNetworknew year"))
	assert.Equal(t, p.Major, uint8(1))
	assert.Equal(t, p.Minor, uint8(1))
	assert.True(t, messageBegin)
	assert.Equal(t, sentTotal, parserTotal)
	assert.True(t, p.Eof())
}
