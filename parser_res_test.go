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
	_, err := p.Execute(&Setting{Status: func(buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func() {
		messageBegin = true
	},
	}, rsp)

	assert.NoError(t, err)
	assert.Equal(t, p.major, uint8(1))
	assert.Equal(t, p.minor, uint8(1))
	assert.True(t, messageBegin)
}

// 测试解析http header
func Test_ParserResponse_HeaderField(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rsp := []byte("HTTP/1.1   200   OK\r\n" +
		"Connection: close\r\n\r\n")
	_, err := p.Execute(&Setting{Status: func(buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func() {
		messageBegin = true
	}, HeaderField: func(buf []byte) {
		assert.Equal(t, buf, []byte("Connection"))
	},
	}, rsp)

	assert.NoError(t, err)
	assert.Equal(t, p.major, uint8(1))
	assert.Equal(t, p.minor, uint8(1))
	assert.True(t, messageBegin)
}

// 测试解析http header 和 http value
func Test_ParserResponse_HeaderValue(t *testing.T) {
	p := New(RESPONSE)

	messageBegin := false
	rsp := []byte("HTTP/1.1   200   OK\r\n" +
		"Connection: close\r\n\r\n")
	_, err := p.Execute(&Setting{Status: func(buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func() {
		messageBegin = true
	}, HeaderField: func(buf []byte) {
		assert.Equal(t, string(buf), "Connection")
	}, HeaderValue: func(buf []byte) {
		assert.Equal(t, string(buf), "close")
	},
	}, rsp)

	assert.NoError(t, err)
	assert.Equal(t, p.major, uint8(1))
	assert.Equal(t, p.minor, uint8(1))
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
	_, err := p.Execute(&Setting{Status: func(buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func() {
		messageBegin = true
	}, HeaderField: func(buf []byte) {
		field = append(field, string(buf))
	}, HeaderValue: func(buf []byte) {
		fieldValue = append(fieldValue, string(buf))
	},
	}, rsp)

	assert.NoError(t, err)
	assert.Equal(t, p.major, uint8(1))
	assert.Equal(t, p.minor, uint8(1))
	assert.Equal(t, field, []string{"Content-Length", "Connection"})
	assert.Equal(t, fieldValue, []string{"10", "close"})
	assert.True(t, messageBegin)
}

// 测试解析Content-Length body数据
func Test_ParserResponse_Content_Length_Body(t *testing.T) {

	p := New(RESPONSE)

	messageBegin := false
	rcvBuf := []byte{}
	setting := &Setting{Status: func(buf []byte) {
		assert.Equal(t, buf, []byte("OK"))
	}, MessageBegin: func() {
		messageBegin = true
	}, HeaderField: func(buf []byte) {

	}, HeaderValue: func(buf []byte) {
	}, Body: func(buf []byte) {
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
	assert.Equal(t, p.major, uint8(1))
	assert.Equal(t, p.minor, uint8(1))
	assert.True(t, messageBegin)
}
