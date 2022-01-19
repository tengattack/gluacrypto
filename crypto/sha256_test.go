package gluacrypto_crypto_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluacrypto"
	lua "github.com/yuin/gopher-lua"
)

func TestSha256(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluacrypto.Preload(L)

	h := sha256.New()
	h.Write(Data)
	hashData := h.Sum(nil)

	script := `
		crypto = require('crypto')
		return crypto.sha256('` + string(Data) + `')
	`
	assert.NoError(L.DoString(script))

	val := getValue(L, L.Get(1))
	err := getValue(L, L.Get(2))
	assert.Nil(err)
	assert.Equal(hex.EncodeToString(hashData), val)
}

func TestSha256Raw(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluacrypto.Preload(L)

	h := sha256.New()
	h.Write(Data)
	hashData := h.Sum(nil)

	script := `
		crypto = require('crypto')
		return crypto.sha256('` + string(Data) + `', true)
	`
	assert.NoError(L.DoString(script))

	val := getValue(L, L.Get(1))
	err := getValue(L, L.Get(2))
	assert.Nil(err)
	assert.Equal(string(hashData), val)
}
