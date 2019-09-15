package repository

import (
	"bounzr/iam/oauth2"
	"bounzr/iam/token"
	"bytes"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"os"
	"testing"
	"time"
)

type TokenDataProvider struct {
	manager         TokenManager
	clientID        uuid.UUID
	ownerID         uuid.UUID
	addRefreshToken bool
	scope           []byte
	state           string
	accessDuration  time.Duration
	refreshDuration time.Duration
	isValidToken    bool
}

var (
	basicTMTest       = &TokenManagerBasic{}
	levelTMTest       = &TokenManagerLeveldb{tokensPath: "../test/token", codesPath: "../test/code"}
	tokenDataProvider = []TokenDataProvider{
		{basicTMTest, uuid.FromStringOrNil("1490c31d-3005-47b4-9bc0-45952a2e5051"), uuid.FromStringOrNil("1490c31d-3005-47b4-9bc0-45952a2e5051"), true, []byte("s1 s2"), "teststate", time.Minute * 5, time.Minute * 10, true},
		{basicTMTest, uuid.FromStringOrNil("28d0dffb-3dbf-4086-965f-33dd5d012b92"), uuid.FromStringOrNil("28d0dffb-3dbf-4086-965f-33dd5d012b92"), false, []byte("s1 s2"), "teststate", time.Minute * 5, time.Minute * 10, true},
		//{basicTMTest, uuid.FromStringOrNil("3490c31b-3dbf-965f-4086-332a2e5129c3"), uuid.FromStringOrNil("3490c31b-3dbf-965f-4086-332a2e5129c3"), false, []byte(""), "", time.Minute * 0, time.Minute * 10, false},
		//{basicTMTest, uuid.FromStringOrNil("48d0dffb-9bc0-4086-3005-33dd52e506a4"), uuid.FromStringOrNil("48d0dffb-9bc0-4086-3005-33dd52e506a4"), false, []byte(""), "", time.Minute * 0, time.Minute * 0, false},
		{levelTMTest, uuid.FromStringOrNil("5490c31d-3005-47b4-9bc0-45952a2e5055"), uuid.FromStringOrNil("5490c31d-3005-47b4-9bc0-45952a2e5055"), true, []byte("s1 s2"), "teststate", time.Minute * 5, time.Minute * 10, true},
		{levelTMTest, uuid.FromStringOrNil("68d0dffb-3dbf-4086-965f-33dd5d012b96"), uuid.FromStringOrNil("68d0dffb-3dbf-4086-965f-33dd5d012b96"), false, []byte("s1 s2"), "teststate", time.Minute * 5, time.Minute * 10, true},
		//{levelTMTest, uuid.FromStringOrNil("7490c31b-3dbf-965f-4086-332a2e5129c7"), uuid.FromStringOrNil("7490c31b-3dbf-965f-4086-332a2e5129c7"), false, []byte(""), "", time.Minute * 0, time.Minute * 10, false},
		//{levelTMTest, uuid.FromStringOrNil("88d0dffb-9bc0-4086-3005-33dd52e506a8"), uuid.FromStringOrNil("88d0dffb-9bc0-4086-3005-33dd52e506a8"), false, []byte(""), "", time.Minute * 0, time.Minute * 0, false},
	}
)

func executeTokenTest(test func(provider TokenDataProvider, access *oauth2.TokenUnit)) {
	os.RemoveAll("../test/")
	log, _ = zap.NewDevelopment()
	token.Init()
	for _, provider := range tokenDataProvider {
		options := &oauth2.AccessTokenOptions{
			ClientID:        provider.clientID,
			AddRefreshToken: provider.addRefreshToken,
			Scope:           provider.scope,
			OwnerID:         provider.ownerID,
			State:           provider.state,
		}
		provider.manager.init()
		access, _ := oauth2.NewTokenSet(options, provider.accessDuration, provider.refreshDuration)
		log.Info("test access token", zap.String("client", provider.clientID.String()), zap.ByteString("token", access.GetToken()))
		provider.manager.setTokenUnit(access)
		test(provider, access)
		provider.manager.close()
	}
}

func TestValidateToken(t *testing.T) {
	test := func(provider TokenDataProvider, access *oauth2.TokenUnit) {
		accessHint := access.GetTokenHint()
		token, ok := provider.manager.getTokenUnit(accessHint)
		if provider.isValidToken != ok {
			t.Errorf("want valid=%v got valid=%v", provider.isValidToken, ok)
		}
		if token != nil && !ok {
			t.Errorf("got %s and got not ok", accessHint.Token)
		}
		if token == nil && ok {
			t.Errorf("got ok and got no token")
		}
		if ok {
			if token.GetToken() == nil {
				t.Errorf("want token got nil")
			}
			t1 := token.GetToken()
			t2 := access.GetToken()
			if bytes.Compare(t1, t2) != 0 {
				t.Errorf("want %v, got %v", t1, t2)
			}
		}
	}
	executeTokenTest(test)
}

//todo all token tests
