package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func HmacMd5(s, key string) []byte {
	h := hmac.New(md5.New, []byte(key))
	_, _ = h.Write([]byte(s))
	return h.Sum(nil)
}

func HmacMd5String(s, key string) string {
	return hex.EncodeToString(HmacMd5(s, key))
}

func HmacSha1(s, key string) []byte {
	h := hmac.New(sha1.New, []byte(key))
	_, _ = h.Write([]byte(s))
	return h.Sum(nil)
}

func HmacSha1String(s, key string) string {
	return hex.EncodeToString(HmacSha1(s, key))
}

func HmacSha256(s, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(s))
	return h.Sum(nil)
}

func HmacSha256String(s, key string) string {
	return hex.EncodeToString(HmacSha256(s, key))
}
