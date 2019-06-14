package protocol

import (
	"encoding/base64"
	"fmt"
	"republicofminer-client-go/crypto"
	"testing"
)

var keys = []string{
	"7r7oFxKhhaH7UvMLpUXlcIEk0WWx7i4nw6BVnrKCmLk=",
	"QkQ07ERvruv0idJ6e0xDX2GbpYcuLB+dPueldNyd5xA=",
	"cyg1lEqj83Xa00KHdtab5tCyvhKVC6q84YkEFnERbPI=",
	"WgqyI72GAn/up/H6uSajmxh3njrbXjDxoxXA56x69JQ=",
	"yNV1hPckx4F/uyh9Y6OwbzlYMDVU/DhJlUs7zkr44DE=",
}

var signatures = []string{
	`AcNFVS/cGY5wRzMLhmhQ/7unELDEwGPkj3I2fPb5j0fwAztZVW+UTDL33183EuO1dBYI+2D161zJvi58qG2WlIk=`,
	`Aehh93+JfZ1HdY/sWw+IENjdBoveYjyXBVgIPCISR2EMUwe1HGWo5evVxlYiNr0JzR838q7V9pLUZbl7mPxFCBA=`,
	`AFvdQB7J+xJqrEKkLfwL1gDGuSP/0Mrx/A06kOEIkpNFdtArURTBMh8l8gHHn91U5CL+AC/7H7Vxtw9b+JH3Kzw=`,
	`AHDDFwCGPUcVUP/ykTlUlEJusme11WbW1om2heeYmzsPfzBZUN5v0BA5YjeQA33fU8xFWgRl1uczXWdoBMkfDR8=`,
	`AchJ1ZqYP89O8Xay8Px3iMh2VMw/d3zNWswuWSSR4W6MTuEPIeuTwMNdaVdnIkFpjfviJDdfcw+bUmx8UDHzo5o=`,
}

func TestSignature(t *testing.T) {
	// get the private key
	for index, base64 := range keys {
		key, err := PrivateKeyFromBase64(base64)
		if err != nil {
			t.Fatal("Error decoding private key from base 64")
		}
		if base64 != key.ToBase64() {
			t.Fatal("Error converting private key to and from base 64")
		}

		// derive the public key
		var pub = key.GetPublicKey()

		// encode the address
		var address = pub.GetAddress()

		var text = []byte("Guillaume is the best !")
		hash := crypto.Hash(text)
		// sign hash
		signed, err := key.SignMessage(hash, Network)
		if err != nil {
			t.Fatal("Error compute the signature :", err)
		}

		// check signature against public key
		if !pub.CheckSignature(hash, signed, Network) {
			t.Fatal("Error signature is invalid !")
		}

		// check public key against address
		if !pub.CheckAddress(address.Encoded) {
			t.Fatal("The address doesn't match the public key!")
		}

		expected := signatures[index]
		actual := signed.ToBase64()
		if actual != expected {
			fmt.Println("expected :", expected, "actual", actual)
			t.Fail()
		}
	}
}

func TestCurrency(t *testing.T) {
	expected := "IRO"
	currency := CurrencyFromSymbol(expected)
	actual := currency.ToSymbol()
	if actual != expected {
		fmt.Println("expected :", expected, "actual", actual)
		t.Fail()
	}
}

func TestSerialization1(t *testing.T) {
	// {"result":0,"type":"GetTransactionResponse","data":{"Transaction":{"Hash":"zIJZB67U0gTUnGq649baM/5ylbUE1ydm5WpJ7xn2XfQ=","Expire":1556277083,"Declarations":[{"Secret":"K/w77o6eeCFUiLzq6jTzKrQAoAn8BrurFaztYzsP68s=","Type":2}],"Inputs":[{"Address":"qg64nhvuzlj2lenndj3mg89gcswkuc3axtq2v40s","Currency":"IRO","Amount":0.00000001}],"Outputs":[{"Address":"qyunuamu8u9axnx8e6y0809qup2599snluyccvd2","Currency":"IRO","Amount":0.00000001}]}},"crid":"1560396570-0"}
	secret, _ := base64.StdEncoding.DecodeString("K/w77o6eeCFUiLzq6jTzKrQAoAn8BrurFaztYzsP68s=")
	transaction := Transaction{
		Expire:       1556277083,
		Declarations: []*TxDeclaration{&TxDeclaration{Type: TxSecret, Declaration: NewSecretRevelation(Secret(secret))}},
		Inputs:       []*TxInput{&TxInput{Address: *DecodeAddress("qg64nhvuzlj2lenndj3mg89gcswkuc3axtq2v40s"), Currency: CurrencyFromSymbol("IRO"), Amount: AmountFromFloat(0.00000001)}},
		Outputs:      []*TxOutput{&TxOutput{Address: *DecodeAddress("qyunuamu8u9axnx8e6y0809qup2599snluyccvd2"), Currency: CurrencyFromSymbol("IRO"), Amount: AmountFromFloat(0.00000001)}},
	}

	expected := "zIJZB67U0gTUnGq649baM/5ylbUE1ydm5WpJ7xn2XfQ="
	actual := transaction.Hash().ToBase64()
	if actual != expected {
		fmt.Println("expected :", expected, "actual", actual)
		t.Fail()
	}
}

