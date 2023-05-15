package messages

import (
	"fmt"
)



type ClientRequestMsg struct {
	MessageType MessageType
	RequsetData RequsetData
}

func (c *ClientRequestMsg) Serialize() ([]byte, error) {
	return serialize(c)
}

func (c *ClientRequestMsg) GetStringifyData() string {

	msg := fmt.Sprintf("target url:%s",c.RequsetData.URL)

	return msg
}

func (c *ClientRequestMsg) GetMessageType() MessageType {
	return c.MessageType
}

func (c *ClientRequestMsg) GetStringifyType() string {
	return "ClientRequestMsg"
}
