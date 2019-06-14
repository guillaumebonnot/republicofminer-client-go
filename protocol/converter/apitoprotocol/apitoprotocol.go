// utility package that converts the object representation from api to protocol
package apitoprotocol

import (
	"encoding/base64"
	"republicofminer-client-go/explorer/api"
	"republicofminer-client-go/protocol"
	"republicofminer-client-go/protocol/bytestream"
)

func ToTransaction(transaction *api.Transaction) *protocol.Transaction {
	declarations := make([]*protocol.TxDeclaration, len(transaction.Declarations))
	for index, d := range transaction.Declarations {
		declarations[index] = ToDeclaration(d)
	}

	inputs := make([]*protocol.TxInput, len(transaction.Inputs))
	for index, d := range transaction.Inputs {
		inputs[index] = ToInput(d)
	}

	outputs := make([]*protocol.TxOutput, len(transaction.Outputs))
	for index, d := range transaction.Outputs {
		outputs[index] = ToOutput(d)
	}

	return &protocol.Transaction{
		Expire:       *transaction.Expire,
		Declarations: declarations,
		Inputs:       inputs,
		Outputs:      outputs,
		Message:      protocol.TransactionMessage([]byte(transaction.Message)),
		Fees:         ToInput(transaction.Fees),
	}
}

func ToDeclaration(transaction *api.TxDeclaration) *protocol.TxDeclaration {
	var declaration bytestream.ByteStreamer

	switch transaction.Type {
	case protocol.TxMultiSignature:
	case protocol.TxHashLock:
	case protocol.TxSecret:
		secret := transaction.Declaration.(api.SecretRevelation)
		decoded, _ := base64.StdEncoding.DecodeString(secret.Secret)
		declaration = protocol.NewSecretRevelation(protocol.Secret(decoded))
	case protocol.TxTimeLock:
	case protocol.TxVendingMachine:
	case protocol.TxLimitOrder:
	case protocol.TxDelegatedAccount:
	default:
		panic(0)
	}

	return &protocol.TxDeclaration{
		Type:        transaction.Type,
		Declaration: declaration,
	}
}

func ToInput(input *api.TxInput) *protocol.TxInput {
	if input == nil {
		return nil
	}

	return &protocol.TxInput{
		Address:  *protocol.DecodeAddress(input.Address),
		Amount:   protocol.AmountFromFloat(input.Amount),
		Currency: protocol.CurrencyFromSymbol(input.Currency),
	}
}

func ToOutput(input *api.TxOutput) *protocol.TxOutput {
	if input == nil {
		return nil
	}

	return &protocol.TxOutput{
		Address:  *protocol.DecodeAddress(input.Address),
		Amount:   protocol.AmountFromFloat(input.Amount),
		Currency: protocol.CurrencyFromSymbol(input.Currency),
	}
}
