package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)
func RSAPKCS1Encrypt(plain, pk string) (string, error){
	block, _ := pem.Decode([]byte(pk))
	if block == nil {
		return "", errors.New("block is empty")
	}
	x, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	rst, err := rsa.EncryptPKCS1v15(rand.Reader, x, []byte(plain))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(rst), nil
}

func RSAPKCS1Decrypt(cipher, pk string) ([]byte, error){
	block, _ := pem.Decode([]byte(pk))
	if block == nil {
		return nil, errors.New("block is empty")
	}
	x, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ct, _ := base64.StdEncoding.DecodeString(cipher)
	rst, err := rsa.DecryptPKCS1v15(rand.Reader, x, ct)
	if err != nil {
		return nil, err
	}
	return rst, nil
}


func RSAPKCS8Decrypt(cipher, pk string) ([]byte, error) {
	block, _ := pem.Decode([]byte(pk))
	if block == nil {
		return nil, errors.New("block is empty")
	}
	x, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ct, _ := base64.StdEncoding.DecodeString(cipher)
	rst, err := rsa.DecryptPKCS1v15(rand.Reader, x.(*rsa.PrivateKey), ct)
	if err != nil {
		return nil, err
	}
	return rst, nil
}
