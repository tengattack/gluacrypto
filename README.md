# LuaCrypto for GopherLua

A native Go implementation of crypto library for the [GopherLua](https://github.com/yuin/gopher-lua) VM.

## Using

### Loading Modules

```go
import (
	"github.com/tengattack/gluacrypto"
)

// Bring up a GopherLua VM
L := lua.NewState()
defer L.Close()

// Preload LuaSocket modules
gluasocket.Preload(L)
```

### Encoding

* base64

### Hashing

* crc32
* md5
* sha1
* sha256
* sha512
* hmac

```lua
crypto.md5(input [, raw])
-- ...crc32, sha1, sha256, sha512
crypto.hmac(dtype, input, key [, raw])
```

If you need raw data output, set `raw` to `true`.

### Encrypt/Decrypt Chiper Method

* des-ecb
* des-cbc
* aes-cbc (key supports 128, 256, etc.)

```lua
crypto.encrypt(input, cipher, key, options, iv)
crypto.decrypt(input, cipher, key, options, iv)
```

If you need raw data input/output, using `crypto.RAW_DATA` as `options`, otherwise set it to 0.

## License

MIT
