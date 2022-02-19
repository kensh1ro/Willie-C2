package config

import (
	"encoding/base64"
	"crypto/rc4"
)

func Decrypt(data string) string {
	cipher, _ := rc4.NewCipher([]byte(RC4Key))
	ct, _ := base64.StdEncoding.DecodeString(data)
	dst := make([]byte, len(ct))
	cipher.XORKeyStream(dst, ct)
	return string(dst)
}