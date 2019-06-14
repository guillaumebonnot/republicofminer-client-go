package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/big"

	"golang.org/x/crypto/sha3"
)

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
	one            = big.NewInt(1)
)

var errInvalidPubkey = errors.New("invalid secp256k1 public key")

type Hash256 []byte

func (hash Hash256) ToBase64() string {
	return base64.StdEncoding.EncodeToString([]byte(hash))
}

func (hash Hash256) ToBytes() []byte {
	return []byte(hash)
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) Hash256 {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// Hash is the double SHA256 used by bitcoin
func Hash(data ...[]byte) Hash256 {
	return SHA256(SHA256(data...))
}

func SHA256(data ...[]byte) Hash256 {
	d := sha256.New()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
