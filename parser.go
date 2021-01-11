package httparser

import (
	"bytes"
	"errors"
	"unicode"
)

var (
	ErrHTTPVersion = errors.New("http version")
)

var (
	strTTPslash = []byte("TTP/")
)

type Parser struct {
	currState state //记录当前状态
	major     uint8
	minor     uint8
}

func New(t ReqOrRsp) *Parser {
	p := &Parser{}
	switch t {
	case REQUEST:
		p.currState = startReq
	case RESPONSE:
		p.currState = startRsp
	case BOTH:
		p.currState = startReqOrRsp
	}
	return p
}

// Execute传递setting参数, 该API 设计成现有形式有如下原因:
// setting如果通过New函数传递, Parser内存占用会多8 * 6的byte
// 为了减小Parser的大小，setting放至Execute函数里面传递

// 请求报文示例
// GET /somedir/page.html HTTP/1.1
// Host: www.someschool.edu
// Connection: close
// User-agent: Mozilla/5.0
// Accept-language: fr

// 响应报文示例
// HTTP/1.1   200   OK
// Connection: close
// Date: Tue,  09 Aug  2011  15:44:04 GMT
// Server: Apache/2.2.3 (CentOS)
// Last-Modified: Tue, 09 Aug 2011 15:11:03 GMT
// Content-Length: 6821
// Content-Type: text/html

// https://tools.ietf.org/html/rfc7230#section-3.1.2 状态行
// status-line = HTTP-version SP status-code SP reason-phrase CRLF
func (p *Parser) Execute(setting *Setting, buf []byte) error {
	currState := p.currState

	for i := 0; i < len(buf); {
		switch currState {
		case startReq:

		case startRsp:
			c := buf[i]
			if c != 'H' {
				return ErrHTTPVersion
			}

			if setting.MessageBegin != nil {
				setting.MessageBegin(p)
			}

			i++
		case rspHTTP:
			if len(buf[i:]) < len(strTTPslash) {
				return ErrHTTPVersion
			}

			if !bytes.Equal(buf[i:], strTTPslash) {
				return ErrHTTPVersion
			}
			i += len(strTTPslash)

		case rspHTTPVersionNum:
			if len(buf[i:]) < 3 || !unicode.IsNumber(buf[0]) || !unicode.IsNumber(buf[2]) { // 1.1 or 1.0 or 0.9
				return ErrHTTPVersion
			}

			p.major = buf[0] - '0'
			p.minor = buf[2] - '0'
			i += 3
		}
	}
}
