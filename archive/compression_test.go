package archive

import (
	"bytes"
	"testing"

	"github.com/thee-engineer/cryptor/crypt"
)

func TestCompression(t *testing.T) {
	t.Parallel()

	// Compress random bytes
	_, err := Compress(crypt.RandomData(100))
	if err != nil {
		t.Error(err)
	}

	// Compress no bytes
	_, err = Compress([]byte{})
	if err != nil {
		t.Error(err)
	}
}

func TestDecompression(t *testing.T) {
	t.Parallel()

	// Create random data and compress it
	initialData := crypt.RandomData(100)
	buffer, err := Compress(initialData)
	if err != nil {
		t.Error(err)
	}

	// Decompress random compressed data
	data, err := Decompress(buffer)
	if err != nil {
		t.Error(err)
	}

	// Compare initial data with uncompressed data
	if !bytes.Equal(data, initialData) {
		t.Error("data mismatch: uncompressed data does not match initial data")
	}
}
