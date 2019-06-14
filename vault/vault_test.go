package vault

import (
	"bytes"
	"republicofminer-client-go/crypto"
	"testing"
)

func TestEncryptDecypt(t *testing.T) {

	secret = crypto.Keccak256([]byte("ansdfsd45f141as41fas1ds1f1"))
	plaintext := []byte("as4da1dd4qd4s1ad7qd54q1d541q4w1d45q154d")
	encrypted := encrypt(plaintext)
	decrypted := decrypt(encrypted)

	if bytes.Equal(encrypted, decrypted) {
		t.Errorf("encrypt + decrypt does not work")
	}
}

func TestVault(t *testing.T) {
	name := "someitem"
	content := []byte("some important stuff")
	ok := Unlock("test", "thisisapassword")

	if !ok {
		t.Errorf("could not unlock an empty vault")
		return
	}

	item, err := Load(name)
	if err != nil {
		t.Errorf("the vault database should be empty")
		return
	}

	err = Save(name, content)
	if err != nil {
		t.Errorf("error saving an item in the vault")
		return
	}

	item, err = Load(name)
	if err == nil {
		t.Errorf("the item should be in the vault database")
		return
	}

	if bytes.Equal(content, item) {
		t.Errorf("encrypt + decrypt does not work")
	}
}
