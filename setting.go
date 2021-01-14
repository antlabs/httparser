package httparser

type Setting struct {
	// 解析开始
	MessageBegin func()
	// url 回调函数
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
	// body结束
	MessageEnd func()
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
	// response状态
	startRsp
	rspHTTP
	rspHTTPVersionNum
	rspStatusCode
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
	// 进入http body
	httpBody
)

type headerState uint8

const (
	hGeneral headerState = iota
	hContentLength
)
