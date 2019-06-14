// This package stores encrypted data in a sql database
package vault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"republicofminer-client-go/crypto"
)

const (
	CHECKITEM   = "check"
	CHECKSTRING = "this is a string to check"
)

// Unlock will unlock the target vault for future use
func Unlock(name, password string) bool {
	database = Database(name)
	// get the sample
	check, err := database.Item(CHECKITEM)

	secret = crypto.Keccak256([]byte(password))

	if err != nil {
		database.SetItem(CHECKITEM, encrypt([]byte(CHECKSTRING)))
		return true
	} else {
		if bytes.Compare(decrypt(check), []byte(CHECKSTRING)) == 0 {
			return true
		}
		fmt.Println("The pasword does not match")
		return false
	}
}

var database *VaultDatabase
var secret []byte

// CheckDatabase : Check if connected to database
// TODO IsUnlocked
func CheckDatabase() error {
	return nil
}

// Load will load and decrypt the requested item from the database
func Load(item string) ([]byte, error) {
	if err := CheckDatabase(); err != nil {
		return nil, err
	}
	data, err := database.Item(item)

	if err != nil {
		return nil, err
	}

	return decrypt(data), nil
}

// Save will save and encrypt the requested item in the database
func Save(item string, bytes []byte) error {
	if err := CheckDatabase(); err != nil {
		return err
	}

	return database.SetItem(item, encrypt([]byte(bytes)))
}

func encrypt(plaintext []byte) []byte {
	block, err := aes.NewCipher(secret)
	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext
}

func decrypt(cyphertext []byte) []byte {
	block, err := aes.NewCipher(secret)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := cyphertext[:nonceSize], cyphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func delete() {
	database = nil
	database.Delete()
}
