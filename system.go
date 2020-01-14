package libsrv

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"strings"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// System is system interface
type System interface {
	//Check system need environment variable is set
	Check()
	ID() string
	//IsProduction check system runing in production environment
	IsProduction() bool
	//Log text with serverity level
	Log(text string, level int32, piyuoID string, userID string)
	//Normal but significant events, such as start up, shut down, or a configuration change.
	Notice(text string)
	//Normal but significant events, such as start up, shut down, or a configuration change.
	NoticeBy(text string, userID string)
	//Warning events might cause problems.
	Warning(text string)
	//Warning events might cause problems.
	WarningBy(text string, userID string)
	//A person must take an action immediately.
	Alert(text string)
	//A person must take an action immediately.
	AlertBy(text string, userID string)
	//One or more systems are unusable.
	Emergency(text string)
	//Routine information, such as ongoing status or performance.
	Info(text string)
	//Routine information, such as ongoing status or performance.
	InfoBy(text string, userID string)
	//log Error
	Error(err error)
	//log Error
	//id is an identifier for the user,piyuoID affected by the error
	ErrorBy(err error, userID string)
	//stack format like "at firstLine (a.js:3)\nat secondLine (b.js:3)"
	//id is an identifier for the user,piyuoID affected by the error
	//language may be flutter, js, go, c#
	ErrorFrom(message string, stack string, piyuoID string, userID string, language string)
	JoinCurrentDir(dir string) string
	GetGoogleCloudCredential(c Credential) (*google.Credentials, error)
}

type system struct {
	googleCred   *google.Credentials
	isProduction bool
}

var instance System

//Sys is system only instance
//
//	Sys().Notify("hello")
func Sys() System {
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

// Credential LOG,DB,...
const (
	NOTICE    int32 = 1 //Normal but significant events, such as start up, shut down, or a configuration change.
	WARNING   int32 = 2 //Warning events might cause problems.
	CRITICAL  int32 = 3 //Critical events cause more severe problems or outages.
	ALERT     int32 = 4 //A person must take an action immediately.
	EMERGENCY int32 = 5 //One or more systems are unusable.
)

func (s *system) JoinCurrentDir(dir string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic("failed to call os.Getwd(), this should not happen")
	}
	return path.Join(currentDir, dir)
}

func (s *system) Check() {
	id := os.Getenv("PIYUO_ID")
	if id == "" {
		panic("need set env var PIYUO_ID=...")
	}
	//id format like piyuo-tw-m-app
	s.isProduction = s.checkProduction(id)
}

func (s *system) checkProduction(id string) bool {
	//id format like PIYUO-TW-M-SYS
	if strings.Contains(id, "-") {
		arg := strings.Split(id, "-")
		if arg[2] == "M" || arg[2] == "m" {
			return true
		}
	}
	return false
}

func (s *system) IsProduction() bool {
	return s.isProduction
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
	if _, err := os.Stat(keyfile); err != nil {
		keyfile = s.JoinCurrentDir("../keys/" + filename)
	}
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

// there is no error return for log
// server log like [PIYUO-TW-M-SYS] store-user: hello
// client log like (piyuo-tw-m-web-index) store-user: hello
func (s *system) Log(text string, level int32, piyuoID string, userID string) {
	log := s.getLogHead(piyuoID, userID) + ": " + text
	fmt.Printf("%v (logged)\n", log)
	ctx := context.Background()
	cred, err := s.GetGoogleCloudCredential(LOG)
	if err != nil {
		fmt.Printf("Log() failed to get google credential.  %v\n", err)
		return
	}

	client, err := logging.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		fmt.Printf("failed to create logging client: %v\n", err)
	}
	severity := logging.Notice
	switch level {
	case WARNING:
		severity = logging.Warning
	case CRITICAL:
		severity = logging.Critical
	case ALERT:
		severity = logging.Alert
	case EMERGENCY:
		severity = logging.Emergency
	}

	logger := client.Logger(piyuoID)
	logger.Log(logging.Entry{Payload: log, Severity: severity})

	if err := client.Close(); err != nil {
		fmt.Printf("failed to close client: %v\n", err)
	}
}

// GetLogHead use UPPER case for server, lower for client app
// server: [PIYUO-M-TW-AUTH] store-user: hello
// client: (piyuo-m-tw-web-page) store-user: hello
func (s *system) getLogHead(piyuoID string, userID string) string {
	displayID := userID
	if displayID != "" {
		displayID = " " + userID
	}
	if piyuoID[0] == 'P' {
		return fmt.Sprintf("[%v]%v", piyuoID, displayID)
	} else if piyuoID[0] == 'd' {
		return fmt.Sprintf("<%v>%v", piyuoID, displayID)
	}
	return fmt.Sprintf("(%v)%v", piyuoID, displayID)
}

func (s *system) Info(text string) {
	s.InfoBy(text, "")
}

func (s *system) InfoBy(text string, userID string) {
	fmt.Printf(s.getLogHead(s.ID(), userID) + ": " + text + "\n")
}

func (s *system) Notice(text string) {
	s.NoticeBy(text, "")
}

func (s *system) NoticeBy(text string, userID string) {
	s.Log(text, NOTICE, s.ID(), userID)
}

func (s *system) Warning(text string) {
	s.WarningBy(text, "")
}

func (s *system) WarningBy(text string, userID string) {
	s.Log(text, WARNING, s.ID(), userID)
}

func (s *system) Alert(text string) {
	s.AlertBy(text, "")
}

func (s *system) AlertBy(text string, userID string) {
	s.Log(text, ALERT, s.ID(), userID)
}

func (s *system) Emergency(text string) {
	s.Log(text, EMERGENCY, s.ID(), "")
}

func (s *system) Error(err error) {
	s.error(err, "", "", s.ID(), "", "")
}

func (s *system) ErrorBy(err error, userID string) {
	s.error(err, "", "", s.ID(), userID, "")
}

func (s *system) ErrorFrom(message string, stack string, piyuoID string, userID string, language string) {
	s.error(nil, message, stack, piyuoID, userID, language)
}

func (s *system) error(targetErr error, targetMessage string, targetStack string, targetPiyuoID string, targetUserID string, targetLanguage string) {
	if targetErr == nil && targetMessage == "" {
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
				fmt.Printf("could not log error: %v\n", err)
			},
		},
		option.WithCredentials(cred))
	if err != nil {
		fmt.Printf("failed to create error reporting client: %v\n", err)
	}
	defer client.Close()

	displayID := s.getLogHead(targetPiyuoID, targetUserID)
	if targetErr != nil {
		errWithID := errors.Wrap(targetErr, displayID)
		client.Report(errorreporting.Entry{
			Error: errWithID, User: displayID,
		})
		fmt.Printf(displayID+": %v\n", targetErr)
		stack := string(debug.Stack())
		fmt.Println(stack)
	} else {
		stack := s.formatStack(targetStack, targetLanguage)
		customErr := errors.New(displayID + ": " + targetMessage)
		client.Report(errorreporting.Entry{
			Error: customErr, User: displayID,
			Stack: []byte(stack),
		})
		fmt.Println(displayID+": ", targetMessage+"\n"+targetStack)
	}
}

func (s *system) formatStack(stack string, language string) string {
	if language == "flutter" {
		return stack
	}
	return stack
}
