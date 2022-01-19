package gluacrypto_crypto_test

import (
	"crypto/md5"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluacrypto"
	lua "github.com/yuin/gopher-lua"
)

func TestMD5(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluacrypto.Preload(L)

	h := md5.New()
	h.Write(Data)
	hashData := h.Sum(nil)

	script := `
		crypto = require('crypto')
		return crypto.md5('` + string(Data) + `')
	`
	assert.NoError(L.DoString(script))

	val := GetValue(L, L.Get(1))
	err := GetValue(L, L.Get(2))
	assert.Nil(err)
	assert.Equal(hex.EncodeToString(hashData), val)
}

func TestMD5Raw(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluacrypto.Preload(L)

	h := md5.New()
	h.Write(Data)
	hashData := h.Sum(nil)

	script := `
		crypto = require('crypto')
		return crypto.md5('` + string(Data) + `', true)
	`
	assert.NoError(L.DoString(script))

	val := GetValue(L, L.Get(1))
	err := GetValue(L, L.Get(2))
	assert.Nil(err)
	assert.Equal(string(hashData), val)
}
