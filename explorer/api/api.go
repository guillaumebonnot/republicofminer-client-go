package api

import (
	"encoding/json"
	"errors"
	api "republicofminer-client-go/common/json"
	"republicofminer-client-go/protocol"
)

// Ledger ...
type Ledger struct {
	Height              int64
	Hash                string
	Timestamp           int64
	Lastledger          string
	Version             byte
	FeeTransactionIndex int32
	Transactions        []*TransactionHeader
}

// TransactionHeader ...
type TransactionHeader struct {
	Index          int      `json:"i"`
	Hash           string   `json:"h"`
	Fee            *float64 `json:"f"`
	HasDeclaration bool     `json:"d"`
}

type TxDeclaration struct {
	Type        protocol.DeclarationType
	Declaration interface{}
}

// Transaction ...
type Transaction struct {
	Hash         string
	Expire       *int64
	Declarations []*TxDeclaration
	Inputs       []*TxInput
	Outputs      []*TxOutput
	Message      string   `json:",omitempty"`
	Fees         *TxInput `json:",omitempty"`
}

// TxInputOutput is the base type for TxInput and TxOutput
type TxInputOutput struct {
	Address  string
	Currency string
	Amount   float64
}

type TxInput TxInputOutput
type TxOutput TxInputOutput

type HashLock struct {
	Address    string
	SecretHash SecretHash
}

type SecretHashType byte

const (
	SHA3   SecretHashType = 0
	SHA256 SecretHashType = 1
)

type SecretHash struct {
	Type SecretHashType
	Hash string
}

type SecretRevelation struct {
	Secret string
}

type MultiSignature struct {
	Address  string
	Signers  []string
	Required int32
}

// GetLedgerRequest ...
type GetLedgerRequest struct {
	Height *int64 `json:",omitempty"`
	Hash   string `json:",omitempty"`
}

// GetLedgerResponse ...
type GetLedgerResponse struct {
	Ledger Ledger
}

// GetTransactionResponse ...
type GetTransactionResponse struct {
	Transaction Transaction
}

// GetTransactionRequest ...
type GetTransactionRequest struct {
	Hash string
}

type SendTransactionRequest struct {
	Transaction *Transaction
	Signatures  []*Signature
}

type Signature struct {
	PublicKey     string `json:"k"`
	SignatureByte string `json:"s"`
}

type SendTransactionResponse struct {
	Hash string
}

type GetAccountRequest struct {
	Address string
}

type GetAccountResponse struct {
	Address     string
	Balance     map[string]float64
	Declaration *TxDeclaration
}

func CreateResponse(t string) (api.Response, bool) {
	switch t {
	case "GetLedgerResponse":
		return &GetLedgerResponse{}, true
	case "GetTransactionResponse":
		return &GetTransactionResponse{}, true
	case "SendTransactionResponse":
		return &SendTransactionResponse{}, true
	case "GetAccountResponse":
		return &GetAccountResponse{}, true
	}
	return nil, false
}

func (declaration *TxDeclaration) UnmarshalJSON(bytes []byte) error {
	var tmp struct {
		Type protocol.DeclarationType
	}

	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}

	declaration.Type = tmp.Type
	declaration.Declaration, err = CreateDeclaration(tmp.Type, bytes)

	return nil
}

func CreateDeclaration(t protocol.DeclarationType, bytes []byte) (interface{}, error) {
	switch t {
	case protocol.TxHashLock:
		var tmp = HashLock{}
		err := json.Unmarshal(bytes, &tmp)
		if err != nil {
			return nil, err
		}
		return tmp, nil
	case protocol.TxMultiSignature:
		var tmp = MultiSignature{}
		err := json.Unmarshal(bytes, &tmp)
		if err != nil {
			return nil, err
		}
		return tmp, nil
	case protocol.TxSecret:
		var tmp = SecretRevelation{}
		err := json.Unmarshal(bytes, &tmp)
		if err != nil {
			return nil, err
		}
		return tmp, nil
	}
	return nil, errors.New("Unknow declaration")
}

func (declaration *TxDeclaration) MarshalJSON() ([]byte, error) {

	var d interface{}
	switch declaration.Type {
	case protocol.TxSecret:
		var tmp struct {
			Type protocol.DeclarationType
			SecretRevelation
		}
		tmp.Type = declaration.Type
		tmp.SecretRevelation.Secret = declaration.Declaration.(*SecretRevelation).Secret
		d = tmp
	default:
		panic(0)
	}

	return json.Marshal(&d)
}
