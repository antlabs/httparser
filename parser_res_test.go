package httparser

import (
	"bytes"
	"reflect"
	"testing"
)

// 测试解析状态行
// TODO status为空的测试数据
func Test_ParserResponse_StatusLine(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rsp := []byte("HTTP/1.1   200   OK\r\n")
	_, err := p.Execute(&Setting{
		Status: func(_ *Parser, buf []byte, _ int) {
			if !bytes.Equal(buf, []byte("OK")) {
				t.Error("status error")
			}
		}, MessageBegin: func(*Parser, int) {
			messageBegin = true
		},
	}, rsp)
	if err != nil {
		t.Error(err)
	}

	if p.Major != 1 {
		t.Errorf("major is %d, expect 1", p.Major)
	}

	if p.Minor != 1 {
		t.Errorf("minor is %d, expect 1", p.Minor)
	}

	if !messageBegin {
		t.Error("message begin is false, expect true")
	}
}

// 测试解析http header
func Test_ParserResponse_HeaderField(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rsp := []byte("HTTP/1.1   200   OK\r\n" +
		"Connection: close\r\n\r\n")
	_, err := p.Execute(&Setting{
		Status: func(_ *Parser, buf []byte, _ int) {
			if !bytes.Equal(buf, []byte("OK")) {
				t.Error("status error")
			}

		}, MessageBegin: func(*Parser, int) {
			messageBegin = true
		}, HeaderField: func(_ *Parser, buf []byte, _ int) {
			if !bytes.Equal(buf, []byte("Connection")) {
				t.Error("header field error")
			}
		},
	}, rsp)
	if err != nil {
		t.Error(err)
	}

	if p.Major != 1 {
		t.Errorf("major is %d, expect 1", p.Major)
	}
	if p.Minor != 1 {
		t.Errorf("minor is %d, expect 1", p.Minor)
	}
	if !messageBegin {
		t.Error("message begin is false, expect true")
	}
}

// 测试解析http header 和 http value
func Test_ParserResponse_HeaderValue(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rsp := []byte("HTTP/1.1   200   OK\r\n" +
		"Connection: close\r\n\r\n")
	_, err := p.Execute(&Setting{
		Status: func(_ *Parser, buf []byte, _ int) {
			if !bytes.Equal(buf, []byte("OK")) {
				t.Error("status error")
			}
		}, MessageBegin: func(*Parser, int) {
			messageBegin = true
		}, HeaderField: func(_ *Parser, buf []byte, _ int) {
			if !bytes.Equal(buf, []byte("Connection")) {
				t.Error("header field error")
			}

		}, HeaderValue: func(_ *Parser, buf []byte, _ int) {
			if !bytes.Equal(buf, []byte("close")) {
				t.Error("header value error")
			}

		},
	}, rsp)
	if err != nil {
		t.Error(err)
	}

	if p.Major != 1 {
		t.Errorf("major is %d, expect 1", p.Major)
	}
	if p.Minor != 1 {
		t.Errorf("minor is %d, expect 1", p.Minor)
	}
	if !messageBegin {
		t.Error("message begin is false, expect true")
	}
}

// 测试解析多个http header 和 http value。
func Test_ParserResponse_Multiple_HeaderValue(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	field := []string{}
	fieldValue := []string{}

	rsp := []byte("HTTP/1.1   200   OK\r\n" +
		"Content-Length: 10\r\n" +
		"Connection: close\r\n\r\n")
	_, err := p.Execute(&Setting{
		Status: func(_ *Parser, buf []byte, _ int) {
			if !bytes.Equal(buf, []byte("OK")) {
				t.Error("status error")
			}
		}, MessageBegin: func(*Parser, int) {
			messageBegin = true
		}, HeaderField: func(_ *Parser, buf []byte, _ int) {
			field = append(field, string(buf))
		}, HeaderValue: func(_ *Parser, buf []byte, _ int) {
			fieldValue = append(fieldValue, string(buf))
		},
	}, rsp)
	if err != nil {
		t.Error(err)
	}

	if p.Major != 1 {
		t.Errorf("major is %d, expect 1", p.Major)
	}
	if p.Minor != 1 {
		t.Errorf("minor is %d, expect 1", p.Minor)
	}
	if !messageBegin {
		t.Error("message begin is false, expect true")
	}

	if !reflect.DeepEqual(field, []string{"Content-Length", "Connection"}) {
		t.Errorf("field is %v, expect [Content-Length, Connection]", field)
	}

	if !reflect.DeepEqual(fieldValue, []string{"10", "close"}) {
		t.Errorf("fieldvalue is %v, expect [10, close]", fieldValue)
	}
}

// 测试解析Content-Length body数据
func Test_ParserResponse_Content_Length_Body(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rcvBuf := []byte{}
	setting := &Setting{
		Status: func(_ *Parser, buf []byte, _ int) {
			if !bytes.Equal(buf, []byte("OK")) {
				t.Error("status error")
			}

		}, MessageBegin: func(_ *Parser, _ int) {
			messageBegin = true
		}, HeaderField: func(_ *Parser, buf []byte, _ int) {
		}, HeaderValue: func(_ *Parser, buf []byte, _ int) {
		}, Body: func(_ *Parser, buf []byte, _ int) {
			rcvBuf = append(rcvBuf, buf...)
		},
	}

	var rsp [3][]byte
	rsp[0] = []byte("HTTP/1.1   200   OK\r\n" +
		"Content-Length: 20\r\n" +
		"Connection: close\r\n\r\n")

	rsp[1] = []byte("123456789a")
	rsp[2] = []byte("abcdefghij")

	for _, buf := range rsp {
		_, err := p.Execute(setting, buf)
		if err != nil {
			t.Error(err)
			return
		}
	}

	if !bytes.Equal(rcvBuf, append(rsp[1], rsp[2]...)) {
		t.Errorf("rcvBuf is %v, expect %v", rcvBuf, append(rsp[1], rsp[2]...))
	}

	if p.Major != 1 {
		t.Errorf("major is %d, expect 1", p.Major)
	}

	if p.Minor != 1 {
		t.Errorf("minor is %d, expect 1", p.Minor)
	}

	if !messageBegin {
		t.Error("message begin is false, expect true")
	}
}

func Test_ParserResponse_Chunked(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rcvBuf := []byte{}
	setting := &Setting{
		Status: func(_ *Parser, buf []byte, _ int) {
			if !bytes.Equal(buf, []byte("OK")) {
				t.Error("status error")
			}

		}, MessageBegin: func(_ *Parser, _ int) {
			messageBegin = true
		}, HeaderField: func(_ *Parser, buf []byte, _ int) {
		}, HeaderValue: func(_ *Parser, buf []byte, _ int) {
		}, Body: func(_ *Parser, buf []byte, _ int) {
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
			t.Error(err)
			return
		}

		parserTotal += rv
		sentTotal += len(buf)
	}

	if !bytes.Equal(rcvBuf, []byte("MozillaDeveloperNetworknew year")) {
		t.Errorf("rcvBuf is %v, expect %v", rcvBuf, []byte("MozillaDeveloperNetworknew year"))
	}

	if p.Major != 1 {
		t.Errorf("major is %d, expect 1", p.Major)
	}

	if p.Minor != 1 {
		t.Errorf("minor is %d, expect 1", p.Minor)
	}
	if !messageBegin {
		t.Error("message begin is false, expect true")
	}

	if sentTotal != parserTotal {
		t.Errorf("sentTotal is %d, parserTotal is %d", sentTotal, parserTotal)
	}

	if !p.EOF() {
		t.Error("EOF is true, expect false")
	}
}
