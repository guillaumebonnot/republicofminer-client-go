package protocol

import (
	"republicofminer-client-go/crypto"
	"republicofminer-client-go/protocol/bytestream"
)

type Transaction struct {
	Expire       int64
	Fees         *TxInput
	Declarations []*TxDeclaration
	Inputs       []*TxInput
	Outputs      []*TxOutput
	Message      TransactionMessage
}

type TransactionMessage []byte

type TxInputOutput struct {
	Address  Address
	Amount   Amount
	Currency Currency
}

type TxInput TxInputOutput
type TxOutput TxInputOutput

type Amount int64

func (amount Amount) ToFloat() float64 {
	return float64(amount) / 100000000
}

func AmountFromFloat(float float64) Amount {
	return Amount(int64(float * 100000000))
}

type Currency int16

func CurrencyFromSymbol(symbol string) Currency {
	sum := 0
	multiplier := 1
	for i := 0; i < 3; i++ {
		sum += int((symbol[2-i] - 'A')) * multiplier
		multiplier *= 26
	}
	return Currency(int16(sum))
}

func (currency Currency) ToSymbol() string {
	buffer := int16(currency)
	array := make([]rune, 3)
	for i := 0; i < 3; i++ {
		var rest = buffer % 26
		array[2-i] = rune(byte(rest) + 'A')
		buffer /= 26
	}
	return string(array)
}

func (transaction *Transaction) Write(stream *bytestream.ByteStream) {
	stream.WriteInt64(transaction.Expire)
	stream.WriteNullable(bytestream.ByteStreamer(transaction.Fees.This()))
	stream.WriteList(len(transaction.Declarations), func(i int) bytestream.ByteStreamer { return transaction.Declarations[i] })
	stream.WriteList(len(transaction.Inputs), func(i int) bytestream.ByteStreamer { return transaction.Inputs[i].This() })
	stream.WriteList(len(transaction.Outputs), func(i int) bytestream.ByteStreamer { return transaction.Outputs[i].This() })
	stream.WriteNullable(bytestream.ByteStreamer(transaction.Message))
}

func (io *TxInputOutput) Write(stream *bytestream.ByteStream) {
	io.Address.Write(stream)
	io.Currency.Write(stream)
	io.Amount.Write(stream)
}

func (address *Address) Write(stream *bytestream.ByteStream) {
	stream.WriteByte(byte(address.Type))
	stream.WriteBytes(address.hash)
}

func (currency Currency) Write(stream *bytestream.ByteStream) {
	stream.WriteInt16(int16(currency))
}

func (amount Amount) Write(stream *bytestream.ByteStream) {
	stream.WriteInt64(int64(amount))
}

func (io *TxInput) This() *TxInputOutput {
	if io == nil {
		return nil
	}
	output := *io
	base := TxInputOutput(output)
	return &base
}

func (io *TxOutput) This() *TxInputOutput {
	if io == nil {
		return nil
	}
	output := *io
	base := TxInputOutput(output)
	return &base
}

func (declaration *TxDeclaration) Write(stream *bytestream.ByteStream) {
	stream.WriteByte(byte(declaration.Type))
	declaration.Declaration.Write(stream)
}

func (secret *SecretRevelation) Write(stream *bytestream.ByteStream) {
	stream.WriteBytes([]byte(secret.Secret))
}

func (secret TransactionMessage) Write(stream *bytestream.ByteStream) {
	stream.WriteBytes([]byte(secret))
}

func (transaction *Transaction) Hash() crypto.Hash256 {
	return crypto.Keccak256(bytestream.Write(transaction))
}
