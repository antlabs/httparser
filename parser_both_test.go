package httparser

import (
	"bytes"
	"testing"
)

// 测试请求body
func Test_ParserResponse_RequestBody_BOTH(t *testing.T) {
	p := New(BOTH)

	data := []byte(
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
	value2 := "github.com" +
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
	field2 := []byte(
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

	setting := Setting{
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
		HeaderValue: func(p *Parser, buf []byte, _ int) {
			value = append(value, buf...)
		},
		HeadersComplete: func(*Parser, int) {
			headersComplete = true
		},
		Body: func(p *Parser, buf []byte, _ int) {
			body = append(body, buf...)
		},
		MessageComplete: func(p *Parser, _ int) {
			messageComplete = true
		},
	}

	i, err := p.Execute(&setting, data)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(url, []byte("/joyent/http-parser")) {
		t.Fatal("url error")
	}

	if i != len(data) {
		t.Fatal("data length error")
	}

	// 总数据长度
	if i != len(data) {
		t.Errorf("i(%d) != len(data)(%d)\n", i, len(data))
		return
	}

	if !bytes.Equal(field, field2) {
		t.Errorf("field error: %s\n", field)
		return
	}
	if string(value) != value2 {
		t.Errorf("value error: %s\n", value)
		return
	}

	if string(body) != body2 {
		t.Errorf("body error: %s\n", body)
		return
	}

	if !messageBegin {
		t.Error("message begin is false, expect true")
		return
	}

	if !messageComplete {
		t.Error("message complete is false, expect true")
		return
	}

	if !headersComplete {
		t.Error("headers complete is false, expect true")
		return
	}
	if !p.EOF() {
		t.Error("EOF is false, expect true")
		return
	}

	// fmt.Printf("##:%s", stateTab[p.currState])
}

func Test_ParserResponse_Chunked_Both(t *testing.T) {
	p := New(BOTH)

	messageBegin := false
	rcvBuf := []byte{}
	setting := &Setting{
		Status: func(p *Parser, buf []byte, _ int) {
			if !bytes.Equal(buf, []byte("OK")) {
				t.Error("status error")
			}
		}, MessageBegin: func(p *Parser, _ int) {
			messageBegin = true
		}, HeaderField: func(p *Parser, buf []byte, _ int) {
		}, HeaderValue: func(p *Parser, buf []byte, _ int) {
		}, Body: func(p *Parser, buf []byte, _ int) {
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
		if err != nil {
			t.Errorf("err:%s", err)
			return
		}

		parserTotal += rv
		sentTotal += len(buf)
	}

	if !bytes.Equal(rcvBuf, []byte("MozillaDeveloperNetworknew year")) {
		t.Errorf("rcvBuf:%s", rcvBuf)
		return
	}
	if p.Major != 1 || p.Minor != 1 {
		t.Errorf("major:%d, minor:%d", p.Major, p.Minor)
		return
	}

	if !messageBegin {
		t.Error("message begin is false, expect true")
		return
	}

	if sentTotal != parserTotal {
		t.Errorf("sendTotal:%d, parserTotal:%d", sentTotal, parserTotal)
		return
	}

	if !p.EOF() {
		t.Error("EOF is false, expect true")
		return
	}
}
