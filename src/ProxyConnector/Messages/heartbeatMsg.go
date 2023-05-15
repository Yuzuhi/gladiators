package messages

type HeartBeatMsg struct {
	MessageType MessageType
}

func (msg *HeartBeatMsg) GetMessageType() MessageType {
	return msg.MessageType
}

func (msg *HeartBeatMsg) GetStringifyData() string {
	return "heartbeat"
}

func (msg *HeartBeatMsg) Serialize() ([]byte,error) {
	return serialize(msg)
}

func (msg *HeartBeatMsg) GetStringifyType() string {
	return "HeartBeatMsg"
}
