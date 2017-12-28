package assembler_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/thee-engineer/cryptor/cachedb"
	"github.com/thee-engineer/cryptor/cachedb/ldbcache"
	"github.com/thee-engineer/cryptor/chunker"
	"github.com/thee-engineer/cryptor/chunker/assembler"
	"github.com/thee-engineer/cryptor/crypt"
	"github.com/thee-engineer/cryptor/crypt/aes"
	"github.com/thee-engineer/cryptor/crypt/hashing"
)

func TestEChunk(t *testing.T) {
	t.Parallel()

	// Create temporary dir for test
	tmpDir, err := ioutil.TempDir("/tmp", "assembler")
	defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Error(err)
	}

	// Create temp cache
	db, err := ldbcache.New(tmpDir, 0, 0)
	if err != nil {
		t.Error(err)
	}
	cache := ldbcache.NewManager(cachedb.DefaultManagerConfig, db)

	// Test data
	var buffer bytes.Buffer
	data := crypt.RandomData(520)
	buffer.Write(data)

	// Create chunker
	chunkHash, err := chunker.ChunkFrom(&buffer, 1024, cache, aes.NullKey)
	if err != nil {
		t.Error(err)
	}

	// Read encrypted chunk
	eChunk := assembler.GetEChunk(chunkHash, db)
	dChunk, err := eChunk.Decrypt(aes.NullKey)
	if err != nil {
		t.Error(err)
	}

	// Invalid hash
	if !dChunk.IsValid() {
		t.Log(crypt.EncodeString(dChunk.Header.Hash))
		t.Log(crypt.EncodeString(hashing.SHA256Digest(dChunk.Content)))
		t.Error("chunk: is not valid")
	}

	// Chunk should be the tail (as it is the only chunk)
	if !dChunk.IsLast() {
		t.Error("chunk: is not last")
	}

	// Compare initial data with data after encryption, storage and decryption
	if !bytes.Equal(dChunk.Content, data) {
		t.Log("init", data)
		t.Log("decr", dChunk.Content)
		t.Error("data mismatch: initial package data and assembled chunks")
	}
}
