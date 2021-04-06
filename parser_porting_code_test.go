package httparser

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type message struct {
	name  string
	hType ReqOrRsp
	raw   string
	//httpMethod      method
	statusCode     int
	responseStatus string
	requestPath    string
	requestUrl     string
	body           string
	bodySize       string
	host           string
	userinfo       string
	port           uint16
	numHeaders     int
	//enum { NONE=0, FIELD, VALUE } last_header_element; //TODO最近移值
	headers         [][]string
	shouldKeepAlive bool

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
	messageCompleteOnEof    bool
	bodyIsFinal             int
	allowChunkedLength      bool
}

func (m *message) eq(m2 *message, t *testing.T) {
	b := assert.Equal(t, m.headers, m2.headers)
	if !b {
		return
	}

	assert.Equal(t, m.httpMajor, m2.httpMajor)
	assert.Equal(t, m.httpMinor, m2.httpMinor)
	assert.Equal(t, m.hType, m2.hType)
}

var requests = []message{
	{
		name:  "curl get",
		hType: REQUEST,
		raw: "GET /test HTTP/1.1\r\n" +
			"User-Agent: curl/7.18.0 (i486-pc-linux-gnu) libcurl/7.18.0 OpenSSL/0.9.8g zlib/1.2.3.3 libidn/1.1\r\n" +
			"Host: 0.0.0.0=5000\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
		shouldKeepAlive:      true,
		messageCompleteOnEof: false,
		httpMajor:            1,
		httpMinor:            1,
		//method: HTTP_GET,
		requestUrl:    "/test",
		contentLength: math.MaxUint64,
		headers: [][]string{
			{"User-Agent", "curl/7.18.0 (i486-pc-linux-gnu) libcurl/7.18.0 OpenSSL/0.9.8g zlib/1.2.3.3 libidn/1.1"},
			{"Host", "0.0.0.0=5000"},
			{"Accept", "*/*"},
		},
	},
}

func test_Message(t *testing.T, m *message) {
	for msg1len := 0; msg1len < len(m.raw); msg1len++ {
		//p := New(m.hType)
	}
}
