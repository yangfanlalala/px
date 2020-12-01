package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func HmacMd5(s, key string) string {
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha1(s, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha256(s, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
