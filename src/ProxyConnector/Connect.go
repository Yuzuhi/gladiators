package ProxyConnector

import (
	"encoding/binary"
	"fmt"
	"gladiators/src/MyError"
	"gladiators/src/ProxyConnector/Messages"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type ProxyConnector struct {
	addr           string
	connectionType string
	conn           net.Conn
	wg             *sync.WaitGroup
}

func NewProxyConnector(connectionType, addr string) *ProxyConnector {
	return &ProxyConnector{
		addr:           addr,
		connectionType: connectionType,
		wg:             &sync.WaitGroup{},
	}
}

func (c *ProxyConnector) sendMessage(message messages.ProxyClientMsg) error {

	if c.conn == nil {
		return MyError.ErrConnectionClosed
	}

	messageBytes,err := message.Serialize()

	if err != nil{
		return err
	}

	if _, err := c.conn.Write(messageBytes); err != nil {
		return err
	}

	return nil
}

func (c *ProxyConnector) readMessage() (messages.MessageType, []byte, error) {
	// read header
	header := make([]byte, 1+4)

	if	_, err := io.ReadFull(c.conn, header);err != nil{
		return 0, []byte{}, err
	}

	// 获取message类型
	messageType := messages.MessageType(header[0])

	if !messages.IsValidMessageType(messageType) {
		return messageType,[]byte{}, MyError.ErrUndefinedMessageType(string(messageType))
	}

	// 获取jsonbytes长度
	contentLength := int(binary.BigEndian.Uint32(header[1:]))

	jsonBytes := make([]byte, contentLength)
	// 读取jsonbytes
	if _, err := io.ReadFull(c.conn, jsonBytes);err != nil{
		return 0, []byte{}, err
	}

	return messageType, jsonBytes, nil

}


func (ps *ProxyConnector) handleConnection(readDoneSignal, sendDoneSignal chan bool, internalMessageChan chan messages.ProxyClientMsg) {

	heartBeatMsg, err := messages.CreateNewMessage(messages.HeartBeatMsgType)

	if err != nil {
		ps.conn.Close()
		panic(err)
		return
	}

	// create a goroutine to send message

	ps.wg.Add(1)

	go func(readDoneSignal, sendDoneSignal chan bool, internalMessageChan chan messages.ProxyClientMsg) {

		//  Set the interval to send heartbeat packets
		heartbeatTicker := time.NewTicker(5 * time.Second)

		defer func() {
			ps.conn.Close()
			heartbeatTicker.Stop()
			ps.wg.Done()
		}()

		for {
			select {
			case <-sendDoneSignal:
				return
			case <-heartbeatTicker.C:
				if err := ps.sendMessage(heartBeatMsg); err != nil {
					log.Print(err)
					close(readDoneSignal)
					return
				}
			case msg := <-internalMessageChan:
				fmt.Println("准备转发：", msg.GetStringifyData())
				if !messages.IsValidMessageType(msg.GetMessageType()) {
					log.Panic(MyError.ErrUndefinedMessageType(string(msg.GetMessageType())))
					close(readDoneSignal)
					return
				}
				if err := ps.sendMessage(msg); err != nil {
					log.Print(err)
					close(readDoneSignal)
					return
				}
				fmt.Println("转发请求：", msg.GetStringifyData())

			}

		}

	}(readDoneSignal, sendDoneSignal, internalMessageChan)

	// create a goroutine to receive message from connection

	ps.wg.Add(1)

	go func(readDoneSignal, sendDoneSignal chan bool, internalMessageChan chan messages.ProxyClientMsg) {

		defer ps.wg.Done()

		for {
			select {
			case <-readDoneSignal:
				return
			default:
				// read message from proxy
				messageType, content, err := ps.readMessage()
				if err != nil {
					close(sendDoneSignal)
					log.Print(err)
					return
				}

				switch messageType {

				case messages.HeartBeatMsgType:
					fmt.Printf("received:%s\n", content)

				case messages.ProxyResponseMsgType:
					fmt.Printf("get normal msg:%s\n", content)
				default:
					err := MyError.ErrUndefinedMessageType(string(messageType))
					close(sendDoneSignal)
					log.Print(err)
					return

				}
			}
		}
	}(readDoneSignal, sendDoneSignal, internalMessageChan)

}

func (ps *ProxyConnector) HandleProxyConnection(internalMessageChan chan messages.ProxyClientMsg) error {
	var err error
	var count int

	for {

		if count >= 3 {
			return err
		}

		ps.conn, err = net.Dial(ps.connectionType, ps.addr)

		if err != nil {
			fmt.Printf("执行第%d次重连\n", count+1)
			count++
			time.Sleep(10 * time.Second)
			continue
		}

		count = 0

		readDoneSignal := make(chan bool)
		sendDoneSignal := make(chan bool)

		ps.handleConnection(readDoneSignal, sendDoneSignal, internalMessageChan)

		ps.wg.Wait()

	}

}
