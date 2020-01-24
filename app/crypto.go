package app

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

//Crypto interface
type Crypto interface {
	// encrypt string
	//
	//	crypted, err := crypto.Encrypt("hello")
	Encrypt(text string) (string, error)

	// decrypt string
	//
	//	result, err := crypto.Decrypt(crypted)
	Decrypt(crypted string) (string, error)
}

//NewCrypto create a crypto
//
//	crypto := NewCrypto()
//	crypted, err := crypto.Encrypt("hello")
func NewCrypto() Crypto {
	return &crypto{}
}

type crypto struct {
	key          []byte
	block        cipher.Block
	blockSize    int
	encBlockMode cipher.BlockMode
	decBlockMode cipher.BlockMode
}

//use 128 bit aes key generate by online Encryption Key Generator
func (s *crypto) initCryptoKey() error {
	if s.key != nil {
		return nil
	}
	keyPath, err := KeyPath("crypto")
	if err != nil {
		return errors.Wrap(err, "crypto.key not found")
	}

	file, err := os.Open(keyPath)
	if err != nil {
		return errors.Wrap(err, "failed to open key file "+keyPath)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "failed to read key file "+keyPath)
	}
	s.key = bytes
	return nil
}

func (s *crypto) initCipher() error {
	if s.block != nil {
		return nil
	}

	err := s.initCryptoKey()
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return errors.Wrap(err, "failed to create cipher, make sure you have 128bit key in crypto.key, "+fmt.Sprintf("NewCipher(%d bytes) = %s", len(s.key), err))
	}
	s.block = block
	s.blockSize = block.BlockSize()
	return nil
}

func (s *crypto) encryptBlockMode() cipher.BlockMode {
	if s.encBlockMode == nil {
		s.encBlockMode = cipher.NewCBCEncrypter(s.block, s.key[:s.blockSize])
	}
	return s.encBlockMode
}

func (s *crypto) decryptBlockMode() cipher.BlockMode {
	if s.decBlockMode == nil {
		s.decBlockMode = cipher.NewCBCDecrypter(s.block, s.key[:s.blockSize])
	}
	return s.decBlockMode
}

func (s *crypto) Encrypt(text string) (string, error) {
	err := s.initCipher()
	if err != nil {
		return "", errors.Wrap(err, "failed to init cipher")
	}
	textData := []byte(text)
	textData = pkcs7Padding(textData, s.blockSize)
	cryted := make([]byte, len(textData))
	s.encryptBlockMode().CryptBlocks(cryted, textData)
	return base64.StdEncoding.EncodeToString(cryted), nil
}

//EnvDecrypt decrypt string use crypto.key
func (s *crypto) Decrypt(crypted string) (string, error) {
	err := s.initCipher()
	if err != nil {
		return "", errors.Wrap(err, "failed to init cipher")
	}
	crytedByte, err := base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64 from cryted string")
	}
	if len(crytedByte) == 0 {
		return "", errors.New("crypted string can not be empty")
	}
	orig := make([]byte, len(crytedByte))
	s.decryptBlockMode().CryptBlocks(orig, crytedByte)
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
