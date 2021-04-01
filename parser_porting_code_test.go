package httparser

import "testing"

type message struct {
	name  string
	raw   string
	hType ReqOrRsp
	//httpMethod      method
	statusCode      int
	response_status string
	request_path    string
	request_url     string
	fragment        string
	queryString     string
	body            string
	bodySize        string
	host            string
	userinfo        string
	port            uint16
	numHeaders      int
	//enum { NONE=0, FIELD, VALUE } last_header_element; //TODO最近移值
	headers         [][]string
	shouldKeepAlive int

	numChunks         int
	numChunksComplete int
	chunkLengths      []int

	upgrade string // upgraded body

	httpMajor     uint16
	httpMinor     uint16
	contentLength uint64

	messageBeginCbCalled    int
	headersCompleteCbCalled int
	messageCompleteCbCalled int
	statusCbCalled          int
	messageCompleteOnEof    int
	bodyIsFinal             int
	allowChunkedLength      bool
}

func test_Message(t *testing.T, m *message) {
	for msg1len := 0; msg1len < len(m.raw); msg1len++ {
		//p := New(m.hType)
	}
}
