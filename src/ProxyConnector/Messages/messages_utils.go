package messages

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gladiators/src/MyError"
	"net/http"
	"net/url"
)

func IsValidMessageType(mt MessageType) bool {
	numTypes := len(MessageTypes)

	for i := 0; i < numTypes; i++ {
		if mt == MessageTypes[i] {
			return true
		}

	}

	return false
}

func serialize(message ProxyClientMsg) ([]byte, error) {
	// Serialize the message to JSON
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	// Calculate the length of the serialized message
	messageLength := uint32(len(jsonBytes))

	// Create a 5-byte header, type + messageLength
	header := make([]byte, 1+4)

	header[0] = byte(message.GetMessageType())

	// Store the message length in the header using big-endian byte order
	binary.BigEndian.PutUint32(header[1:], messageLength)

	// Combine the header and the serialized message into a single byte slice
	bytesMessage := append(header, jsonBytes...)

	return bytesMessage, nil
}

func CreateNewMessage(messageType MessageType, data ...interface{}) (ProxyClientMsg, error) {

	switch messageType {
	case HeartBeatMsgType:
		if len(data) != 0 {
			err := fmt.Errorf("do not pass second parameter if you are going to creat HeartBeatType message")
			return nil, err
		}
		return &HeartBeatMsg{
			MessageType: HeartBeatMsgType,
		}, nil

	case ClientRequestType:
		if len(data) == 0 {
			err := fmt.Errorf("no request data")
			return nil, err
		}

		req, ok := data[0].(*http.Request)

		if !ok {
			err := MyError.ErrInvalidRequsetType
			return nil, err
		}

		return BuildRequestData(req)

	case ProxyResponseMsgType:
		if len(data) == 0 {
			err := fmt.Errorf("no response data")
			return nil, err
		}

		rep, ok := data[0].(*http.Response)

		if !ok {
			err := MyError.ErrInvalidRequsetType
			return nil, err
		}

		return &ProxyResponseMsg{
			messageType: ProxyResponseMsgType,
			Rep:         rep,
		}, nil
	}

	err := MyError.ErrUndefinedMessageType(string(messageType))
	return nil, err
}

func BuildRequestData(req *http.Request) (*ClientRequestMsg, error) {

	clientReq := &ClientRequestMsg{
		MessageType: ClientRequestType,
		RequsetData: RequsetData{
			URL:        req.URL.String(),
			Method:     req.Method,
			Header:     req.Header,
			Cookies:    make([]*http.Cookie, len(req.Cookies())),
			FormValues: url.Values{},
		},
	}

	// 复制 Cookies
	copy(clientReq.RequsetData.Cookies, req.Cookies())

	// 处理 Body
	if req.Method == "POST" && req.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := req.ParseForm()
		if err != nil {
			return clientReq, err
		}
		clientReq.RequsetData.FormValues = req.PostForm
	}

	return clientReq, nil

}
