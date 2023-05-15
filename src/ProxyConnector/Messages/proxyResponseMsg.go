package messages

import "net/http"

type ProxyResponseMsg struct {
	messageType MessageType
	Rep   *http.Response
}

func (msg *ProxyResponseMsg) GetMessageType() MessageType {
	return msg.messageType
}

func (msg *ProxyResponseMsg) GetStringifyData() string {
	return "this is response message from proxyserver"
}

func (msg *ProxyResponseMsg) Serialize()([]byte,error)  {
	
	return serialize(msg)
}

func (msg *ProxyResponseMsg) GetStringifyType() string {
	return "ProxyResponseMsg"
}

