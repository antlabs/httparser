package httparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试解析状态行
// TODO status为空的测试数据
func Test_ParserResponse_StatusLine(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rsp := []byte("HTTP/1.1   200   OK\r\n")
	_, err := p.Execute(&Setting{Status: func(_ *Parser, buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func(*Parser) {
		messageBegin = true
	},
	}, rsp)

	assert.NoError(t, err)
	assert.Equal(t, p.Major, uint8(1))
	assert.Equal(t, p.Minor, uint8(1))
	assert.True(t, messageBegin)
}

// 测试解析http header
func Test_ParserResponse_HeaderField(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rsp := []byte("HTTP/1.1   200   OK\r\n" +
		"Connection: close\r\n\r\n")
	_, err := p.Execute(&Setting{Status: func(_ *Parser, buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func(*Parser) {
		messageBegin = true
	}, HeaderField: func(_ *Parser, buf []byte) {
		assert.Equal(t, buf, []byte("Connection"))
	},
	}, rsp)

	assert.NoError(t, err)
	assert.Equal(t, p.Major, uint8(1))
	assert.Equal(t, p.Minor, uint8(1))
	assert.True(t, messageBegin)
}

// 测试解析http header 和 http value
func Test_ParserResponse_HeaderValue(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rsp := []byte("HTTP/1.1   200   OK\r\n" +
		"Connection: close\r\n\r\n")
	_, err := p.Execute(&Setting{Status: func(_ *Parser, buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func(*Parser) {
		messageBegin = true
	}, HeaderField: func(_ *Parser, buf []byte) {
		assert.Equal(t, string(buf), "Connection")
	}, HeaderValue: func(_ *Parser, buf []byte) {
		assert.Equal(t, string(buf), "close")
	},
	}, rsp)

	assert.NoError(t, err)
	assert.Equal(t, p.Major, uint8(1))
	assert.Equal(t, p.Minor, uint8(1))
	assert.True(t, messageBegin)
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
	_, err := p.Execute(&Setting{Status: func(_ *Parser, buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func(*Parser) {
		messageBegin = true
	}, HeaderField: func(_ *Parser, buf []byte) {
		field = append(field, string(buf))
	}, HeaderValue: func(_ *Parser, buf []byte) {
		fieldValue = append(fieldValue, string(buf))
	},
	}, rsp)

	assert.NoError(t, err)
	assert.Equal(t, p.Major, uint8(1))
	assert.Equal(t, p.Minor, uint8(1))
	assert.Equal(t, field, []string{"Content-Length", "Connection"})
	assert.Equal(t, fieldValue, []string{"10", "close"})
	assert.True(t, messageBegin)
}

// 测试解析Content-Length body数据
func Test_ParserResponse_Content_Length_Body(t *testing.T) {

	p := New(RESPONSE)

	messageBegin := false
	rcvBuf := []byte{}
	setting := &Setting{Status: func(_ *Parser, buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func(_ *Parser) {
		messageBegin = true
	}, HeaderField: func(_ *Parser, buf []byte) {

	}, HeaderValue: func(_ *Parser, buf []byte) {
	}, Body: func(_ *Parser, buf []byte) {
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
		assert.NoError(t, err)
		if err != nil {
			return
		}
	}

	assert.Equal(t, rcvBuf, append(rsp[1], rsp[2]...))
	assert.Equal(t, p.Major, uint8(1))
	assert.Equal(t, p.Minor, uint8(1))
	assert.True(t, messageBegin)
}

func Test_ParserResponse_Chunked(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rcvBuf := []byte{}
	setting := &Setting{Status: func(_ *Parser, buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func(_ *Parser) {
		messageBegin = true
	}, HeaderField: func(_ *Parser, buf []byte) {

	}, HeaderValue: func(_ *Parser, buf []byte) {
	}, Body: func(_ *Parser, buf []byte) {
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
	assert.True(t, p.EOF())
}
