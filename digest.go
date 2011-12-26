package cstore

import (
	"crypto"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

// Compute an SHA256 digest for a string.
func Digest(data string) string {
	hash := crypto.SHA256.New()
	if _, err := hash.Write([]byte(data)); err != nil {
		panic("Writing to a hash should never fail")
	}
	return hex.EncodeToString(hash.Sum())
}

// A reader which wraps another reader and computes a running hash code
// on the data as it is read.
type DigestReader struct {
	nested io.Reader
	hash hash.Hash
}

// Create a new DigestReader wrapping the specified reader.
func NewDigestReader(r io.Reader) *DigestReader {
	return &DigestReader{nested: r, hash: crypto.SHA256.New()}
}

// Read data from the nested reader and update the hash data.
func (dr *DigestReader) Read(p []byte) (n int, err os.Error) {
	n, err = dr.nested.Read(p)
	if n > 0 {
		_, werr := dr.hash.Write(p[:n])
		if werr != nil {
			panic("Writing to a hash should never fail")
		}
	}
	return
}

// Get the digest for all data read so far.
func (dr *DigestReader) Digest() string {
	return hex.EncodeToString(dr.hash.Sum())
}
