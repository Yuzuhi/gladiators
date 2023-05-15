package messages

import "net/http"

type ProxyClientMsg interface {
	Serialize() ([]byte, error)
	GetStringifyData() string
	GetMessageType() MessageType
	GetStringifyType() string
}

type ProxyClientWithHTTPRequestMsg interface {
	GetHTTPRequest() *http.Request
}
