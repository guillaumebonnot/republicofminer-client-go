package address32

import "testing"

func TestEncoding(t *testing.T) {
	// Private 7r7oFxKhhaH7UvMLpUXlcIEk0WWx7i4nw6BVnrKCmLk=
	address := "qyl68tygnjx6qqwrsmynmejmc9wxlw7almv3397j"
	typ, hash, err := Decode(address)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := Encode(typ, hash)
	if err != nil {
		t.Fatal(err)
	}

	if address != encoded {
		t.Fatal("encoding and decoding does not work")
	}
}
