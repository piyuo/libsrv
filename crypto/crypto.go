package app

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	key "github.com/piyuo/libsrv/key"

	"github.com/pkg/errors"
)

var blockSize int

// getBlock return cipher block from /keys/crypto.key, block will be cached after read from file
//
func getBlock() ([]byte, cipher.Block, error) {
	var err error
	cachedKey, err := key.Bytes("crypto.key") // key will cache key content
	if err != nil {
		return nil, nil, errors.Wrap(err, "/keys/crypto.key not found")
	}

	block, err := aes.NewCipher(cachedKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create cipher, make sure you have 128bit key in crypto.key, "+fmt.Sprintf("NewCipher(%d bytes) = %s",
			len(cachedKey), err))
	}
	if blockSize == 0 {
		blockSize = block.BlockSize()
	}
	return cachedKey, block, nil
}

// getEncrypter return encrypter
//
//	encrypter, err := crypto.Encrypter()
//
func getEncrypter() (cipher.BlockMode, error) {
	cryptoKey, block, err := getBlock()
	if err != nil {
		return nil, err
	}

	return cipher.NewCBCEncrypter(block, cryptoKey[:blockSize]), nil
}

// getDecrypter return decrypter
//
//	decrypter, err := crypto.Decrypter()
//
func getDecrypter() (cipher.BlockMode, error) {
	cryptoKey, block, err := getBlock()
	if err != nil {
		return nil, err
	}
	return cipher.NewCBCDecrypter(block, cryptoKey[:blockSize]), nil
}

// Encrypt string, if you need encrypt multiple string at the same time using getEncrypter()
//
//	crypted1, err := crypto.Encrypt("hello1")
//
func Encrypt(text string) (string, error) {
	encrypter, err := getEncrypter()
	if err != nil {
		return "", errors.Wrap(err, "failed to init encrypter")
	}

	textData := []byte(text)
	textData = pkcs7Padding(textData, blockSize)
	cryted := make([]byte, len(textData))
	encrypter.CryptBlocks(cryted, textData)
	return base64.StdEncoding.EncodeToString(cryted), nil
}

// Decrypt string, if you need decrypt multiple string at the same time using getEncrypter()
//
//	result, err := crypto.Decrypt(crypted)
//
func Decrypt(crypted string) (string, error) {
	decrypter, err := getDecrypter()
	if err != nil {
		return "", errors.Wrap(err, "failed to init decrypter")
	}

	crytedByte, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64 from cryted string")
	}
	if len(crytedByte) == 0 {
		return "", errors.New("crypted string can not be empty")
	}
	orig := make([]byte, len(crytedByte))
	decrypter.CryptBlocks(orig, crytedByte)
	orig = pkcs7UnPadding(orig)
	return string(orig), nil
}

func pkcs7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
