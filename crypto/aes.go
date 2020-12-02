package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
)

func AesEncrypt(plaintext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext = pkcs7padding(plaintext, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	blockModel.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func AesDecrypt(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origintext := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origintext, ciphertext)
	origintext = pkcs7unpadding(origintext)
	return origintext, nil
}

func AesQuickEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, block.BlockSize())
	rand.Read(iv)
	ciphertext, _ := AesEncrypt(plaintext, key, iv)
	datamap := map[string][]byte{"iv": iv, "value": ciphertext}
	data, err := json.Marshal(datamap)
	if err != nil {
		return nil, err
	}
	rst := base64.StdEncoding.EncodeToString(data)
	return []byte(rst), nil
}

func AesQuickDecrypt(cipher, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(cipher))
	if err != nil {
		return nil, errors.New("BASE64 解码错误")
	}
	x := make(map[string]string)
	json.Unmarshal(ciphertext, &x)
	iv, ok := x["iv"]
	if !ok {
		return nil, errors.New("无法确定iv")
	}
	value, _ := x["value"]
	ivDecode, _ := base64.StdEncoding.DecodeString(iv)
	vaDecode, _ := base64.StdEncoding.DecodeString(value)
	return AesDecrypt(vaDecode, key, ivDecode)
}


func pkcs7padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func pkcs7unpadding(plaintext []byte) []byte {
	length := len(plaintext)
	padding := int(plaintext[length-1])
	return plaintext[:(length - padding)]
}