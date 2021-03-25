package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"strconv"

	key "github.com/piyuo/libsrv/src/key"
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
		return nil, nil, errors.Wrapf(err, "must use 128bit key in crypto.key, got %d ", len(cachedKey))
	}
	if blockSize == 0 {
		blockSize = block.BlockSize()
	}
	return cachedKey, block, nil
}

// newEncrypter return encrypter
//
//	encrypter, err := newEncrypter()
//
func newEncrypter() (cipher.BlockMode, error) {
	cryptoKey, block, err := getBlock()
	if err != nil {
		return nil, err
	}
	return cipher.NewCBCEncrypter(block, cryptoKey[:blockSize]), nil
}

// newDecrypter return decrypter
//
//	decrypter, err := newDecrypter()
//
func newDecrypter() (cipher.BlockMode, error) {
	cryptoKey, block, err := getBlock()
	if err != nil {
		return nil, err
	}
	return cipher.NewCBCDecrypter(block, cryptoKey[:blockSize]), nil
}

// Encrypt string, if you need encrypt multiple string at the same time using getEncrypter()
//
//	crypted1, err := Encrypt("hello1")
//
func Encrypt(text string) (string, error) {
	encrypter, err := newEncrypter()
	if err != nil {
		return "", errors.Wrap(err, "new encrypter")
	}

	textData := []byte(text)
	textData = pkcs7Padding(textData, blockSize)
	cryted := make([]byte, len(textData))
	encrypter.CryptBlocks(cryted, textData)
	return base64.RawStdEncoding.EncodeToString(cryted), nil
}

// Decrypt string, if you need decrypt multiple string at the same time using getEncrypter()
//
//	result, err := cDecrypt(crypted)
//
func Decrypt(crypted string) (string, error) {
	decrypter, err := newDecrypter()
	if err != nil {
		return "", errors.Wrap(err, "new decrypter")
	}

	crytedByte, err := base64.RawStdEncoding.DecodeString(crypted)
	if err != nil {
		return "", errors.Wrap(err, "decode base64")
	}
	lenCryted := len(crytedByte)
	if lenCryted == 0 {
		return "", errors.New("input must not empty")
	}
	if lenCryted%blockSize != 0 {
		return "", errors.New("input not full blocks, block size must be " + strconv.Itoa(blockSize))
	}

	orig := make([]byte, lenCryted)
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