func TestSerialization2(t *testing.T) {
	// {"type":"SendTransactionRequest","data":{"Transaction":{"Hash":"KAapdGf1unoM8dSsN+SHkqsQKXDP2Y962RnkanRFYcg=","Expire":1560404881,"Declarations":[{"Type":2,"Secret":"NXYBRplyY/bDfjBzVNppa/PzPvIDlOxK3j3urVVh4Jk="}],"Inputs":[{"Address":"qgefrlzgsx998sj9lvj4hw39plh22llxwlj4tuvp","Currency":"WOD","Amount":1e-8}],"Outputs":[{"Address":"qy2t4fvr6q5k0235p5xg5wu64tn883ks20cg424c","Currency":"WOD","Amount":1e-8}]},"Signatures":[{"k":"BKvF06nSi8SbX5pYVapqwQ57EwLbM11bOQAnN616U+llRwlVHgW8/4znqMtqk03eqhq9LzE2cbt9oGzoVKrq2oY=","s":"AEGihBiwr8l99RWLKxMFDnu9rHrkGFKqjj63xWXrRTTmCI/P4sEOkznSdSezz7pW15zEYXREri5K9DRGc3qV680="}]},"crid":"1560404279-0"}
	secret, _ := base64.StdEncoding.DecodeString("NXYBRplyY/bDfjBzVNppa/PzPvIDlOxK3j3urVVh4Jk=")
	transaction := Transaction{
		Expire:       1560404881,
		Declarations: []*TxDeclaration{&TxDeclaration{Type: TxSecret, Declaration: NewSecretRevelation(Secret(secret))}},
		Inputs:       []*TxInput{&TxInput{Address: *DecodeAddress("qgefrlzgsx998sj9lvj4hw39plh22llxwlj4tuvp"), Currency: CurrencyFromSymbol("WOD"), Amount: AmountFromFloat(0.00000001)}},
		Outputs:      []*TxOutput{&TxOutput{Address: *DecodeAddress("qy2t4fvr6q5k0235p5xg5wu64tn883ks20cg424c"), Currency: CurrencyFromSymbol("WOD"), Amount: AmountFromFloat(0.00000001)}},
	}

	expected := "KAapdGf1unoM8dSsN+SHkqsQKXDP2Y962RnkanRFYcg="
	actual := transaction.Hash().ToBase64()
	if actual != expected {
		fmt.Println("expected :", expected, "actual", actual)
		t.Fail()
	}
}

