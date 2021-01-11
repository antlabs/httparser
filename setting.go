package httparser

type Setting struct {
	// 解析开始
	MessageBegin func(*Parser)
	// url 回调函数
	URL func(*Parser, []byte)
	// http field 回调函数
	HeaderField func(*Parser, []byte)
	// http value 回调函数
	HeaderValue func(*Parser, []byte)
	// http 解析完成之后的回调函数
	HeadersComplete func(*Parser)
	// body结束
	MessageEnd func(*Parser)
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
	// request or response状态，这里让解析器自己选择
	startReqOrRsp
)
