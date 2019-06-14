package address32

import (
	"fmt"
	"strings"
)

const (
	RAW_SIZE     int = 21
	ENCODED_SIZE int = 40
	DECODED_SIZE int = 34
)

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var gen = []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

// Encode encodes the address type part and the data part returning an address32 encoded string
func Encode(code byte, data []byte) (string, error) {
	header, err := bytes8to5([]byte{code})
	if err != nil {
		return "", fmt.Errorf("failed to encode base32")
	}

	body, err := bytes8to5(data)
	if err != nil {
		return "", fmt.Errorf("failed to encode base32")
	}

	encoded, err := encode(append(header, body...))
	if err != nil {
		return "", fmt.Errorf("failed to encode base32")
	}

	return encoded, nil
}

func encode(decoded []byte) (string, error) {
	if len(decoded) != DECODED_SIZE {
		return "", fmt.Errorf("invalid address32 decoded length %d", len(decoded))
	}

	raw := append(decoded, createChecksum(decoded)...)

	encoded, err := toChars(raw)
	if err != nil {
		return "", fmt.Errorf("failed to encode data")
	}

	return encoded, nil
}

// Decode decodes a address32 encoded string,
// returning the address type part and the data part excluding the checksum.
func Decode(encoded string) (byte, []byte, error) {
	raw, err := decode(encoded)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to decode base32")
	}

	header, err := bytes5to8(raw[:2])
	if err != nil {
		return 0, nil, fmt.Errorf("failed to decode base32")
	}

	body, err := bytes5to8(raw[2:])
	if err != nil {
		return 0, nil, fmt.Errorf("failed to decode base32")
	}

	return header[0], body, nil
}

// decodes the base32 encoded string and remove the checksum
func decode(encoded string) ([]byte, error) {
	if len(encoded) != ENCODED_SIZE {
		return nil, fmt.Errorf("invalid address32 string length %d", len(encoded))
	}

	// We'll work with the lowercase string from now on.
	encoded, err := checkformat(encoded)
	if err != nil {
		return nil, err
	}

	// Each character corresponds to the byte with value of the index in
	// 'charset'.
	decoded, err := toBytes(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed converting data to bytes: %v", err)
	}

	if !verifyChecksum(decoded) {
		moreInfo := ""
		checksum := decoded[len(decoded)-6:]
		expected, err := toChars(createChecksum(decoded[:len(decoded)-6]))
		if err == nil {
			moreInfo = fmt.Sprintf("Expected %v, got %v.", expected, checksum)
		}
		return nil, fmt.Errorf("checksum failed. " + moreInfo)
	}

	// We exclude the last 6 bytes, which is the checksum.
	return decoded[:len(decoded)-6], nil
}

// For more details on the checksum verification, please refer to BIP 173.
func verifyChecksum(data []byte) bool {
	integers := make([]int, len(data))
	for i, b := range data {
		integers[i] = int(b)
	}
	return polymod(integers) == 1
}

// For more details on the polymod calculation, please refer to BIP 173.
func polymod(values []int) int {
	chk := 1
	for _, v := range values {
		b := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ v
		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 {
				chk ^= gen[i]
			}
		}
	}
	return chk
}

// For more details on the checksum calculation, please refer to BIP 173.
func createChecksum(data []byte) []byte {
	// Convert the bytes to list of integers, as this is needed for the
	// checksum calculation.
	integers := make([]int, len(data))
	for i, b := range data {
		integers[i] = int(b)
	}
	values := append(integers, []int{0, 0, 0, 0, 0, 0}...)
	polymod := polymod(values) ^ 1
	var res []byte
	for i := 0; i < 6; i++ {
		res = append(res, byte((polymod>>uint(5*(5-i)))&31))
	}
	return res
}

func checkformat(encoded string) (string, error) {
	// Only	ASCII characters between 33 and 126 are allowed.
	for i := 0; i < len(encoded); i++ {
		if encoded[i] < 33 || encoded[i] > 126 {
			return "", fmt.Errorf("invalid character in string: '%c'", encoded[i])
		}
	}

	// The characters must be either all lowercase or all uppercase.
	lower := strings.ToLower(encoded)
	upper := strings.ToUpper(encoded)
	if encoded != lower && encoded != upper {
		return "", fmt.Errorf("string not all lowercase or all uppercase")
	}

	return lower, nil
}

// toBytes converts each character in the string 'chars' to the value of the
// index of the correspoding character in 'charset'.
func toBytes(chars string) ([]byte, error) {
	decoded := make([]byte, 0, len(chars))
	for i := 0; i < len(chars); i++ {
		index := strings.IndexByte(charset, chars[i])
		if index < 0 {
			return nil, fmt.Errorf("invalid character not part of charset: %v", chars[i])
		}
		decoded = append(decoded, byte(index))
	}
	return decoded, nil
}

// toChars converts the byte slice 'data' to a string where each byte in 'data'
// encodes the index of a character in 'charset'.
func toChars(data []byte) (string, error) {
	result := make([]byte, 0, len(data))
	for _, b := range data {
		if int(b) >= len(charset) {
			return "", fmt.Errorf("invalid data byte: %v", b)
		}
		result = append(result, charset[b])
	}
	return string(result), nil
}

func bytes5to8(data []byte) ([]byte, error) {
	return ConvertBits(data, 5, 8, true)
}

func bytes8to5(data []byte) ([]byte, error) {
	return ConvertBits(data, 8, 5, true)
}

// ConvertBits converts a byte slice where each byte is encoding fromBits bits,
// to a byte slice where each byte is encoding toBits bits.
func ConvertBits(data []byte, fromBits, toBits uint8, pad bool) ([]byte, error) {
	if fromBits < 1 || fromBits > 8 || toBits < 1 || toBits > 8 {
		return nil, fmt.Errorf("only bit groups between 1 and 8 allowed")
	}

	// The final bytes, each byte encoding toBits bits.
	var regrouped []byte

	// Keep track of the next byte we create and how many bits we have
	// added to it out of the toBits goal.
	nextByte := byte(0)
	filledBits := uint8(0)

	for _, b := range data {

		// Discard unused bits.
		b = b << (8 - fromBits)

		// How many bits remaining to extract from the input data.
		remFromBits := fromBits
		for remFromBits > 0 {
			// How many bits remaining to be added to the next byte.
			remToBits := toBits - filledBits

			// The number of bytes to next extract is the minimum of
			// remFromBits and remToBits.
			toExtract := remFromBits
			if remToBits < toExtract {
				toExtract = remToBits
			}

			// Add the next bits to nextByte, shifting the already
			// added bits to the left.
			nextByte = (nextByte << toExtract) | (b >> (8 - toExtract))

			// Discard the bits we just extracted and get ready for
			// next iteration.
			b = b << toExtract
			remFromBits -= toExtract
			filledBits += toExtract

			// If the nextByte is completely filled, we add it to
			// our regrouped bytes and start on the next byte.
			if filledBits == toBits {
				regrouped = append(regrouped, nextByte)
				filledBits = 0
				nextByte = 0
			}
		}
	}

	// We pad any unfinished group if specified.
	if pad && filledBits > 0 {
		nextByte = nextByte << (toBits - filledBits)
		regrouped = append(regrouped, nextByte)
		filledBits = 0
		nextByte = 0
	}

	// Any incomplete group must be <= 4 bits, and all zeroes.
	if filledBits > 0 && (filledBits > 4 || nextByte != 0) {
		return nil, fmt.Errorf("invalid incomplete group")
	}

	return regrouped, nil
}
