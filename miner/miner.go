package miner

import (
	"bytes"
	"encoding/base64"
	"math/rand"
	"republicofminer-client-go/explorer"
	"republicofminer-client-go/explorer/api"
	"republicofminer-client-go/protocol"
	"republicofminer-client-go/protocol/converter/protocoltoapi"
	"republicofminer-client-go/republicofminer"
	"republicofminer-client-go/wallet"
	"time"
)

func Run() {
	go explorer.Connect()
	go republicofminer.Connect()
	wallet.Load()

	for {
		task := republicofminer.GetMiningTask(wallet.Address.Encoded, resource())
		hash, _ := base64.StdEncoding.DecodeString(task.SecretHash)
		mask, _ := base64.StdEncoding.DecodeString(task.Mask)
		secret := mine(hash, mask)
		address := protocol.DecodeAddress(task.Address)
		amount := protocol.Amount(task.Amount)
		currency := protocol.CurrencyFromSymbol(task.Currency)
		transaction := claim(*address, *wallet.Address, amount, currency, secret)
		pub, signature := wallet.Sign(transaction.Hash().ToBytes())
		explorer.SendTransaction(protocoltoapi.ToTransaction(transaction), []*api.Signature{&api.Signature{
			PublicKey:     pub.ToBase64(),
			SignatureByte: signature.ToBase64(),
		}})
	}
}

// we try to find the secret matching with the given secret hash
func mine(secret []byte, mask []byte) *protocol.SecretRevelation {
	complexity := protocol.SECRET_SIZE - len(mask)
	// the mask is the first part of the secret
	buffer := append(mask, make([]byte, complexity)...)
	last := len(buffer) - 1
	for {
		// randomize the unknown part
		rand.Read(buffer[len(buffer)-complexity:])

		// brute force the last byte
		for b := byte(0); ; b++ {
			buffer[last] = b
			h := protocol.NewSecretRevelation(protocol.Secret(buffer))
			// check if the hash matches with the secret
			if bytes.Equal(secret, h.Hash) {
				// fmt.Println("Secret Hash found !")
				return h
			}
			if b == 255 {
				break
			}
		}
	}
}

func claim(sender protocol.Address, receiver protocol.Address, amount protocol.Amount, currency protocol.Currency, secret *protocol.SecretRevelation) *protocol.Transaction {
	transaction := protocol.Transaction{
		Expire:       time.Now().Add(time.Minute * 10).Unix(),
		Declarations: []*protocol.TxDeclaration{&protocol.TxDeclaration{Type: protocol.TxSecret, Declaration: secret}},
		Inputs:       []*protocol.TxInput{&protocol.TxInput{Address: sender, Amount: amount, Currency: currency}},
		Outputs:      []*protocol.TxOutput{&protocol.TxOutput{Address: receiver, Amount: amount, Currency: currency}},
	}
	return &transaction
}

// TODO strategy to decide which resource to mine
func resource() string {
	// get one at random
	candidates := []string{"WOD", "STN", "IRO"}
	return candidates[rand.Intn(len(candidates))]
}
