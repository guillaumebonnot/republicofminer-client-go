package explorer

import (
	"republicofminer-client-go/common/websocket"
	"republicofminer-client-go/explorer/api"
)

var client *websocket.WebSocketClient

func Connect() {
	client = websocket.Client(api.CreateResponse)
	client.Connect("data.republicofminer.com:2030")
}

func GetTransaction(hash string) *api.Transaction {
	request := api.GetTransactionRequest{Hash: hash}
	response := <-client.Request(client.RequestMessage(&request, "GetTransactionRequest"))
	data := response.Data.(*api.GetTransactionResponse)
	return &data.Transaction
}

func GetLedgerByHash(hash string) *api.Ledger {
	return GetLedger(&api.GetLedgerRequest{Hash: hash})
}

func GetLedgerByHeight(height int64) *api.Ledger {
	return GetLedger(&api.GetLedgerRequest{Height: &height})
}

func GetLedger(request *api.GetLedgerRequest) *api.Ledger {
	response := <-client.Request(client.RequestMessage(request, "GetLedgerRequest"))
	data := response.Data.(*api.GetLedgerResponse)
	return &data.Ledger
}

func SendTransaction(transaction *api.Transaction, signatures []*api.Signature) string {
	request := api.SendTransactionRequest{Transaction: transaction, Signatures: signatures}
	response := <-client.Request(client.RequestMessage(&request, "SendTransactionRequest"))
	data := response.Data.(*api.SendTransactionResponse)
	return data.Hash
}

func GetAccount(encoded string) map[string]float64 {
	request := api.GetAccountRequest{Address: encoded}
	response := <-client.Request(client.RequestMessage(&request, "GetAccountRequest"))
	data := response.Data.(*api.GetAccountResponse)
	return data.Balance
}
