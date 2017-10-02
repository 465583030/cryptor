package chunker

import (
	"bytes"

	"github.com/thee-engineer/cryptor/crypt"
)

// NullByteArray is used for the last chunk header.Next
var NullByteArray [32]byte

// Chunk combines the chunk content which is a []byte of header.Size and the
// chunk header which contains information about the this chunk and next chunk.
type Chunk struct {
	Header  *ChunkHeader
	Content []byte
}

// NewChunk creates a new chunk with given size content
func NewChunk(size uint32) *Chunk {
	return &Chunk{
		Header:  NewChunkHeader(),
		Content: make([]byte, size),
	}
}

// Bytes returns the chunk header and content as []byte
func (c Chunk) Bytes() []byte {
	var buffer bytes.Buffer

	buffer.Write(c.Header.Bytes()) // Write header bytes
	buffer.Write(c.Content)        // Write chunk content

	return buffer.Bytes()
}

// IsValid compares the header hash with the content hash
func (c Chunk) IsValid() bool {
	return bytes.Equal(c.Header.Hash, crypt.SHA256Data(c.Content).Sum(nil))
}

// IsLast checks if the next chunk hash is the NullByteArray
func (c Chunk) IsLast() bool {
	return bytes.Equal(c.Header.Next, NullByteArray[:])
}
