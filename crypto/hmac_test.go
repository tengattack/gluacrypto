package gluacrypto_crypto_test

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluacrypto"
	lua "github.com/yuin/gopher-lua"
)

var (
	algorithms = []string{"md5", "sha1", "sha256", "sha512"}
)

func getHasher(algorithm string) func() hash.Hash {
	var h func() hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New
	case "sha1":
		h = sha1.New
	case "sha256":
		h = sha256.New
	case "sha512":
		h = sha512.New
	}
	return h
}

func TestHmac(t *testing.T) {
	assert := assert.New(t)

	// test start
	for _, algorithm := range algorithms {
		L := lua.NewState()
		defer L.Close()
		gluacrypto.Preload(L)

		h := hmac.New(getHasher(algorithm), Key)
		h.Write(Data)
		hashData := h.Sum(nil)

		script := `
		  crypto = require('crypto')
		  return crypto.hmac('` + algorithm + `', '` + string(Data) + `', '` + string(Key) + `')
	  `
		assert.NoError(L.DoString(script))

		val := getValue(L, L.Get(1))
		err := getValue(L, L.Get(2))
		assert.Nil(err)
		assert.Equal(hex.EncodeToString(hashData), val)
	}
}

func TestHmacRaw(t *testing.T) {
	assert := assert.New(t)

	// test start
	for _, algorithm := range algorithms {
		L := lua.NewState()
		defer L.Close()
		gluacrypto.Preload(L)

		h := hmac.New(getHasher(algorithm), Key)
		h.Write(Data)
		hashData := h.Sum(nil)

		script := `
		  crypto = require('crypto')
		  return crypto.hmac('` + algorithm + `', '` + string(Data) + `', '` + string(Key) + `', true)
	  `
		assert.NoError(L.DoString(script))

		val := getValue(L, L.Get(1))
		err := getValue(L, L.Get(2))
		assert.Nil(err)
		assert.Equal(string(hashData), val)
	}
}

func TestHmacFail(t *testing.T) {
	assert := assert.New(t)

	L := lua.NewState()
	defer L.Close()
	gluacrypto.Preload(L)

	script := `
	  crypto = require('crypto')
	  return crypto.hmac('unknown', '` + string(Data) + `', '` + string(Key) + `')
	`
	assert.NoError(L.DoString(script))

	val := getValue(L, L.Get(1))
	err := getValue(L, L.Get(2))
	assert.NotNil(err)
	assert.Nil(val)
}
