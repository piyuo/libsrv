package libsrv

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime/debug"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// System is system interface
type System interface {
	Check()
	ID() string
	SetDebug(value bool)
	IsDebug() bool

	//Normal but significant events, such as start up, shut down, or a configuration change.
	Log(text string)

	//Warning events might cause problems.
	Warning(text string)

	//A person must take an action immediately.
	Alert(text string)

	//Routine information, such as ongoing status or performance.
	Info(text string)

	//log Error
	Error(err error)

	JoinCurrentDir(dir string) string

	GetGoogleCloudCredential(c Credential) (*google.Credentials, error)
}

type system struct {
	googleCred *google.Credentials
}

// DEBUG indicate debug or release mode
var DEBUG = true
var instance System

//CurrentSystem keep default system settings and provide log, error functions
//
//	System().Log("hello")
func CurrentSystem() System {
	if instance == nil {
		instance = &system{}
	}
	return instance
}

// Credential type, LOG,DB...
type Credential int

// Credential LOG,DB,...
const (
	LOG Credential = 0
	DB  Credential = 1
)

//LogLevel such
type LogLevel int

// Credential LOG,DB,...
const (
	NOTICE    LogLevel = 300 //Normal but significant events, such as start up, shut down, or a configuration change.
	WARNING   LogLevel = 400 //Warning events might cause problems.
	CRITICAL  LogLevel = 600 //Critical events cause more severe problems or outages.
	ALERT     LogLevel = 700 //A person must take an action immediately.
	EMERGENCY LogLevel = 800 //One or more systems are unusable.
)

func (s *system) JoinCurrentDir(dir string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic("failed to call os.Getwd(), this should not happen")
	}
	return path.Join(currentDir, dir)
}

func (s *system) Check() {
	if os.Getenv("PIYUO_ID") == "" {
		panic("need set env var PIYUO_ID=...")
	}
}

func (s *system) ID() string {
	return os.Getenv("PIYUO_ID")
}

// return filename and scope from credential
func (s *system) getAttributesFromCredential(c Credential) (string, string) {
	filename := ""
	scope := ""
	switch c {
	case LOG:
		filename = "log.key"
		scope = "https://www.googleapis.com/auth/cloud-platform"
	case DB:
		filename = "db.key"
		scope = "https://www.googleapis.com/auth/datastore"
	}
	if filename == "" {
		panic("credential type not support type by GoogleCloudCredentials(). " + string(c))
	}
	return filename, scope
}

// return key filename and scrope from credential
func (s *system) initGoogleCloudCredential(c Credential) (*google.Credentials, error) {
	filename, scope := s.getAttributesFromCredential(c)

	keyfile := s.JoinCurrentDir("keys/" + filename)

	jsonfile, err := NewJSONFile(keyfile)
	if err != nil {
		return nil, errors.Wrap(err, "can no open key file "+"keys/"+filename)
	}
	defer jsonfile.Close()

	text, err := jsonfile.Text()
	if err != nil {
		return nil, errors.Wrap(err, " keyfile content maybe empty or wrong format. "+"keys/"+filename)
	}

	ctx := context.Background()
	creds, err := google.CredentialsFromJSON(ctx, []byte(text), scope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert json to google credentials.\n"+text)
	}
	return creds, nil
}

// return key filename and scrope from credential
func (s *system) GetGoogleCloudCredential(c Credential) (*google.Credentials, error) {
	if s.googleCred == nil {
		cred, err := s.initGoogleCloudCredential(c)
		if err != nil {
			filename, _ := s.getAttributesFromCredential(c)
			return nil, errors.Wrap(err, "google cloud credential ini fail. "+filename)
		}
		s.googleCred = cred
	}
	return s.googleCred, nil
}

func (s *system) SetDebug(value bool) {
	DEBUG = value
}

func (s *system) IsDebug() bool {
	return DEBUG
}

// there is no error return for log
func (s *system) log(text string, loglevel LogLevel) {
	ctx := context.Background()
	cred, err := s.GetGoogleCloudCredential(LOG)
	if err != nil {
		fmt.Printf("Log() failed to get google credential.  %v", err)
		return
	}

	client, err := logging.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		fmt.Printf("failed to create logging client: %v", err)
	}
	severity := logging.Notice
	switch loglevel {
	case WARNING:
		severity = logging.Warning
	case CRITICAL:
		severity = logging.Critical
	case ALERT:
		severity = logging.Alert
	case EMERGENCY:
		severity = logging.Emergency
	}

	logger := client.Logger(s.ID())
	log := s.ID() + ": " + text
	logger.Log(logging.Entry{Payload: log, Severity: severity})
	fmt.Printf("%v (logged)", log)

	if err := client.Close(); err != nil {
		fmt.Printf("failed to close client: %v", err)
	}
}

func (s *system) Info(text string) {
	fmt.Printf(s.ID() + ": " + text)
}

func (s *system) Log(text string) {
	s.log(text, NOTICE)
}

func (s *system) Warning(text string) {
	s.log(text, WARNING)
}

func (s *system) Alert(text string) {
	s.log(text, ALERT)
}

func (s *system) Error(targetErr error) {
	if targetErr == nil {
		return
	}
	ctx := context.Background()
	cred, err := s.GetGoogleCloudCredential(LOG)
	if err != nil {
		fmt.Printf("Log() failed to get google credential.  %v\n", err)
		return
	}

	client, err := errorreporting.NewClient(ctx,
		cred.ProjectID,
		errorreporting.Config{
			ServiceName: s.ID(),
			OnError: func(err error) {
				fmt.Printf("could not log error: %v", err)
			},
		},
		option.WithCredentials(cred))
	if err != nil {
		fmt.Printf("failed to create error reporting client: %v", err)
	}
	defer client.Close()

	id := s.ID()
	client.Report(errorreporting.Entry{
		Error: targetErr, User: id,
	})
	fmt.Printf(s.ID()+"%v", err)
	stack := string(debug.Stack())
	fmt.Println(stack)
}
