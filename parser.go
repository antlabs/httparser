package httparser

import (
	"bytes"
	"errors"
	"strconv"
	"unicode"
)

var (
	ErrHTTPVersion     = errors.New("http version")
	ErrHTTPVersionNum  = errors.New("http version number")
	ErrHeaderOverflow  = errors.New("http header overflow")
	ErrNoEndLF         = errors.New("http there is no end symbol")
	ErrChunkSize       = errors.New("http wrong chunk size")
	ErrReqMethod       = errors.New("http request wrong method")
	ErrRequestLineCRLF = errors.New("http request line wrong CRLF")
)

var (
	strTTPslash = []byte("TTP/")
)

var (
	contentLength          = []byte("Content-Length")
	transferEncoding       = []byte("Transfer-Encoding")
	bytesChunked           = []byte("chunked")
	bytesConnection        = []byte("Connection")
	bytesClose             = []byte("close")
	bytesUpgrade           = []byte("upgrade")
	MaxHeaderSize    int32 = 4096 //默认http header单行最大限制为4k
)

// http 1.1 or http 1.0解析器
type Parser struct {
	hType               ReqOrRsp     //解析器的类型，解析请求还是响应
	currState           state        //记录当前状态
	headerCurrState     headerState  //记录http field状态
	Major               uint8        //主版本号
	Minor               uint8        //次版本号
	MaxHeaderSize       int32        //最大头长度
	contentLength       int32        //content-length 值
	StatusCode          uint16       //状态码
	hasContentLength    bool         //设置Content-Length头部
	hasTransferEncoding bool         //transferEncoding头部
	hasClose            bool         // Connection: close
	hasUpgrade          bool         // Connection: Upgrade
	trailing            trailerState //trailer的状态
	userData            interface{}
}

// 解析器构造函数
func New(t ReqOrRsp) *Parser {
	p := &Parser{}
	p.Init(t)
	return p
}

// 解析器Init函数
func (p *Parser) Init(t ReqOrRsp) {

	p.currState = newState(t)

	p.hType = t
	p.Major = 0
	p.Minor = 0
	p.MaxHeaderSize = MaxHeaderSize

}

// 一般情况，可以使用Setting里面函数闭包特性捕获调用者私有变量
// 好像提供SetUserData没有必要性
// 但是从
// 1.节约内存角度
// 2.数据和行为解耦的角度，提供一个Setting函数集，通过SetUserData，保存调用者私有变量
// 通过GetUserData拿需要的变量，还是比较爽的，不需要写一堆闭包
func (p *Parser) SetUserData(d interface{}) {
	p.userData = d
}

