package protocol

import (
	"encoding/base64"
	"republicofminer-client-go/crypto"
	"republicofminer-client-go/protocol/bytestream"
	"republicofminer-client-go/protocol/format/address32"
)

type AddressType byte

const (
	ECDSA               AddressType = 0x1
	MultiSignatureECDSA AddressType = 0x2
	HashLock            AddressType = 0x3
	TimeLock            AddressType = 0x4
	VendingMachine      AddressType = 0x5
	LimitOrder          AddressType = 0x6
	DelegatedAccount    AddressType = 0x7
)

type Address struct {
	Encoded string
	Type    AddressType
	hash    []byte
}

func CreateAddress(typ AddressType, hash []byte) *Address {
	encoded, _ := address32.Encode(byte(typ), hash)
	return &Address{Type: typ, hash: hash, Encoded: encoded}
}

func DecodeAddress(encoded string) *Address {
	typ, hash, _ := address32.Decode(encoded)
	return &Address{Type: AddressType(typ), hash: hash, Encoded: encoded}
}

type DeclarationType byte

const (
	TxMultiSignature   DeclarationType = 0x0
	TxHashLock         DeclarationType = 0x1
	TxSecret           DeclarationType = 0x2
	TxTimeLock         DeclarationType = 0x3
	TxVendingMachine   DeclarationType = 0x4
	TxLimitOrder       DeclarationType = 0x5
	TxDelegatedAccount DeclarationType = 0x6
)

type TxDeclaration struct {
	Type        DeclarationType
	Declaration bytestream.ByteStreamer
}

type SecretRevelation struct {
	Secret Secret
	Hash   []byte
}

const SECRET_SIZE = 32

type Secret []byte

func (secret Secret) ToBase64() string {
	return base64.StdEncoding.EncodeToString([]byte(secret))
}

func SecretFromBase64(encoded string) Secret {
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	return Secret(decoded)
}

func NewSecretRevelation(secret Secret) *SecretRevelation {
	return &SecretRevelation{Secret: secret, Hash: crypto.Keccak256([]byte(secret))}
}
