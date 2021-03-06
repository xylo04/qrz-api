package qrz

import (
	"context"
	"errors"
	"github.com/antihax/optional"
)

const agent = "k0swe-go-1.0"

var cachedUser = ""
var cachedSession = ""

func Lookup(ctx context.Context, user *string, pw *string, call *string) (*QrzDatabase, error) {
	config := NewConfiguration()
	config.UserAgent = agent
	client := NewAPIClient(config)

	sessionKey, err := login(ctx, user, pw, client)
	if err != nil {
		return nil, err
	}

	lookupResp, err := lookupInner(ctx, sessionKey, call, client)
	if err != nil {
		cachedSession = ""
		return nil, err
	}
	return lookupResp, nil
}

func login(ctx context.Context, user *string, pw *string, client *APIClient) (string, error) {
	if cachedUser == *user && cachedSession != "" {
		return cachedSession, nil
	}
	req := new(RootGetOpts)
	req.Username = optional.NewString(*user)
	req.Password = optional.NewString(*pw)
	req.Agent = optional.NewString(agent)
	sessResp, _, err := client.DefaultApi.RootGet(ctx, req)
	if err != nil {
		return "", err
	}
	sessionKey := sessResp.Session.Key
	if sessionKey == "" {
		return "", errors.New(sessResp.Session.Error)
	}
	cachedUser = *user
	cachedSession = sessionKey
	return sessionKey, err
}

func lookupInner(ctx context.Context, sessionKey string, call *string, client *APIClient) (*QrzDatabase, error) {
	req := new(RootGetOpts)
	req.S = optional.NewString(sessionKey)
	req.Agent = optional.NewString(agent)
	req.Callsign = optional.NewString(*call)
	lookupResp, _, err := client.DefaultApi.RootGet(ctx, req)
	if err != nil {
		return nil, err
	}
	return &lookupResp, nil
}
