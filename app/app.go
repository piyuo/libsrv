package app

import (
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	file "github.com/piyuo/libsrv/file"
	"github.com/pkg/errors"
)

//production -1 mean value not set
//0 debug
//1 in cloud
var production int8 = -1

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
//	keyPath, err := KeyPath("log")
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

//RegionKeyPath get region key path from name
//
//	keyPath, err := RegionKeyPath("log")
func RegionKeyPath(name string) (string, error) {
	return KeyPath("/regions/" + name)
}

//Key get key file content from name
//
//	key, err := Key("log")
func Key(name string) (string, error) {
	keyPath, err := KeyPath(name)
	if err != nil {
		return "", err
	}
	f, err := file.Open(keyPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return f.Text(), nil
}

//RegionKey get region key file content from name
//
//	regionKey, err := RegionKey("log")
func RegionKey(name string) (string, error) {
	return Key("/regions/" + name)
}

//Check only trigger in local debug environment
//
//	app.Check()
func Check() {
	//app format like piyuo-beta-sample-jp
	app := os.Getenv("PIYUO_APP")
	if app == "" {
		panic("need env PIYUO_APP like piyuo-beta-sample-jp")
	}
	//region format like piyuo-beta-sample-jp
	region := os.Getenv("PIYUO_REGION")
	if region == "" {
		panic("need env PIYUO_APP like us")
	}
	//slow warning, usually 12 seconds
	slow := os.Getenv("PIYUO_SLOW")
	if slow == "" {
		panic("need env PIYUO_SLOW like 12000")
	}
	//time to meet context deadline, this will stop all service, usually 20 seconds
	deadline := os.Getenv("PIYUO_DEADLINE")
	if deadline == "" {
		panic("need env PIYUO_DEADLINE like 20000")
	}
}

//IsDebug return true in local debug mode
func IsDebug() bool {
	if production == -1 {
		id := PiyuoID()
		if strings.Contains(id, "-m-") || strings.Contains(id, "-b-") || strings.Contains(id, "-a-") || strings.Contains(id, "-t-") || strings.Contains(id, "-sys-") {
			production = 1
		} else {
			production = 0
		}
	}
	return production == 0
}

//PiyuoID return environment variable PIYUO_APP
//
//	app.PiyuoID()
func PiyuoID() string {
	return os.Getenv("PIYUO_APP")
}

//ContextDateline get context deadline
//
//dateline should not greater than 10 min.
//
//	dateline,err := ContextDateline()
func ContextDateline() time.Time {
	text := os.Getenv("PIYUO_DEADLINE")
	ms, err := strconv.Atoi(text)
	if err != nil {
		panic("PIYUO_DEADLINE must be int")
	}
	return time.Now().Add(time.Duration(ms) * time.Millisecond)
}

//IsSlow check execution time is greater than slow definition,if so return slow limit, other return 0
//
//	So(IsSlow(5), ShouldBeFalse)
func IsSlow(executionTime int) int {
	text := os.Getenv("PIYUO_SLOW")
	ms, err := strconv.Atoi(text)
	if err != nil {
		panic("PIYUO_SLOW must be int")
	}
	if executionTime > ms {
		return ms
	}
	return 0
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
