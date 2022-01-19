package gluacrypto_crypto_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tengattack/gluacrypto"
	crypto "github.com/tengattack/gluacrypto/crypto"
	lua "github.com/yuin/gopher-lua"
)

var (
	methods = []string{"des-ecb", "des-cbc", "aes-cbc"}
)

func TestEncrypt(t *testing.T) {
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
		  return crypto.encrypt('` + string(Data) + `', '` + method + `', '` + string(key) + `', 0, '` + string(iv) + `')
	  `
		assert.NoError(L.DoString(script))

		val := GetValue(L, L.Get(1))
		serr := GetValue(L, L.Get(2))
		assert.Nil(serr)
		assert.Equal(hex.EncodeToString(out), val)
	}
}

func TestEncryptRaw(t *testing.T) {
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
		  return crypto.encrypt('` + string(Data) + `', '` + method + `', '` + string(key) + `', crypto.RAW_DATA, '` + string(iv) + `')
	  `
		assert.NoError(L.DoString(script))

		val := GetValue(L, L.Get(1))
		serr := GetValue(L, L.Get(2))
		assert.Nil(serr)
		assert.Equal(string(out), val)
	}
}

func TestEncryptFail(t *testing.T) {
	assert := assert.New(t)

	L := lua.NewState()
	defer L.Close()
	gluacrypto.Preload(L)

	script := `
	  crypto = require('crypto')
		return crypto.encrypt('` + string(Data) + `', 'unknown', '` + string(Key) + `', 0, '` + string(IV) + `')
	`
	assert.NoError(L.DoString(script))

	val := GetValue(L, L.Get(1))
	err := GetValue(L, L.Get(2))
	assert.NotNil(err)
	assert.Nil(val)
}
