package log

import (
	"context"
	"net/http"
)

// keyContext define key used in ctx
//
type keyContext int

const (
	// REQUEST is context key name for request
	//
	keyRequest keyContext = iota

	// TOKEN is context key name for token
	//
	keyToken
)

// getID return current command context id where log happen
//
//	id := getID(ctx) // user-store
//
func getID(ctx context.Context) string {
	id := ""
	value := ctx.Value(keyToken) //get token from command ctx
	if value != nil {
		m := value.(map[string]string)
		id = m["id"]
	}
	return id
}

// getRequest return current command context request where log happen
//
//	req := getRequest(ctx) // user-store
//
func getRequest(ctx context.Context) *http.Request {
	req := ctx.Value(keyRequest)
	if req == nil {
		return nil
	}
	return req.(*http.Request)
}
