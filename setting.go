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

// 查阅#6 看设计变更原因
type Setting struct {
	// 解析开始
	MessageBegin func(*Parser)
	// url 回调函数, 只有在request包才会回调
	// 解析一个包时,URL回调可能会多次调用
	URL func(*Parser, []byte)
	// 状态短语
	// 解析一个包时, Status回调可能会多次调用
	Status func(*Parser, []byte)
	// http field 回调函数
	HeaderField func(*Parser, []byte)
	// http value 回调函数
	HeaderValue func(*Parser, []byte)
	// http 解析完成之后的回调函数
	HeadersComplete func(*Parser)
	// body的回调函数
	Body func(*Parser, []byte)
	// 所有消息成功解析
	MessageComplete func(*Parser)
}

type ReqOrRsp uint8

const (
	REQUEST ReqOrRsp = iota + 1
	RESPONSE
	BOTH
)

type state uint8

func (s state) String() string {
	return stateTab[s]
}

const (
	// request状态
	startReq state = iota + 1
	// reqMethod状态
	reqMethod
	// reqMethod 后面的SP
	reqMethodAfterSP
	// 请求URL
	reqURL
	// 请求URL后面的SP
	reqURLAfterSP
	//
	reqHTTPVersion
	// HTTP-Version中的major
	reqHTTPVersionMajor
	// HTTP-Version中的.
	reqHTTPVersionDot
	// HTTP-Version中的minor
	reqHTTPVersionMinor

	// request-line \r的位置
	reqRequestLineAlomstDone
	// response状态
	startRsp
	// HTTP
	rspHTTP
	// response 版本号数字
	rspHTTPVersionNum
	// response 版本号后面的空格
	rspHTTPVersionNumAfterSP
	// response 状态吗
	rspStatusCode
	// statuscode后面的sp符号
	rspStatusCodeAfterSP
	// response状态短语
	rspStatus
	// 状态短语后面的SP符号
	rspStatusAfterSP
	// request or response状态，这里让解析器自己选择
	startReqOrRsp

	// http header解析结束
	headersDone
	// 解析http field
	headerField
	// 进入http header分隔符号
	headerValueDiscardWs
	// 进入http value
	headerValue
	// 刚开始进入http value后面的OWS
	headerValueStartOWS
	// 快要离开http value后面的OWS
	headerValueOWS
	// 进入http body
	httpBody
	// 开始进入到chunked 数字解析
	chunkedSizeStart
	// 进入到chunked 数字
	chunkedSize
	// chunked size结束
	chunkedSizeAlmostDone
	// chunked ext
	chunkedExt
	// chunked data
	chunkedData
	// chunked 检查是否真的结束
	chunkedDataAlmostDone
	// chunked data结束
	chunkedDataDone
	// 快要结束
	messageAlmostDone
	// 一直读到socket eof
	bodyIdentityEof
	// 解析结束
	messageDone
)

// debug使用
var stateTab = []string{
	startReq:                 "startReq",
	reqMethod:                "reqMethod",
	reqMethodAfterSP:         "reqMethodAfterSP",
	reqURL:                   "reqURL",
	reqURLAfterSP:            "reqURLAfterSP",
	reqHTTPVersion:           "reqHTTPVersion",
	reqHTTPVersionMajor:      "reqHTTPVersionMajor",
	reqHTTPVersionDot:        "reqHTTPVersionDot",
	reqHTTPVersionMinor:      "reqHTTPVersionMinor",
	reqRequestLineAlomstDone: "reqRequestLineAlomstDone",
	startRsp:                 "startRsp",
	rspHTTP:                  "rspHTTP",
	rspHTTPVersionNum:        "rspHTTPVersionNum",
	rspHTTPVersionNumAfterSP: "rspHTTPVersionNumAfterSP",
	rspStatusCode:            "rspStatusCode",
	rspStatusCodeAfterSP:     "rspStatusCodeAfterSP",
	rspStatus:                "rspStatus",
	rspStatusAfterSP:         "rspStatusCodeAfterSP",
	startReqOrRsp:            "startReqOrRsp",
	headersDone:              "headersDone",
	headerField:              "headerField",
	headerValueDiscardWs:     "headerValueDiscardWs",
	headerValue:              "headerValue",
	headerValueStartOWS:      "headerValueStartOWS",
	headerValueOWS:           "headerValueOWS",
	httpBody:                 "httpBody",
	chunkedSizeStart:         "chunkedSizeStart",
	chunkedSize:              "chunkedSize",
	chunkedSizeAlmostDone:    "chunkedSizeAlmostDone",
	chunkedExt:               "chunkedExt",
	chunkedData:              "chunkedData",
	chunkedDataAlmostDone:    "chunkedDataAlmostDone",
	chunkedDataDone:          "chunkedDataDone",
	messageAlmostDone:        "messageAlmostDone",
	bodyIdentityEof:          "bodyIdentityEof",
	messageDone:              "messageDone",
}

type headerState uint8

const (
	hGeneral headerState = iota
	hContentLength
	hTransferEncoding
	hConnection
)
