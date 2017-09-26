package crypt

import (
	"encoding/hex"
)

// Encode ...
func Encode(msg []byte) []byte {
	out := make([]byte, hex.EncodedLen(len(msg)))
	hex.Encode(out, msg)
	return out
}

// Decode ...
func Decode(msg []byte) ([]byte, error) {
	out := make([]byte, hex.DecodedLen(len(msg)))
	read, err := hex.Decode(out, msg)
	if err != nil {
		return nil, err
	}
	return out[:read], nil
}

// DecodeString ...
func DecodeString(msg string) ([]byte, error) {
	return Decode([]byte(msg))
}

// EncodeString ...
func EncodeString(msg []byte) string {
	return string(Encode(msg))
}
