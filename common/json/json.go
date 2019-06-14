package json

import "encoding/json"

// RequestMessage ...
type RequestMessage struct {
	Type string  `json:"type"`
	Data Request `json:"data"`
	CRID string  `json:"crid"`
}

// ResponseMessage ...
type ResponseMessage struct {
	Type       string          `json:"type"`
	RawData    json.RawMessage `json:"data"`
	Data       Response        `json:"-"`
	CRID       string          `json:"crid"`
	ResultCode byte            `json:"result"`
}

// Request sent to server and will receive a response with same crid
type Request interface{}

// Response received from the server following the request with the associated crid, result code = 0 -> OK
type Response interface{}

// Notification push message from the server
type Notification interface{}
