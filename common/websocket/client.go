// The wesocket package offers a websocket based client that will connect and communicate with the server is a request/response protocol
package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	api "republicofminer-client-go/common/json"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	requests  chan []byte
	crids     map[string]chan *api.ResponseMessage
	id        uint32
	seed      int64
	increment int
	factory   func(t string) (api.Response, bool)
}

var unique = uint32(0)

// instanciates a WebSocketClient that will use the factory function to instanciate the response for a given type
func Client(factory func(t string) (api.Response, bool)) *WebSocketClient {
	client := &WebSocketClient{}
	client.requests = make(chan []byte)
	client.crids = make(map[string]chan *api.ResponseMessage)
	client.seed = time.Now().Unix()
	client.increment = 0
	client.factory = factory
	atomic.AddUint32(&client.id, uint32(1))
	return client
}

func send(client *WebSocketClient, bytes []byte) {
	client.requests <- bytes
}

func receive(client *WebSocketClient, bytes []byte) {
	// we parse the header
	var response api.ResponseMessage
	err := json.Unmarshal(bytes, &response)
	if err != nil {
		log.Println("Unmarshal message header error:", err)
		return
	}

	// we get the data
	data, success := client.factory(response.Type)
	if !success {
		log.Println("Unknow response type :", response.Type)
	} else {
		err = json.Unmarshal(response.RawData, data)
		if err != nil {
			log.Println("Unmarshal message data error:", err)
			return
		}
		response.Data = data
	}

	// TODO
	// it is a notification
	if response.CRID == "" {
		return
	}

	// we notify the requester channel
	channel, ok := client.crids[response.CRID]
	if !ok {
		log.Println("unknown CRID :", response.CRID)
	} else {
		channel <- &response
	}
}

// Connect to the server and blocks the thread while the connection is opened
func (client *WebSocketClient) Connect(uri string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: uri, Path: ""}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("error on read:", err)
				return
			}
			log.Printf("recv: %s", message)
			receive(client, message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case request := <-client.requests:
			err := c.WriteMessage(websocket.TextMessage, request)
			if err != nil {
				log.Println("error on write:", err)
				return
			}
			log.Println("write:", string(request))
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// TODO make thread safe
func (client *WebSocketClient) crid() string {
	crid := fmt.Sprintf("%d-%d-%d", client.seed, client.id, client.increment)
	client.increment++
	return crid
}

// Request is serialized and sent to the server
func (client *WebSocketClient) Request(request *api.RequestMessage) chan *api.ResponseMessage {
	channel := make(chan *api.ResponseMessage)
	marshaled, _ := json.Marshal(request)
	send(client, marshaled)
	client.crids[request.CRID] = channel
	return channel
}

// RequestMessage is a wrapper for request
func (client *WebSocketClient) RequestMessage(request api.Request, t string) *api.RequestMessage {
	return &api.RequestMessage{CRID: client.crid(), Type: t, Data: request}
}
