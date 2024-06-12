// Copyright 2021 guonaihong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httparser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

var (
	// ErrMethod 错误的method
	ErrMethod = errors.New("http method fail")
	// ErrStatusLineHTTP 状态行前面的HTTP错误
	ErrStatusLineHTTP = errors.New("http status line http")
	// ErrHTTPVersionNum 错误的http版本号
	ErrHTTPVersionNum = errors.New("http version number")
	// ErrHeaderOverflow http header包太大
	ErrHeaderOverflow = errors.New("http header overflow")
	// ErrNoEndLF 没有\n
	ErrNoEndLF = errors.New("http there is no end symbol")
	// ErrChunkSize chunked表示长度的字符串错了
	ErrChunkSize = errors.New("http wrong chunk size")
	// ErrReqMethod 错误的请求方法字符串
	ErrReqMethod = errors.New("http request wrong method")
	// ErrRequestLineLF 请求行没有\n
	ErrRequestLineLF = errors.New("http request line wrong LF")
)

var (
	strTTPslash = []byte("TTP/")
)

var (
	bytesCommaSep         = []byte(",")
	bytesContentLength    = []byte("Content-Length")
	bytesTransferEncoding = []byte("Transfer-Encoding")
	bytesChunked          = []byte("chunked")
	bytesConnection       = []byte("Connection")
	bytesClose            = []byte("close")
	bytesUpgrade          = []byte("upgrade")
	bytesSpace            = []byte(" ")
	// MaxHeaderSize 表示 http header单行最大限制为4k
	MaxHeaderSize int32 = 4096
)

const unused = -1

// Parser http 1.1 or http 1.0解析器
type Parser struct {
	hType                ReqOrRsp    //解析器的类型，解析请求还是响应
	Method               Method      //记录request的method
	currState            state       //记录当前状态
	headerCurrState      headerState //记录http field状态
	Major                uint8       //主版本号
	Minor                uint8       //次版本号
	MaxHeaderSize        int32       //最大头长度
	contentLength        int32       //content-length 值
	StatusCode           uint16      //状态码
	hasContentLength     bool        //设置Content-Length头部
	hasTransferEncoding  bool        //transferEncoding头部
	hasConnectionClose   bool        //Connection: close
	hasUpgrade           bool        //Upgrade: xx
	hasConnectionUpgrade bool        //Connection: Upgrade
	hasTrailing          bool        //有trailer的包
	callMessageComplete  bool        //记录MessageComplete是否被调用

	Upgrade bool //从http升级为别的协议, 比如websocket

	userData interface{}
}

// New 解析器构造函数
func New(t ReqOrRsp) *Parser {
	p := &Parser{}
	p.Init(t)
	return p
}

// Init 解析器Init函数
func (p *Parser) Init(t ReqOrRsp) {

	p.currState = newState(t)

	p.hType = t
	p.Major = 0
	p.Minor = 0
	p.contentLength = unused
	p.MaxHeaderSize = MaxHeaderSize

}

// ReadyUpgradeData 如果ReadyUpgradeData为true 说明已经有Upgrade Data数据, 并且http数据已经成功解析完成
func (p *Parser) ReadyUpgradeData() bool {
	return p.callMessageComplete && p.Upgrade
}

func (p *Parser) complete(s *Setting, pos int) {
	p.callMessageComplete = true
	if s.MessageComplete != nil {
		s.MessageComplete(p, pos)
	}
}

// SetUserData 的设计出发点和作用
// 一般情况，可以使用Setting里面函数闭包特性捕获调用者私有变量
// 好像提供SetUserData没有必要性
// 但是从
// 1.节约内存角度
// 2.数据和行为解耦的角度，提供一个Setting函数集，通过SetUserData，保存调用者私有变量
// 通过GetUserData拿需要的变量，还是比较爽的，不需要写一堆闭包
func (p *Parser) SetUserData(d interface{}) {
	p.userData = d
}

// GetUserData 获取SetuserData函数设置的私有变量
func (p *Parser) GetUserData() interface{} {
	return p.userData
}

// Execute传递setting参数, 该API 设计成现有形式有如下原因:
// setting如果通过New函数传递, Parser内存占用会多8 * 8的byte
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

