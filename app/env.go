package libsrv

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/lithammer/shortuuid/v3"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
)

// Credential type, LOG,DB...
type Credential int

// Credential LOG,DB,...
const (
	LOG Credential = 0
	DB  Credential = 1
)

var logCred *google.Credentials
var dbCred *google.Credentials
var isProduction bool

//EnvJoinCurrentDir join dir with current dir
func EnvJoinCurrentDir(dir string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic("failed to call os.Getwd(), this should not happen")
	}
	return path.Join(currentDir, dir)
}

//EnvKeyPath get key path from name
//
//	keyPath, err := EnvKeyPath("log")
func EnvKeyPath(name string) (string, error) {
	keyPath := ""
	keyDir := "keys/"
	for i := 0; i < 5; i++ {
		keyPath = EnvJoinCurrentDir(keyDir + name + ".key")
		if _, err := os.Stat(keyPath); err == nil {
			return keyPath, nil
		}
		keyDir = "../" + keyDir
	}
	return "", errors.New("failed to find " + name + ".key in keys/ or ../keys/")
}

//EnvCheck check environment variable is set properly
func EnvCheck() {
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

//EnvProduction return true if is production environment
func EnvProduction() bool {
	return isProduction
}

//EnvPiyuoApp return environment variable PIYUO_APP
func EnvPiyuoApp() string {
	var env = os.Getenv("PIYUO_APP")
	if env == "" {
		panic("need  environment var PIYUO_APP")
	}
	return env
}

//EnvGoogleCredential return log or db google credential
func EnvGoogleCredential(ctx context.Context, c Credential) (*google.Credentials, error) {
	switch c {
	case LOG:
		if logCred == nil {
			cred, err := createGoogleCloudCredential(ctx, c)
			if err != nil {
				name, _ := getAttributesFromCredential(c)
				return nil, errors.Wrap(err, "failed to create log google cloud credential "+name+".key")
			}
			logCred = cred
		}
		return logCred, nil
	case DB:
		if dbCred == nil {
			cred, err := createGoogleCloudCredential(ctx, c)
			if err != nil {
				name, _ := getAttributesFromCredential(c)
				return nil, errors.Wrap(err, "failed to create db google cloud credential "+name+".key")
			}
			dbCred = cred
		}
		return dbCred, nil
	}
	panic("need add credential in switch above")
}

func createGoogleCloudCredential(ctx context.Context, c Credential) (*google.Credentials, error) {
	keyname, scope := getAttributesFromCredential(c)

	keyPath, err := EnvKeyPath(keyname)
	if err != nil {
		return nil, errors.Wrap(err, keyname+".key not found")
	}
	jsonfile, err := NewJSONFile(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "can no open key file "+keyPath)
	}
	defer jsonfile.Close()

	text, err := jsonfile.Text()
	if err != nil {
		return nil, errors.Wrap(err, " keyfile content maybe empty or wrong format. "+keyPath)
	}

	creds, err := google.CredentialsFromJSON(ctx, []byte(text), scope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert json to google credentials.\n"+text)
	}
	return creds, nil
}

// return filename and scope from credential
func getAttributesFromCredential(c Credential) (string, string) {
	filename := ""
	scope := ""
	switch c {
	case LOG:
		filename = "log"
		scope = "https://www.googleapis.com/auth/cloud-platform"
	case DB:
		filename = "db"
		scope = "https://www.googleapis.com/auth/datastore"
	}
	if filename == "" {
		panic("credential type not support type by getAttributesFromCredential(). " + string(c))
	}
	return filename, scope
}

//UUID generates concise, unambiguous, URL-safe UUIDs
func UUID() string {
	return shortuuid.New()
}
