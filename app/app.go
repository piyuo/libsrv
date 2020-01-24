package app

import (
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// Credential LOG,DB,...
var isProduction bool

//JoinCurrentDir join dir with current dir
func JoinCurrentDir(dir string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic("failed to call os.Getwd(), this should not happen")
	}
	return path.Join(currentDir, dir)
}

//KeyPath get key path from name
//
//	keyPath, err := EnvKeyPath("log")
func KeyPath(name string) (string, error) {
	keyPath := ""
	keyDir := "keys/"
	for i := 0; i < 5; i++ {
		keyPath = JoinCurrentDir(keyDir + name + ".key")
		if _, err := os.Stat(keyPath); err == nil {
			return keyPath, nil
		}
		keyDir = "../" + keyDir
	}
	return "", errors.New("failed to find " + name + ".key in keys/ or ../keys/")
}

//Check check environment variable is set properly
func Check() {
	//id format like piyuo-tw-m-app
	id := os.Getenv("PIYUO_APP")
	if id == "" {
		panic("need set env var PIYUO_APP=...")
	}
	isProduction = false
	if strings.Contains(strings.ToLower(id), "-m-") {
		isProduction = true
	}
}

//IsProduction return true if is production environment
func IsProduction() bool {
	return isProduction
}

//PiyuoID return environment variable PIYUO_APP
func PiyuoID() string {
	var env = os.Getenv("PIYUO_APP")
	if env == "" {
		panic("need  environment var PIYUO_APP")
	}
	return env
}

var cryptoInstance Crypto

//Encrypt text using default crypto
//
//	cryped, err := app.Encrypt("hello")
func Encrypt(text string) (string, error) {
	if cryptoInstance == nil {
		cryptoInstance = NewCrypto()
	}
	return cryptoInstance.Encrypt(text)
}

//Decrypt text using default crypto
//
//	text, err := app.Decrypt(cryped)
func Decrypt(crypted string) (string, error) {
	if cryptoInstance == nil {
		cryptoInstance = NewCrypto()
	}
	return cryptoInstance.Decrypt(crypted)
}