// Execute 执行解析器
func (p *Parser) Execute(setting *Setting, buf []byte) (success int, err error) {
	currState := p.currState

	chunkDataStartIndex := 0
	urlStartIndex := 0
	reasonPhraseIndex := unused

	i := 0
	c := byte(0)

	if len(buf) == 0 {
		switch currState {
		case bodyIdentityEOF:
			p.complete(setting, i)
			return 0, nil
		default:
			return 0, nil
		}
	}

	for ; i < len(buf); i++ {
		c = buf[i]

		// fmt.Printf("---->debug state(%s):(%s)method(%#v)\n", currState, buf[i:], p.Method)
	reExec:
		switch currState {
		case startReqOrRsp:
			if c == '\r' || c == '\n' {
				continue
			}

			if c == 'H' {
				if setting.MessageBegin != nil {
					setting.MessageBegin(p, i)
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

			pos := bytes.Index(buf[i:], bytesSpace)
			if pos == -1 {
				p.currState = startReq
				return i, nil
			}

			if setting.MessageBegin != nil {
				setting.MessageBegin(p, i)
			}

			buf2 := buf[i : i+pos]
			switch {
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "GET"):
				p.Method = GET
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "HEAD"):
				p.Method = HEAD
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "POST"):
				p.Method = POST
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "PUT"):
				p.Method = PUT
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "DELETE"):
				p.Method = DELETE
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "CONNECT"):
				p.Method = CONNECT
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "OPTIONS"):
				p.Method = OPTIONS
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "TRACE"):
				p.Method = TRACE
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "ACL"):
				p.Method = ACL
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "BIND"):
				p.Method = BIND
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "COPY"):
				p.Method = COPY
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "CHECKOUT"):
				p.Method = CHECKOUT
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "LOCK"):
				p.Method = LOCK
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "UNLOCK"):
				p.Method = UNLOCK
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "LINK"):
				p.Method = LINK
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "MKCOL"):
				p.Method = MKCOL
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "MOVE"):
				p.Method = MOVE
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "MKACTIVITY"):
				p.Method = MKACTIVITY
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "MERGE"):
				p.Method = MERGE
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "M-SEARCH"):
				p.Method = MSEARCH
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "MKCALENDAR"):
				p.Method = MKCALENDAR
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "NOTIFY"):
				p.Method = NOTIFY
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "PROPFIND"):
				p.Method = PROPFIND
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "PROPPATCH"):
				p.Method = PROPPATCH
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "PATCH"):
				p.Method = PATCH
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "PURGE"):
				p.Method = PURGE
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "REPORT"):
				p.Method = REPORT
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "REBIND"):
				p.Method = REBIND
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "SUBSCRIBE"):
				p.Method = SUBSCRIBE
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "SEARCH"):
				p.Method = SEARCH
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "SOURCE"):
				p.Method = SOURCE
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "UNSUBSCRIBE"):
				p.Method = UNSUBSCRIBE
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "UNBIND"):
				p.Method = UNBIND
			case strings.EqualFold(*(*string)(unsafe.Pointer(&buf2)), "UNLINK"):
				p.Method = UNLINK
			default:
				return 0, fmt.Errorf("%w:%s", ErrMethod, buf2)
			}

			i += pos
			currState = reqMethodAfterSP

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
					setting.URL(p, buf[urlStartIndex:i], i)
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
				return 0, ErrRequestLineLF
			}

			currState = headerField

		case startRsp:
			if c != 'H' {
				return 0, ErrStatusLineHTTP
			}

			if setting.MessageBegin != nil {
				setting.MessageBegin(p, i)
			}

			currState = rspHTTP

		case rspHTTP:
			if len(buf[i:]) < len(strTTPslash) {
				p.currState = currState
				return i, nil
			}

			if !bytes.Equal(buf[i:i+len(strTTPslash)], strTTPslash) {
				return 0, ErrStatusLineHTTP
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
			if c == ' ' {
				continue
			}

			currState = rspStatus
			goto reExec
		case rspStatus:
			if reasonPhraseIndex == unused {
				reasonPhraseIndex = i
			}

			if c == '\r' {
				if setting.Status != nil {
					setting.Status(p, buf[reasonPhraseIndex:i], i)
				}
				currState = rspStatusAfterSP
				continue
			}

			if c == '\n' {
				if setting.Status != nil {
					setting.Status(p, buf[reasonPhraseIndex:i], i)
				}
				currState = headerField
			}

		case rspStatusAfterSP:
			currState = headerField

		case headerField:
			if c == '\r' {
				currState = headersDone
				continue
			}

			// 如果http包只使用'\n'作为分隔符号, 将会进入到这个if里面
			if c == '\n' {
				currState = headersDone
				goto reExec
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
				setting.HeaderField(p, field, i+pos)
			}

			field = bytes.TrimRight(field, " ")
			c2 := c | 0x20
			if c2 == 'c' || c2 == 't' {
				if bytes.EqualFold(field, bytesContentLength) {
					// Content-Length
					p.headerCurrState = hContentLength
					p.contentLength = 0
				} else if bytes.EqualFold(field, bytesTransferEncoding) {
					// Transfer-Encoding
					p.headerCurrState = hTransferEncoding
				} else if bytes.EqualFold(field, bytesConnection) {
					// Connection
					p.headerCurrState = hConnection
				} else {
					// general
					p.headerCurrState = hGeneral
				}
			} else if bytes.EqualFold(field, bytesUpgrade) {
				p.hasUpgrade = true
				p.headerCurrState = hGeneral
			} else {
				p.headerCurrState = hGeneral
			}

			i += pos
			currState = headerValueDiscardWs
		case headerValueDiscardWs:
			// 只跳过一个' ' or '\t'
			// 下个状态可能会跳出, 所以这里先把状态刷到parser里面
			p.currState, currState = headerValue, headerValue
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
				return i, nil
			}

			hValue := buf[i : i+end]
			if setting.HeaderValue != nil {
				setting.HeaderValue(p, hValue, i+end)
			}

			err2 := Split(hValue, bytesCommaSep, func(hValue []byte) error {

				hValue = bytes.TrimSpace(hValue)
				switch p.headerCurrState {
				case hConnection:
					switch {
					case bytes.Contains(hValue, bytesClose):
						p.hasConnectionClose = true
					case bytes.EqualFold(hValue, bytesUpgrade):
						p.hasConnectionUpgrade = true
					}
				case hContentLength:
					n, err := strconv.Atoi(BytesToString(bytes.TrimSpace(hValue)))
					if err != nil {
						return err
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
				return nil
			})
			if err2 != nil {
				return i, err2
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

			// 不是'\n'就是headerField的数据
			goto reExec

		case headersDone:
			if c != '\n' {
				return i, ErrNoEndLF
			}

			if p.hasUpgrade && p.hasConnectionUpgrade {
				p.Upgrade = p.hType == REQUEST || p.StatusCode == 101
			} else {
				//ReadyUpgradeData 函数需要使用
				p.Upgrade = p.Method == CONNECT
			}

			hasBody := p.hasTransferEncoding || p.hasContentLength && p.contentLength != unused

			//fmt.Printf("p.Upgrade:%t, hasBody:%t, hasTrailing:%t\n", p.Upgrade, hasBody, p.hasTrailing)
			if p.Upgrade && !hasBody || p.Method == CONNECT {
				p.complete(setting, i)

				p.currState = p.newMessage()
				return i + 1, nil
			}

			if p.hasTrailing {
				p.complete(setting, i)

				currState = messageDone
				goto reExec
			}

			if setting.HeadersComplete != nil {
				setting.HeadersComplete(p, i)
			}

			// TODO hasContentLength, hasTransferEncoding同时为true
			if p.hasContentLength {
				// 如果contentLength 等于0，说明body的内容为空，可以直接退出
				if p.contentLength == 0 {
					p.complete(setting, i)
					return i, nil
				}
				currState = httpBody
				continue
			}

			if p.hasTransferEncoding {
				currState = chunkedSizeStart
				continue
			}

			if p.EOF() {
				currState = messageDone

				p.complete(setting, i)
				continue
			}
			//一直读到socket eof
			currState = bodyIdentityEOF
		case httpBody:
			if p.hasContentLength {
				nread := min(int32(len(buf[i:])), p.contentLength)
				if setting.Body != nil && nread > 0 {
					setting.Body(p, buf[i:int32(i)+nread], i+int(nread))
				}

				p.contentLength -= nread

				if p.contentLength == 0 {
					currState = messageDone
				}

				i += int(nread) - 1

				// httpBody没有后继结点,所以这里发现数据消费完, 调用下MessageComplete方法
				if p.contentLength == 0 {
					p.complete(setting, i)
				}
			}

		case bodyIdentityEOF:
			if setting.Body != nil {
				setting.Body(p, buf[i:], i+len(buf[i:]))
				i = len(buf) - 1
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
				if c == ';' || c == ' ' {
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

				// 不管有没有trailing数据包, 先当它有
				p.hasTrailing = true
				currState = headerField

				continue
			}

			chunkDataStartIndex = i + 1
			currState = chunkedData

		case chunkedData:
			nread := min(int32(len(buf[i:])), p.contentLength)
			if setting.Body != nil && nread > 0 {
				setting.Body(p, buf[chunkDataStartIndex:int32(chunkDataStartIndex)+nread], chunkDataStartIndex+int(nread))
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
			if p.Upgrade {
				return i, nil
			}

			currState = newState(p.hType)
			p.Reset()
			goto reExec
		}

	}

	switch currState {
	case reqURL:
		if setting.URL != nil && len(buf[urlStartIndex:]) > 0 {
			setting.URL(p, buf[urlStartIndex:], len(buf))
		}

	case rspStatus:
		if setting.Status != nil && len(buf[reasonPhraseIndex:]) > 0 {
			setting.Status(p, buf[reasonPhraseIndex:], len(buf))
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

// Reset 重置状态
func (p *Parser) Reset() {
	p.currState = newState(p.hType)
	p.headerCurrState = hGeneral
	p.Major = 0
	p.Minor = 0
	//p.MaxHeaderSize
	p.contentLength = unused
	p.StatusCode = 0
	p.hasContentLength = false
	p.hasTransferEncoding = false
	p.hasConnectionClose = false
	p.hasUpgrade = false
	p.hasConnectionUpgrade = false
	p.hasTrailing = false
	p.callMessageComplete = false
	p.Upgrade = false
}

// Status debug专用
func (p *Parser) Status() string {
	return stateTab[p.currState]
}

// EOF 表示结束
func (p *Parser) EOF() bool {
	if p.hType == REQUEST {
		return true
	}

	return p.currState == messageDone
}

func (p *Parser) shouldKeepAlive() bool {
	if p.Major > 0 && p.Minor > 0 {
		if p.hasConnectionClose {
			return false
		}
		//TODO
	}

	return p.EOF()
}

func (p *Parser) newMessage() state {
	if p.shouldKeepAlive() {
		if p.hType == REQUEST {
			return startReq
		}
		return startRsp
	}

	return dead
}

func min(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}
