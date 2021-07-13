package gluacrypto_crypto_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluacrypto"
	crypto "github.com/tengattack/gluacrypto/crypto"
	luautil "github.com/tengattack/gluacrypto/util"
	lua "github.com/yuin/gopher-lua"
)

func encodeRawDataToLuaString(data []byte) string {
	out := make([]byte, 0, len(data)*4)
	for _, b := range data {
		if (b >= 'A' && b <= 'Z') ||
			(b >= 'a' && b <= 'z') {
			out = append(out, b)
		} else {
			out = append(out, []byte(fmt.Sprintf("\\%d", int(b)))...)
		}
	}
	return string(out)
}

func TestDecrypt(t *testing.T) {
	assert := assert.New(t)

	// test start
	for _, method := range methods {
		L := lua.NewState()
		defer L.Close()
		gluacrypto.Preload(L)

		var key []byte
		var iv []byte
		if strings.HasPrefix(method, "des-") {
			key = Key
			iv = IV
		} else {
			key = Key128
			iv = IV128
		}
		out, err := crypto.Encrypt(Data, method, key, iv)
		assert.NoError(err)

		script := `
		  crypto = require('crypto')
		  return crypto.decrypt('` + hex.EncodeToString(out) + `', '` + method + `', '` + string(key) + `', 0, '` + string(iv) + `')
	  `
		assert.NoError(L.DoString(script))

		val := luautil.GetValue(L, L.Get(1))
		serr := luautil.GetValue(L, L.Get(2))
		assert.Nil(serr)
		assert.Equal(string(Data), val)
	}
}

func TestDecryptRaw(t *testing.T) {
	assert := assert.New(t)

	// test start
	for _, method := range methods {
		L := lua.NewState()
		defer L.Close()
		gluacrypto.Preload(L)

		var key []byte
		var iv []byte
		if strings.HasPrefix(method, "des-") {
			key = Key
			iv = IV
		} else {
			key = Key128
			iv = IV128
		}
		out, err := crypto.Encrypt(Data, method, key, iv)
		assert.NoError(err)

		script := `
		  crypto = require('crypto')
		  return crypto.decrypt('` + encodeRawDataToLuaString(out) + `', '` + method + `', '` + string(key) + `', crypto.RAW_DATA, '` + string(iv) + `')
	  `
		assert.NoError(L.DoString(script))

		val := luautil.GetValue(L, L.Get(1))
		serr := luautil.GetValue(L, L.Get(2))
		assert.Nil(serr)
		assert.Equal(string(Data), val)
	}
}

func TestDecryptFail(t *testing.T) {
	assert := assert.New(t)

	L := lua.NewState()
	defer L.Close()
	gluacrypto.Preload(L)

	script := `
	  crypto = require('crypto')
		return crypto.decrypt('` + string(Data) + `', 'unknown', '` + string(Key) + `', 0, '` + string(IV) + `')
	`
	assert.NoError(L.DoString(script))

	val := luautil.GetValue(L, L.Get(1))
	err := luautil.GetValue(L, L.Get(2))
	assert.NotNil(err)
	assert.Nil(val)
}
