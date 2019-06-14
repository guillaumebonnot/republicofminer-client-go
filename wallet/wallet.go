package wallet

import (
	"fmt"
	"log"
	"republicofminer-client-go/protocol"
	"republicofminer-client-go/vault"
)

var Privatekey *protocol.PrivateKey
var Publickey *protocol.PublicKey
var Address *protocol.Address

func Load() {
	vault.Unlock("republicofminer", "8dLyWpyupBty")
	pk, err := vault.Load("wallet")
	if err != nil {
		err := vault.Save("wallet", Privatekey.ToBytes())
		if err != nil {
			fmt.Println("Save Failed", err)
			return
		}
	} else {
		Privatekey = protocol.PrivateKeyFromBytes(pk)
	}
	Publickey = Privatekey.GetPublicKey()
	Address = Publickey.GetAddress()

	log.Println("Loaded wallet :", Address.Encoded)
	// fmt.Println("Private key :", Privatekey.ToBase64())
}

func Sign(data []byte) (*protocol.PublicKey, *protocol.Signature) {
	signature, _ := Privatekey.SignMessage(data, protocol.Network)
	return Publickey, signature
}
