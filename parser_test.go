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
		"Connection: close\r\n")
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
		"Connection: close\r\n")
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
