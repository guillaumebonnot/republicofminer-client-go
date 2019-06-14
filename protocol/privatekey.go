package protocol

import (
	"encoding/base64"
	"math/big"
	"republicofminer-client-go/crypto"

	"github.com/btcsuite/btcd/btcec"
)

const KEY_SIZE int = 32

var Network []byte = []byte("republicofminer.com")

type PrivateKey struct {
	key *btcec.PrivateKey
}

type PublicKey struct {
	key *btcec.PublicKey
}

type Signature struct {
	signature *btcec.Signature
}

func PrivateKeyFromBase64(encoded string) (*PrivateKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), decoded)

	return &PrivateKey{priv}, nil
}

func PrivateKeyFromBytes(bytes []byte) *PrivateKey {
	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), bytes)
	return &PrivateKey{priv}
}

func GeneratePrivateKey() *PrivateKey {
	pk, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		panic(err)
	}
	return &PrivateKey{pk}
}

func (key *PrivateKey) ToBytes() []byte {
	return key.key.Serialize()
}

func (key *PrivateKey) ToBase64() string {
	return base64.StdEncoding.EncodeToString(key.ToBytes())
}

func (key *PrivateKey) GetPublicKey() *PublicKey {
	return &PublicKey{key.key.PubKey()}
}

func (key *PrivateKey) SignMessage(message []byte, network []byte) (*Signature, error) {
	s, err := key.key.Sign(prepare(message, network))
	return &Signature{s}, err
}

func PublicKeyFromBase64(encoded string) (*PublicKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	pub, _ := btcec.ParsePubKey(decoded, btcec.S256())

	return &PublicKey{pub}, nil
}

func (key *PublicKey) GetAddress() *Address {
	serialized := key.key.SerializeUncompressed()
	hash := crypto.Keccak256(serialized)

	// take the last 20 bytes;
	return CreateAddress(ECDSA, hash[len(hash)-20:])
}

func (key *PublicKey) CheckSignature(data []byte, signature *Signature, network []byte) bool {
	return signature.signature.Verify(prepare(data, network), key.key)
}

func (key *PublicKey) CheckAddress(encoded string) bool {
	return key.GetAddress().Encoded == encoded
}

func prepare(data []byte, network []byte) []byte {
	return append(network, data...)
}

func (key *PublicKey) ToBytes() []byte {
	return key.key.SerializeUncompressed()
}

func (key *PublicKey) ToBase64() string {
	return base64.StdEncoding.EncodeToString(key.ToBytes())
}

func SignatureFromBase64(encoded string) (*Signature, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	left := decoded[1:33]
	right := decoded[33:]

	r := FromByteArray(left)
	s := FromByteArray(right)

	return &Signature{&btcec.Signature{R: r, S: s}}, nil
}

func (sig *Signature) ToBytes() []byte {
	overflow := byte(0)
	left, isOverflow := ToByteArray(sig.signature.R)
	if isOverflow {
		overflow |= 0x1
	}

	right, isOverflow := ToByteArray(sig.signature.S)
	if isOverflow {
		overflow |= 0x2
	}

	buffer := make([]byte, 1)
	buffer[0] = overflow
	buffer = append(append(buffer, left...), right...)
	return buffer
}

func ToByteArray(i *big.Int) ([]byte, bool) {
	bytes := i.Bytes()
	length := len(bytes)
	if length == 33 {
		// should not happen because we dont have the sign byte included
		panic(1)
		// return bytes[1:], true
	}
	if length == 32 {
		return bytes, overflow(bytes[0])
	}
	if length < 32 {
		return append(make([]byte, 32-length), bytes...), false
	}
	panic(0)
}

func FromByteArray(bytes []byte) *big.Int {
	i := &big.Int{}
	return i.SetBytes(bytes)
}

// we try emulate adding the sign to be compatible with bouncy castle
func overflow(b byte) bool {
	return b > 0x7f
}

func (sig *Signature) ToBase64() string {
	return base64.StdEncoding.EncodeToString(sig.ToBytes())
}
