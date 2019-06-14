package api

import api "republicofminer-client-go/common/json"

type GetMiningTaskRequest struct {
	Address  string
	Resource string
}

type GetMiningTaskResponse struct {
	Task *MiningTask
}

type MiningTask struct {
	Address    string
	SecretHash string
	Mask       string
	Currency   string
	Amount     float64
}

type ClaimMiningRequest struct {
	TaskAddress string
	Secret      string
	Receiver    string
}

type ClaimMiningResponse struct {
	TransactionHash string
}

func CreateResponse(t string) (api.Response, bool) {
	switch t {
	case "GetMiningTaskResponse":
		return &GetMiningTaskResponse{}, true
	case "ClaimMiningResponse":
		return &ClaimMiningResponse{}, true
	}
	return nil, false
}
