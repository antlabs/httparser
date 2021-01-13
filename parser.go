package httparser

import (
	"bytes"
	"errors"
	"strconv"
	"unicode"
)

var (
	ErrHTTPVersion    = errors.New("http version")
	ErrHTTPVersionNum = errors.New("http version number")
	ErrHTTPStatus     = errors.New("http status")
	ErrRspStatusLine  = errors.New("http rsp status line")
	ErrHeaderOverflow = errors.New("http header overflow")
	ErrNoEndLF        = errors.New("http there is no end symbol")
)

var (
	strTTPslash = []byte("TTP/")
)

var (
	contentLength       = []byte("Content-Length")
	maxHeaderSize int32 = 4096 //默认http header限制为4k
)

// http 1.1 or http 1.0解析器
type Parser struct {
	currState        state       //记录当前状态
	headerCurrState  headerState //记录http field状态
	major            uint8       //主版本号
	minor            uint8       //次版本号
	maxHeaderSize    int32       //最大头长度
	contentLength    int32       //content-length 值
	StatusCode       uint16      //状态码
	hasContentLength bool        //设置Content-Length头部
}

// 解析器构造函数
func New(t ReqOrRsp) *Parser {
	p := &Parser{}
	p.Init(t)
	return p
}

// 解析器Init函数
func (p *Parser) Init(t ReqOrRsp) {
	switch t {
	case REQUEST:
		p.currState = startReq
	case RESPONSE:
		p.currState = startRsp
	case BOTH:
		p.currState = startReqOrRsp
	}

	p.major = 0
	p.minor = 0
	p.maxHeaderSize = maxHeaderSize

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

// 注意
// 调用必须保证status-line的数据包是完整的
// (mtu 大约是1530左右，而status-line不会超过1个mtu)。
func (p *Parser) Execute(setting *Setting, buf []byte) (success int, err error) {
	currState := p.currState

	for i := 0; i < len(buf); i++ {
		c := buf[i]
	next:
		switch currState {
		case startReq:

		case startRsp:
			if c != 'H' {
				return 0, ErrHTTPVersion
			}

			if setting.MessageBegin != nil {
				setting.MessageBegin()
			}

			currState = rspHTTP

		case rspHTTP:
			if len(buf[i:]) < len(strTTPslash) {
				return 0, ErrHTTPVersion
			}

			if !bytes.Equal(buf[i:len(strTTPslash)+1], strTTPslash) {
				return 0, ErrHTTPVersion
			}
			i += len(strTTPslash) - 1
			currState = rspHTTPVersionNum

		case rspHTTPVersionNum:
			if len(buf[i:]) < 3 || !unicode.IsNumber(rune(buf[i])) || !unicode.IsNumber(rune(buf[i+2])) { // 1.1 or 1.0 or 0.9
				return 0, ErrHTTPVersionNum
			}

			p.major = buf[i] - '0'
			p.minor = buf[i+2] - '0'
			i += 2 // 3-1
			currState = rspStatusCode

		case rspStatusCode:
			for ; (buf[i] == ' ' || buf[i] == '\r' || buf[i] == '\n') && i < len(buf); i++ {
			}

			for ; buf[i] >= '0' && buf[i] <= '9'; i++ {
				p.StatusCode = uint16(int(p.StatusCode)*10 + int(buf[i]-'0'))
			}

			if i >= len(buf) {
				return 0, ErrHTTPStatus
			}

			currState = rspStatus
			goto next
		case rspStatus:
			start := i
			for ; (buf[start] == ' ' || buf[start] == '\r' || buf[start] == '\n') && start < len(buf); start++ {
			}

			end := start
			for ; !(buf[end] == ' ' || buf[end] == '\r' || buf[end] == '\n') && end < len(buf); end++ {
			}

			if end >= len(buf) || end+1 >= len(buf) {
				return 0, ErrRspStatusLine
			}

			switch {
			case buf[end] == '\r' && buf[end+1] == '\n':
				i = end + 1
			case buf[end] == '\r' || buf[end] == '\n':
				i = end

			}

			if setting.Status != nil {
				setting.Status(buf[start:end])
			}
			currState = headerField

		case headerField:
			if c == '\r' || c == '\n' {
				currState = headerDone
				continue
			}

			pos := bytes.IndexByte(buf[i:], ':')
			if pos == -1 {
				if int32(len(buf[i:])) > p.maxHeaderSize {
					return 0, ErrHeaderOverflow
				}
				return i, nil
			}

			if setting.HeaderField != nil {
				setting.HeaderField(buf[i : i+pos])
			}

			if bytes.Equal(buf[i:i+pos], contentLength) {
				p.headerCurrState = hContentLength
			}

			i += pos
			currState = headerValueDiscardWs
		case headerValueDiscardWs:
			// 只跳过一个' ' or '\t'
			if c == ' ' || c == '\t' {
				currState = headerValue
				continue
			}

			currState = headerValue

		case headerValue:
			pos := bytes.IndexAny(buf[i:], "\r\n")
			if pos == -1 {
				if int32(len(buf[i:])) > p.maxHeaderSize {
					return 0, ErrHeaderOverflow
				}
				return i, nil
			}
			if setting.HeaderValue != nil {
				setting.HeaderValue(buf[i : i+pos])
			}

			switch p.headerCurrState {
			case hContentLength:
				n, err := strconv.Atoi(BytesToString(buf[i : i+pos]))
				if err != nil {
					return i, err
				}

				p.contentLength = int32(n)
				p.hasContentLength = true
			}

			i += pos
			currState = headerField

		case headerDone:
			if c != '\n' {
				//return i, ErrNoEndLF
			}

			if setting.HeadersComplete != nil {
				setting.HeadersComplete()
			}

			if p.hasContentLength {
				if p.contentLength == 0 {
					if setting.MessageComplete != nil {
						setting.MessageComplete()
						return i, nil
					}
				} else {
					currState = httpBody
				}
			}

		case httpBody:
		}
	}

	p.currState = currState

	return 0, nil
}

func (p *Parser) SetMaxHeaderSize(size int32) {
	p.maxHeaderSize = size
}

func (p *Parser) Eof() bool {
	return true
}
