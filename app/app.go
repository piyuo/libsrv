package app

import (
	"os"
	"path"
	"strconv"
	"time"

	file "github.com/piyuo/libsrv/file"
	"github.com/pkg/errors"
)

// production -1 mean value not set
// 0 debug
// 1 in cloud
//
var production int8 = -1

//JoinCurrentDir join dir with current dir
//
func JoinCurrentDir(dir string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(currentDir, dir), nil
}

//KeyPath get key real path from name, key path is "key/" which can be place under /src or /src/project
//
//	keyPath, err := KeyPath("log")
func KeyPath(name string) (string, error) {
	keyDir := "keys/"
	for i := 0; i < 5; i++ {
		keyPath, err := JoinCurrentDir(keyDir + name + ".json")
		if err != nil {
			return "", errors.Wrap(err, "failed to make KeyPath: "+name)
		}
		if _, err = os.Stat(keyPath); err == nil {
			return keyPath, nil
		}
		keyDir = "../" + keyDir
	}
	return "", errors.New("failed to find " + name + ".json in keys/ or ../keys/")
}

//RegionKeyPath get region key path from name
//
//	keyPath, err := RegionKeyPath("log")
func RegionKeyPath(name string) (string, error) {
	return KeyPath("/region/" + name)
}

//Key get key file content from name
//
//	key, err := Key("log")
func Key(name string) (string, error) {
	keyPath, err := KeyPath(name)
	if err != nil {
		return "", err
	}
	return file.ReadText(keyPath)
}

//RegionKey get region key file content from name
//
//	regionKey, err := RegionKey("log")
func RegionKey(name string) (string, error) {
	return Key("/region/" + name)
}

//Check only trigger in local debug environment
//
//	app.Check()
func Check() {
	//app format like piyuo-beta-sample-jp
	app := os.Getenv("PIYUO_APP")
	if app == "" {
		panic("need env PIYUO_APP=\"sample-jp\"")
	}
	//region format like piyuo-beta-sample-jp
	region := os.Getenv("PIYUO_REGION")
	if region == "" {
		panic("need env PIYUO_APP=\"us\"")
	}
	//slow warning, usually 12 seconds
	slow := os.Getenv("PIYUO_SLOW")
	if slow == "" {
		panic("need env PIYUO_SLOW=\"12000\"")
	}
	//time to meet context deadline, this will stop all service, usually 20 seconds
	deadline := os.Getenv("PIYUO_DEADLINE")
	if deadline == "" {
		panic("need env PIYUO_DEADLINE=\"20000\"")
	}
}

//stop support IsDebug
/*
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
*/

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
