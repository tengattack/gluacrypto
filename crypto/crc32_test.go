package gluacrypto_crypto_test

import (
	"encoding/hex"
	"hash/crc32"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluacrypto"
	"github.com/tengattack/tgo/luautil"
	lua "github.com/yuin/gopher-lua"
)

func TestCRC32(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluacrypto.Preload(L)

	h := crc32.NewIEEE()
	h.Write(Data)
	hashData := h.Sum(nil)

	script := `
		crypto = require('crypto')
		return crypto.crc32('` + string(Data) + `')
	`
	assert.NoError(L.DoString(script))

	val := luautil.GetValue(L, 1)
	err := luautil.GetValue(L, 2)
	assert.Nil(err)
	assert.Equal(hex.EncodeToString(hashData), val)
}

func TestCRC32Raw(t *testing.T) {
	assert := assert.New(t)

	// test start
	L := lua.NewState()
	defer L.Close()
	gluacrypto.Preload(L)

	h := crc32.NewIEEE()
	h.Write(Data)
	hashData := h.Sum(nil)

	script := `
		crypto = require('crypto')
		return crypto.crc32('` + string(Data) + `', true)
	`
	assert.NoError(L.DoString(script))

	val := luautil.GetValue(L, 1)
	err := luautil.GetValue(L, 2)
	assert.Nil(err)
	assert.Equal(string(hashData), val)
}
