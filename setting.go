package httparser

type Setting struct {
	// 解析开始
	MessageBegin func()
	// url 回调函数, 只有在request包才会回调
	URL func([]byte)
	// 状态短语
	Status func([]byte)
	// http field 回调函数
	HeaderField func([]byte)
	// http value 回调函数
	HeaderValue func([]byte)
	// http 解析完成之后的回调函数
	HeadersComplete func()
	// body的回调函数
	Body func([]byte)
	// 所有消息成功解析
	MessageComplete func()
}

type ReqOrRsp uint8

const (
	REQUEST ReqOrRsp = iota + 1
	RESPONSE
	BOTH
)

type state uint8

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
	// request-line \r的位置
	reqRequestLineAlomstDone
	// response状态
	startRsp
	// HTTP
	rspHTTP
	// response 版本号数字
	rspHTTPVersionNum
	// response 状态吗
	rspStatusCode
	// response状态短语
	rspStatus
	// request or response状态，这里让解析器自己选择
	startReqOrRsp

	// http header解析结束
	headerDone
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
	reqRequestLineAlomstDone: "reqRequestLineAlomstDone",
	startRsp:                 "startRsp",
	rspHTTP:                  "rspHTTP",
	rspHTTPVersionNum:        "rspHTTPVersionNum",
	rspStatusCode:            "rspStatusCode",
	rspStatus:                "rspStatus",
	startReqOrRsp:            "startReqOrRsp",
	headerDone:               "headerDone",
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
	messageDone:              "messageDone",
}

type headerState uint8

const (
	hGeneral headerState = iota
	hContentLength
	hTransferEncoding
)