// 获取SetuserData函数设置的私有变量
func (p *Parser) GetUserData() interface{} {
	return p.userData
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

// 响应行
// https://tools.ietf.org/html/rfc7230#section-3.1.2 状态行
// status-line = HTTP-version SP status-code SP reason-phrase CRLF

// 请求行
// https://tools.ietf.org/html/rfc7230#section-3.1.1
// method SP request-target SP HTTP-version CRLF

// 设计思路修改
// 为了适应流量解析的场景，状态机的状态会更碎一点

func (p *Parser) Execute(setting *Setting, buf []byte) (success int, err error) {
	currState := p.currState

	chunkDataStartIndex := 0
	urlStartIndex := 0
	reasonPhraseIndex := 0

	i := 0
	c := byte(0)

	for ; i < len(buf); i++ {
		c = buf[i]

		//fmt.Printf("---->debug state(%s):(%s)\n", currState, buf[i:])
	reExec:
		switch currState {
		case startReqOrRsp:
			if c == 'H' {
				if setting.MessageBegin != nil {
					setting.MessageBegin(p)
				}
				currState = rspHTTP
				continue
			}
			currState = startReq
			fallthrough
		case startReq:
			if c == '\r' || c == '\n' {
				continue
			}

			if token[c] == 0 {
				return 0, ErrReqMethod
			}

			currState = reqMethod
			if setting.MessageBegin != nil {
				setting.MessageBegin(p)
			}

		case reqMethod:
			if token[c] == 0 {
				if c == ' ' || c == '\t' {
					currState = reqMethodAfterSP
					continue
				}

				return i, ErrReqMethod
			}

			// 维持reqMethod状态不变
		case reqMethodAfterSP:
			if c != ' ' && c != '\t' {
				urlStartIndex = i
				currState = reqURL
			}

		case reqURL:
			if c == ' ' || c == '\t' {
				currState = reqURLAfterSP
				if setting.URL != nil {
					setting.URL(p, buf[urlStartIndex:i])
				}
			}

		case reqURLAfterSP:
			if c != ' ' && c != '\t' {
				currState = reqHTTPVersion
			}
		case reqHTTPVersion:
			if c == '/' {
				currState = reqHTTPVersionMajor
			}
		case reqHTTPVersionMajor:
			p.Major = c - '0'
			currState = reqHTTPVersionDot
		case reqHTTPVersionDot:
			currState = reqHTTPVersionMinor
		case reqHTTPVersionMinor:
			if c == '\r' {
				currState = reqRequestLineAlomstDone
				continue
			}
			p.Minor = c - '0'

		case reqRequestLineAlomstDone:
			if c != '\n' {
				return 0, ErrRequestLineCRLF
			}

			currState = headerField

		case startRsp:
			if c != 'H' {
				return 0, ErrHTTPVersion
			}

			if setting.MessageBegin != nil {
				setting.MessageBegin(p)
			}

			currState = rspHTTP

		case rspHTTP:
			if len(buf[i:]) < len(strTTPslash) {
				p.currState = currState
				return i, nil
			}

			if !bytes.Equal(buf[i:i+len(strTTPslash)], strTTPslash) {
				return 0, ErrHTTPVersion
			}

			i += len(strTTPslash) - 1
			currState = rspHTTPVersionNum

		case rspHTTPVersionNum:
			// 1.1 or 1.0 or 0.9
			if len(buf[i:]) < 3 {
				p.currState = currState
				return i, nil
			}

			if !unicode.IsNumber(rune(buf[i])) || !unicode.IsNumber(rune(buf[i+2])) {
				return 0, ErrHTTPVersionNum
			}

			p.Major = buf[i] - '0'
			p.Minor = buf[i+2] - '0'
			i += 2 // 3-1
			currState = rspHTTPVersionNumAfterSP

		case rspHTTPVersionNumAfterSP:
			if c == ' ' || c == '\r' || c == '\n' {
				continue
			}

			currState = rspStatusCode
			goto reExec
		case rspStatusCode:

			if c >= '0' && c <= '9' {
				p.StatusCode = uint16(int(p.StatusCode)*10 + int(c-'0'))
				continue
			}

			currState = rspStatusCodeAfterSP
			goto reExec
		case rspStatusCodeAfterSP:
			if c == ' ' || c == '\r' || c == '\n' {
				continue
			}

			currState = rspStatus
			goto reExec
		case rspStatus:
			if reasonPhraseIndex == 0 {
				reasonPhraseIndex = i
			}

			if c == '\r' {
				if setting.Status != nil {
					setting.Status(p, buf[reasonPhraseIndex:i])
				}
				currState = rspStatusAfterSP
				continue
			}

			if c == '\n' {
				if setting.Status != nil {
					setting.Status(p, buf[reasonPhraseIndex:i])
				}
				currState = headerField
			}

		case rspStatusAfterSP:
			currState = headerField

		case headerField:
			if c == '\r' || c == '\n' {
				currState = headerDone
				continue
			}

			pos := bytes.IndexByte(buf[i:], ':')
			if pos == -1 {
				if int32(len(buf[i:])) > p.MaxHeaderSize {
					return 0, ErrHeaderOverflow
				}

				p.currState = headerField
				return i, nil
			}

			field := buf[i : i+pos]
			if setting.HeaderField != nil {
				setting.HeaderField(p, field)
			}

			c2 := c | 0x20
			if c2 == 'c' || c2 == 't' {
				if bytes.EqualFold(field, contentLength) {
					// Content-Length
					p.headerCurrState = hContentLength
				} else if bytes.EqualFold(field, transferEncoding) {
					// Transfer-Encoding
					p.headerCurrState = hTransferEncoding
				} else if bytes.EqualFold(field, bytesConnection) {
					// Connection
					p.headerCurrState = hConnection
				} else {
					// general
					p.headerCurrState = hGeneral
				}
			} else {
				p.headerCurrState = hGeneral
			}

			i += pos
			currState = headerValueDiscardWs
		case headerValueDiscardWs:
			// 只跳过一个' ' or '\t'
			currState = headerValue
			if c == ' ' || c == '\t' {
				continue
			}
			goto reExec

			// 解析http value
		case headerValue:
			end := bytes.IndexAny(buf[i:], "\r\n")
			if end == -1 {
				if int32(len(buf[i:])) > p.MaxHeaderSize {
					return 0, ErrHeaderOverflow
				}
				p.currState = headerValueDiscardWs
				return i, nil
			}

			hValue := buf[i : i+end]
			if setting.HeaderValue != nil {
				setting.HeaderValue(p, hValue)
			}

			switch p.headerCurrState {
			case hConnection:
				switch {
				case bytes.Index(hValue, bytesClose) != -1:
					p.hasClose = true
				case bytes.Index(hValue, bytesUpgrade) != -1:
					p.hasUpgrade = true
				}
			case hContentLength:
				n, err := strconv.Atoi(BytesToString(hValue))
				if err != nil {
					return i, err
				}

				p.contentLength = int32(n)
				p.hasContentLength = true
				p.headerCurrState = hGeneral
			case hTransferEncoding:
				pos := bytes.Index(hValue, bytesChunked)
				// 没有chunked值，归类到通用http header
				if pos == -1 {
					p.headerCurrState = hGeneral
				}
				p.hasTransferEncoding = true
			}

			i += end

			c = buf[i]
			currState = headerValueStartOWS
			// 进入header value的OWS
			fallthrough
		case headerValueStartOWS:
			if c == '\r' {
				currState = headerValueOWS
				continue
			}

			// 不是'\r'的情况，继续往下判断
			fallthrough
		case headerValueOWS:
			currState = headerField
			if c == '\n' {
				continue
			}

			// 不是'\n'也许是headerField的数据
			goto reExec

		case headerDone:
			if c != '\n' {
				return i, ErrNoEndLF
			}

			if p.trailing == needParserTrailer {
				if setting.MessageComplete != nil {
					setting.MessageComplete(p)
				}

				currState = messageDone
				goto reExec
			}

			if setting.HeadersComplete != nil {
				setting.HeadersComplete(p)
			}

			// TODO hasContentLength, hasTransferEncoding同时为true
			if p.hasContentLength {
				// 如果contentLength 等于0，说明body的内容为空，可以直接退出
				if p.contentLength == 0 {
					if setting.MessageComplete != nil {
						setting.MessageComplete(p)
						return i, nil
					}
				}
				currState = httpBody
				continue
			}

			if p.hasTransferEncoding {
				currState = chunkedSizeStart
				continue
			}

			if p.hasClose {
				currState = messageDone
				continue
			}

			if p.contentLength == 0 {
				if setting.MessageComplete != nil {
					currState = newState(p.hType)
					setting.MessageComplete(p)
				}
				continue
			}

			if p.Eof() {
				currState = messageDone
			}
		case httpBody:
			if p.hasContentLength {
				nread := min(int32(len(buf[i:])), p.contentLength)
				if setting.Body != nil && nread > 0 {
					setting.Body(p, buf[i:int32(i)+nread])
				}

				p.contentLength -= nread

				if p.contentLength == 0 {
					currState = messageDone
				}

				i += int(nread) - 1

				// httpBody没有后继结点,所以这里发现数据消费完, 调用下MessageComplete方法
				if p.contentLength == 0 {
					if setting.MessageComplete != nil {
						setting.MessageComplete(p)
					}
				}
			}

		case chunkedSizeStart:
			l := unhex[c]
			if l == -1 {
				return 0, ErrChunkSize
			}

			p.contentLength = int32(l)
			currState = chunkedSize

		case chunkedSize:
			if c == '\r' {
				currState = chunkedSizeAlmostDone
				continue
			}

			l := unhex[c]
			if l == -1 {
				if c == ';' {
					currState = chunkedExt
					continue
				}

				return 0, ErrChunkSize
			}

			p.contentLength = p.contentLength*16 + int32(l)

		case chunkedExt:
			// 忽略chunked ext
			if c == '\r' {
				currState = chunkedSizeAlmostDone
			}

		case chunkedSizeAlmostDone:
			if p.contentLength == 0 {

				p.trailing = needParserTrailer
				currState = headerField

				continue
			}

			chunkDataStartIndex = i + 1
			currState = chunkedData

		case chunkedData:
			nread := min(int32(len(buf[i:])), p.contentLength)
			if setting.Body != nil && nread > 0 {
				setting.Body(p, buf[chunkDataStartIndex:int32(chunkDataStartIndex)+nread])
			}

			p.contentLength -= nread

			if p.contentLength == 0 {
				currState = chunkedDataAlmostDone
			}

			if nread > 0 {
				i += int(nread) - 1
			}
		case chunkedDataAlmostDone:
			currState = chunkedDataDone
		case chunkedDataDone:
			currState = chunkedSizeStart
		case messageDone:
			// 规范的chunked包是以\r\n结尾的
			if c == '\r' || c == '\n' {
				continue
			}

			currState = newState(p.hType)
			p.Reset()
		}

	}

	switch currState {
	case reqURL:
		if setting.URL != nil && len(buf[urlStartIndex:len(buf)]) > 0 {
			setting.URL(p, buf[urlStartIndex:len(buf)])
		}

	case rspStatus:
		if setting.Status != nil && len(buf[reasonPhraseIndex:i]) > 0 {
			setting.Status(p, buf[urlStartIndex:len(buf)])
		}
	}

	p.currState = currState
	return i, nil
}

func newState(t ReqOrRsp) state {
	switch t {
	case REQUEST:
		return startReq
	case RESPONSE:
		return startRsp
	case BOTH:
		return startReqOrRsp
	}
	return startReqOrRsp
}

func (p *Parser) Reset() {
	p.currState = newState(p.hType)
	p.headerCurrState = hGeneral
	p.Major = 0
	p.Minor = 0
	//p.MaxHeaderSize
	p.contentLength = 0
	p.StatusCode = 0
	p.hasContentLength = false
	p.hasTransferEncoding = false
	p.trailing = defaultTrailer
}

// debug专用
func (p *Parser) Status() string {
	return stateTab[p.currState]
}

func (p *Parser) Eof() bool {
	if p.hType == REQUEST {
		return true
	}

	return p.currState == messageDone
}

func min(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}
