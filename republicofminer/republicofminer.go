package republicofminer

import (
	"republicofminer-client-go/common/websocket"
	"republicofminer-client-go/republicofminer/api"
)

var client *websocket.WebSocketClient

func Connect() {
	client = websocket.Client(api.CreateResponse)
	client.Connect("game.republicofminer.com:2026")
}

func GetMiningTask(address string, resource string) *api.MiningTask {
	request := api.GetMiningTaskRequest{Address: address, Resource: resource}
	response := <-client.Request(client.RequestMessage(&request, "GetMiningTaskRequest"))
	data := response.Data.(*api.GetMiningTaskResponse)
	return data.Task
}
