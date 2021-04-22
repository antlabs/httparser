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
	//enum { NONE=0, FIELD, VALUE } last_header_element; //TODO最近移值
	headers         [][2]string
	shouldKeepAlive bool

	numChunks         int
	numChunksComplete int
	chunkLengths      []int

	upgrade string // upgraded body

	httpMajor     uint16
	httpMinor     uint16
	contentLength uint64

	messageBeginCbCalled    bool
	headersCompleteCbCalled bool
	messageCompleteCbCalled bool
	statusCbCalled          bool
	messageCompleteOnEof    bool
	bodyIsFinal             int
	allowChunkedLength      bool
}

func (m *message) eq(t *testing.T, m2 *message) bool {
	b := assert.Equal(t, m.headers, m2.headers)
	if !b {
		return false
	}

	b = assert.Equal(t, m.httpMajor, m2.httpMajor, "major")
	if !b {
		return false
	}
	b = assert.Equal(t, m.httpMinor, m2.httpMinor, "minor")
	if !b {
		return false
	}
	b = assert.Equal(t, m.hType, m2.hType)
	if !b {
		return false
	}
	return true
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
		headers: [][2]string{
			{"User-Agent", "curl/7.18.0 (i486-pc-linux-gnu) libcurl/7.18.0 OpenSSL/0.9.8g zlib/1.2.3.3 libidn/1.1"},
			{"Host", "0.0.0.0=5000"},
			{"Accept", "*/*"},
		},
	},
}

var settingTest Setting = Setting{
	MessageBegin: func(p *Parser) {
		m := p.GetUserData().(*message)
		m.messageBeginCbCalled = true
	},
	URL: func(p *Parser, url []byte) {
		m := p.GetUserData().(*message)
		m.requestUrl = string(url)
	},
	Status: func(p *Parser, status []byte) {
		m := p.GetUserData().(*message)
		m.responseStatus = string(status)
	},
	HeaderField: func(p *Parser, headerField []byte) {
		m := p.GetUserData().(*message)
		m.headers = append(m.headers, [2]string{string(headerField), ""})
	},
	HeaderValue: func(p *Parser, headerValue []byte) {
		m := p.GetUserData().(*message)
		m.headers[len(m.headers)-1][1] = string(headerValue)
	},
	HeadersComplete: func(p *Parser) {
		m := p.GetUserData().(*message)
		m.headersCompleteCbCalled = true
	},
	Body: func(p *Parser, body []byte) {
		m := p.GetUserData().(*message)
		m.body += string(body)
	},
	MessageComplete: func(p *Parser) {
		m := p.GetUserData().(*message)
		m.messageCompleteCbCalled = true
		m.httpMajor = uint16(p.Major)
		m.httpMinor = uint16(p.Minor)
	},
}

func parse(p *Parser, data string) error {
	_, err := p.Execute(&settingTest, []byte(data))
	return err
}

func test_Message(t *testing.T, m *message) {
	for msg1len := 0; msg1len < len(m.raw); msg1len++ {
		p := New(m.hType)
		got := &message{}
		p.SetUserData(got)

		msg1Message := m.raw[:msg1len]
		msg2Message := m.raw[msg1len:]

		if msg1len > 0 {
			err := parse(p, msg1Message)
			assert.NoError(t, err)
		}

		err := parse(p, msg2Message)
		assert.NoError(t, err)
		if b := m.eq(t, got); !b {
			t.Logf("test case name:%s\n", m.name)
			break
		}

	}
}

func Test_Message(t *testing.T) {
	for _, req := range requests {
		//test_Message(t, &req)
		_ = req
	}
}
