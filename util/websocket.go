package util

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

// HTTP upgrade configuration of websocket protocol
var wsUpgrader = websocket.Upgrader{
	// Allow all CORS cross-domain requests
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsMessage websocket message
type WsMessage struct {
	MessageType int
	Data        []byte
}

// WsConnection Package WebSocket connection
type WsConnection struct {
	wsSocket *websocket.Conn // The bottom websocket
	inChan   chan *WsMessage // Read queue
	outChan  chan *WsMessage // Send a queue

	mutex     sync.Mutex // Avoid repeatedly closed the pipeline
	isClosed  bool
	CloseChan chan byte // Close notification
}

// Reading coroutine
func (wsConn *WsConnection) wsReadLoop() {
	var (
		msgType int
		data    []byte
		msg     *WsMessage
		err     error
	)
	for {
		// Read a message
		if msgType, data, err = wsConn.wsSocket.ReadMessage(); err != nil {
			fmt.Println("Read error: " + err.Error())
			goto CLOSED
		}
		msg = &WsMessage{
			msgType,
			data,
		}
		// Put in the request queue
		select {
		case wsConn.inChan <- msg:
			if string(data) == "exit" {
				goto CLOSED
			}
		case <-wsConn.CloseChan:
			goto CLOSED
		}
	}
CLOSED:
	wsConn.WsClose()
}

// Sending coroutine
func (wsConn *WsConnection) wsWriteLoop() {
	var (
		msg *WsMessage
		err error
	)
	for {
		select {
		// Take a response
		case msg = <-wsConn.outChan:
			// Write to WebSocket
			if err = wsConn.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				fmt.Println(err)
				goto CLOSED
			}
		case <-wsConn.CloseChan:
			goto CLOSED
		}
	}
CLOSED:
	wsConn.WsClose()
}

// InitWebsocket Initialization WS
func InitWebsocket(resp http.ResponseWriter, req *http.Request) (wsConn *WsConnection, err error) {
	var (
		wsSocket *websocket.Conn
	)
	// Answer the client to inform the upgrade connection as websocket
	if wsSocket, err = wsUpgrader.Upgrade(resp, req, nil); err != nil {
		return
	}
	wsConn = &WsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *WsMessage, 1000),
		outChan:   make(chan *WsMessage, 1000),
		CloseChan: make(chan byte),
		isClosed:  false,
	}

	// Reading coroutine
	go wsConn.wsReadLoop()
	// Writing corporation
	go wsConn.wsWriteLoop()

	return
}

// WsWrite send messages
func (wsConn *WsConnection) WsWrite(messageType int, data []byte) (err error) {
	select {
	case wsConn.outChan <- &WsMessage{messageType, data}:
	case <-wsConn.CloseChan:
		err = errors.New("websocket closed")
	}
	return
}

// WsRead Read message
func (wsConn *WsConnection) WsRead() (msg *WsMessage, err error) {
	select {
	case msg = <-wsConn.inChan:
		return
	case <-wsConn.CloseChan:
		err = errors.New("websocket closed")
	}
	return
}

// WsClose Close the connection
func (wsConn *WsConnection) WsClose() {
	wsConn.wsSocket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.CloseChan)
	}
}
