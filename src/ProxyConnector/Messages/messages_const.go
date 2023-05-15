package messages

import (
	"net/http"
	"net/url"
)

type MessageType uint8

const (
	HeartBeatMsgType MessageType = iota + 1
	ClientRequestType
	ProxyResponseMsgType
)

var MessageTypes = []MessageType{
	HeartBeatMsgType,
	ClientRequestType,
	ProxyResponseMsgType,
}

type RequsetData struct{
	URL         string         `json:"url"`
	Method      string         `json:"method"`
	Header      http.Header    `json:"header"`
	Cookies     []*http.Cookie `json:"cookies"`
	Body        string         `json:"body"`
	FormValues  url.Values     `json:"formValues"`
}