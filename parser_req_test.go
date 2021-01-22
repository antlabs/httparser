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
}