func TestSecretHash(t *testing.T) {
	var secrets = []string{
		"sm191XFhlz6aHH8xB0ZphrnExkdoaykYmYNRozysd0A=",
		"Pt449sVdjyskOYHNzvtMU5p9quNOLOOqSLge83VB0eE=",
		"lOER9EjpxTlRGJjIBoMBRhoOsf+1eCd3uF+gvonq9jk=",
		"2rQ1LAvxkuo/IzgySFzRW5sFC8UWlZkcupo86UkHsGY=",
		"IV/tUk3VOyl7x7+RP/QuW41EIoPfMCJEGhlxcbWooDg=",
		"c9ER6G5/DcbK4klP1QfcxmvgWwXkLIfwPag3NGTJdNM=",
		"JocdbuNF7Q02XLj959NvEj2qBtRNm1RHSiCrxUjZYgA=",
		"bXYoTepu4yXyTc6HdIkJajkk97nDn4WLM9YjwponxVw=",
		"+rePcaohKM7H4OH/zMzSwBZAjCPiqD3O8JwaE0ogiss=",
		"h3GNvPqYvxf8KPLjvYa/ZgNirFmbonpResGADQYVhKY=",
		"VRZ3Wj8Wxn6AdCe6tlkQat20D+h0MC/l9ztBOBexEeY=",
		"pEuNpg5kJu2bLdgne7+bILAWMFIFmzEk55JjBNSNZ88=",
	}

	var hashes = []string{
		"oR/pKR+lTax0EAntbph8539SDmS0gjPoU+ZHXwmU5b4=",
		"+Jzz4Ysu0fsFpgN6VpqqN1ikmDhNlBa/pDz7jEM/pYo=",
		"fFJ6JWkRUr2GXu21PsQa/oWzLL25Cmr6E9OKErjq1/g=",
		"L/NVV+GInHJzyAM58Odhf1T3OkFXWZp+k0sJsNQfUVM=",
		"txAH+9IGwp0QpUfbpLURWUImIj7wHsvENDNFt3WCmUI=",
		"HcBpQQmOfCRHzrAxYEx1AZ7KJV9CxCANWYf3ILQdtJg=",
		"OEU55Q19YtcXOjfNPc6gaz4in9QFw5GeK6H5PQUDYmA=",
		"whUpZhNx6pS3Nsh4i8+g14b+XSwCJB8LuDmjN7iPWwU=",
		"RP2t8bhobT7AGFFyrPZfDm1S0B/DpHYlyRPfGXGkvg4=",
		"GTG37Ie0cWbD0N+ed6Y4cEUdNmTVVV1E553+K2+4QMg=",
		"weOWnxlVynrYTDNk/t47NwZOJHClkVwJvCQ4Du3t4Vo=",
		"9t3Ym1JuvLfXezD8v+zfwL/wttzx/PqAhODNWo/xfT4=",
	}

	for index := 0; index < len(hashes); index++ {
		secret := NewSecretRevelation(SecretFromBase64(secrets[index]))

		expected := hashes[index]
		actual := crypto.Hash256(secret.Hash).ToBase64()
		if actual != expected {
			fmt.Println("expected :", expected, "actual", actual)
			t.Fail()
		}

	}
}

func TestPublicKey(t *testing.T) {
	var pubs = []string{
		"BFaibfgtfJpKDciB6uTzvRRmXGTBKJl4CGdehor3tToAtj3+W8dWWCUFpWRw5a2ADsJoXm9Ja5lYEvwtUq76Vt0=",
		"BJBOzQ+1XuAaevISCJ4eiDG4qDFao/m9cA9aZBtYQepRd3cdi1bmSoGXOtXJ/86E7nLbMAvKOwFxCATJyiB0JDI=",
		"BOb/VI0QxfejcgWbthn4zj67qx7/M0t+/F7ytp7Lr3jhpxUjB27v9zcRXMfS5bqRq4KZDf6Z8VOuzwLzVwrwBZg=",
		"BJNLu2O6pb67DtOoNkUHz2vdCYMbrXZL5UTRpCnomQoraWyX2VX7t6aYUZMS1M9makcYdwuFo2tzAxyhPv3oC/s=",
		"BJrh2sqP9zcTF2bnIdf4n8xV33R4FwXdPeAhsfh+5cuWDjLHkeaz3vDbiA90qhFhkUttZDfG/fihJ8/ftG9r2hc=",
	}

	for index := 0; index < len(keys); index++ {
		pk, _ := PrivateKeyFromBase64(keys[index])
		expected := pubs[index]
		actual := pk.GetPublicKey().ToBase64()
		if actual != expected {
			fmt.Println("expected :", expected, "actual", actual)
			t.Fail()
		}
	}
}

func TestAddress(t *testing.T) {
	var addresses = []string{
		"qyl68tygnjx6qqwrsmynmejmc9wxlw7almv3397j",
		"qyaj20aksyvxlfmznyjdqzxrvvf0w7ca7mamwzll",
		"qymcnhdqy2qth0ls0pel7ln4gmps2vkxcaq78kcn",
		"qyn8zz6ys2wkxr8r5e3gsnjny6v7hfc9qmg7netx",
		"qyy2m8wndvy8aleqgw0z89m2m4k72fvxn6t35zcv",
	}

	for index := 0; index < len(keys); index++ {
		pk, _ := PrivateKeyFromBase64(keys[index])

		expected := addresses[index]
		actual := pk.GetPublicKey().GetAddress().Encoded
		if actual != expected {
			fmt.Println("expected :", expected, "actual", actual)
			t.Fail()
		}
	}
}
