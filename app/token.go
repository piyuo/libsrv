package libsrv

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Token transfer between client and server, work like session in the old days
type Token interface {
	UserID() string

	StoreID() string

	LocationID() string

	// is this token expire? default expired in 1 hour
	Expired() bool

	// Revive token so user won't feel interrpt
	// return true if tooken are revive
	// only token created > 10 min and < 1 hour need revive
	Revive() bool

	// IsUserJustLogin return true if user just enter password to login in 3 minutes
	//
	//	 justLogin := token.IsUserJustLogin()
	IsUserJustLogin() bool

	// Identity like userId-storeId-locationId
	//
	//	 id := token.Identity()
	Identity() string

	//ToString save everything to string
	//
	//	 text := token.ToString()
	ToString() string

	//ToCookie save token to cookie
	//
	//	 err := token.ToCookie(w)
	ToCookie(w http.ResponseWriter) error

	//ToContext set token to context and return new context
	//
	//	 ctx := token.ToContext()
	ToContext(ctx context.Context) context.Context
}

//NewToken create new token
//
//	 token := NewToken("userId", "storeId", "locationId", "permission", time.Now())
func NewToken(userID, storeID, locationID, permission string, login, created time.Time) Token {
	token := &token{
		userID:     userID,
		storeID:    storeID,
		locationID: locationID,
		permission: permission,
		login:      login,
		created:    created,
	}
	return token
}

//TokenDeleteCookie delete token from cookie
func TokenDeleteCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   string(CookieTokenName),
		MaxAge: -1}
	http.SetCookie(w, &cookie)
}

//PiyuoTokenTimeLayout for cookie create time
const piyuoTimeLayout = "200601021504"

type name string

const (
	//ContextTokenName is name use in context to save token
	ContextTokenName name = "token"

	//CookieTokenName is access token name use in cookie
	CookieTokenName name = "piyuo"
)

type token struct {
	userID     string
	storeID    string
	locationID string
	permission string
	login      time.Time // the last time user enter password to login
	created    time.Time // cookie create time, may revive over time
}

func (s *token) Identity() string {
	if s.userID != "" && s.storeID != "" && s.locationID != "" {
		return s.userID + "-" + s.storeID + "-" + s.locationID

	} else if s.userID != "" && s.storeID != "" {
		return s.userID + "-" + s.storeID

	} else if s.userID != "" {
		return s.userID
	}
	return ""
}

func (s *token) UserID() string {
	return s.userID
}

func (s *token) StoreID() string {
	return s.storeID
}

func (s *token) LocationID() string {
	return s.locationID
}

func (s *token) ToString() string {
	createdText := s.created.Format(piyuoTimeLayout)
	data := s.userID + "|" + s.storeID + "|" + s.locationID + "|" + s.permission + "|" + createdText
	return data
}

//TokenFromString load everything from string
//
//	 token,err := TokenFromString()
func TokenFromString(str string) (Token, error) {
	arg := strings.Split(str, "|")
	if len(arg) != 5 {
		return nil, errors.New("failed to split data, it should have 5 arg ")
	}
	createdText := arg[4]
	created, err := time.Parse(piyuoTimeLayout, createdText)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode time in token.created")
	}
	token := &token{
		userID:     arg[0],
		storeID:    arg[1],
		locationID: arg[2],
		permission: arg[3],
		created:    created,
	}
	return token, nil
}

func (s *token) ToCookie(w http.ResponseWriter) error {
	everything := s.ToString()
	crypto := NewCrypto()
	crypted, err := crypto.Encrypt(everything)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt token("+everything+")")
	}
	cookie := http.Cookie{
		Name:     string(CookieTokenName),
		Value:    crypted,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   3600, // expire in one hour
	}
	http.SetCookie(w, &cookie)
	return nil
}

//TokenFromCookie get token from cookie
func TokenFromCookie(r *http.Request) (Token, error) {
	cookie, err := r.Cookie(string(CookieTokenName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get crypted text from cookie[\""+string(CookieTokenName)+"\"]")
	}
	crypto := NewCrypto()
	everything, err := crypto.Decrypt(cookie.Value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt text("+cookie.Value+")")
	}

	return TokenFromString(everything)
}

func (s *token) ToContext(ctx context.Context) context.Context {
	type favContextKey string
	return context.WithValue(ctx, ContextTokenName, s)
}

//TokenFromContext get token from context
func TokenFromContext(ctx context.Context) (Token, error) {
	value := ctx.Value(ContextTokenName)
	if value == nil {
		return nil, errors.New("failed to get token from context[\"" + string(ContextTokenName) + "\"]")
	}
	token := value.(Token)
	return token, nil
}

//check token is expired, token expire in 1 hour
func (s *token) Expired() bool {
	if s.created.IsZero() {
		return true
	}
	diff := time.Now().Sub(s.created)
	if diff.Minutes() > 60 {
		return true
	}
	return false
}

//Revive token, only token created > 10 min and < 1 hour need revive
func (s *token) Revive() bool {
	diff := time.Now().Sub(s.created)
	if diff.Minutes() > 10 && diff.Minutes() < 60 {
		s.created = time.Now()
		return true
	}
	return false
}

// IsUserJustLogin return true if user just enter password to login in 3 minutes
//
//	 justLogin := token.IsUserJustLogin()
func (s *token) IsUserJustLogin() bool {
	if s.login.IsZero() {
		return false
	}
	diff := time.Now().Sub(s.login)
	if diff.Minutes() > 3 {
		return false
	}
	return true
}
