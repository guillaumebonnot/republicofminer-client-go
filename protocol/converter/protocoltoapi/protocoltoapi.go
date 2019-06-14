// utility package that converts the object representation from protocol to api
package protocoltoapi

import (
	"republicofminer-client-go/explorer/api"
	"republicofminer-client-go/protocol"
)

func ToTransaction(transaction *protocol.Transaction) *api.Transaction {
	declarations := make([]*api.TxDeclaration, len(transaction.Declarations))
	for index, d := range transaction.Declarations {
		declarations[index] = ToDeclaration(d)
	}

	inputs := make([]*api.TxInput, len(transaction.Inputs))
	for index, d := range transaction.Inputs {
		inputs[index] = ToInput(d)
	}

	outputs := make([]*api.TxOutput, len(transaction.Outputs))
	for index, d := range transaction.Outputs {
		outputs[index] = ToOutput(d)
	}

	return &api.Transaction{
		Hash:         transaction.Hash().ToBase64(),
		Expire:       &transaction.Expire,
		Declarations: declarations,
		Inputs:       inputs,
		Outputs:      outputs,
		Message:      string([]byte(transaction.Message)),
		Fees:         ToInput(transaction.Fees),
	}
}

func ToDeclaration(transaction *protocol.TxDeclaration) *api.TxDeclaration {
	var declaration interface{}

	switch transaction.Type {
	case protocol.TxMultiSignature:
	case protocol.TxHashLock:
	case protocol.TxSecret:
		secret := transaction.Declaration.(*protocol.SecretRevelation)
		declaration = &api.SecretRevelation{Secret: secret.Secret.ToBase64()}
	case protocol.TxTimeLock:
	case protocol.TxVendingMachine:
	case protocol.TxLimitOrder:
	case protocol.TxDelegatedAccount:
	}

	return &api.TxDeclaration{
		Type:        transaction.Type,
		Declaration: declaration,
	}
}

func ToInput(input *protocol.TxInput) *api.TxInput {
	if input == nil {
		return nil
	}

	return &api.TxInput{
		Address:  input.Address.Encoded,
		Amount:   input.Amount.ToFloat(),
		Currency: input.Currency.ToSymbol(),
	}
}

func ToOutput(input *protocol.TxOutput) *api.TxOutput {
	if input == nil {
		return nil
	}

	return &api.TxOutput{
		Address:  input.Address.Encoded,
		Amount:   input.Amount.ToFloat(),
		Currency: input.Currency.ToSymbol(),
	}
}
