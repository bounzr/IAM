package repository

import (
	"../utils"
	"github.com/gofrs/uuid"
	"testing"
	"os"
	"go.uber.org/zap"
)

type SessionDataProvider struct {
	manager    SessionManager
	repository string
	username   string
	password   string
	id         uuid.UUID
	token      [32]byte
}

var (
	basicSMTest         = &SessionManagerBasic{}
	leveldbSMTest       = &SessionManagerLeveldb{sessionsPath:"../test/session" }
	sessionDataProvider = []SessionDataProvider{
		{basicSMTest, "basic", "testusername", "testuserpwd", uuid.FromStringOrNil("2490c31d-3005-47b4-9bc0-45952a2e505e"), utils.GetRandom32Token()},
		{leveldbSMTest, "leveldb", "otherusername", "otheruserpwd", uuid.FromStringOrNil("68d0dffb-3dbf-4086-965f-33dd5d012995"), utils.GetRandom32Token()},
	}
)

func executeSessionTest(test func(provider SessionDataProvider)) {
	os.RemoveAll("../test/")
	log, _ = zap.NewDevelopment()
	for _, provider := range sessionDataProvider {
		user := &User{
			RepositoryName: provider.repository,
			UserName:       provider.username,
			ID:             provider.id,
			Password:       provider.password,
		}
		provider.manager.init()
		ctx := user.GetUserCtx()
		provider.manager.setSessionContext(provider.token, ctx)
		test(provider)
		provider.manager.close()
	}
}

func TestGetSession(t *testing.T) {
	test := func(provider SessionDataProvider) {
		ctx, err := provider.manager.getSessionContext(provider.token)
		if err != nil {
			t.Errorf("session not found - %s", err.Error())
		}
		if ctx.UserID != provider.id {
			t.Errorf("want %v got %v", provider.id, ctx.UserID)
		}
		err = provider.manager.deleteSessionContext(provider.token)
		if err != nil {
			t.Errorf("session not deleted - %s", err.Error())
		}
		/*
		ctx, err = provider.Manager.getSessionContext(provider.token)
		if err == nil {
			t.Errorf("session not deleted test failed in repository %s", provider.repository)
		} else {
			log.Info("the test returned a session error as expected", zap.Error(err))
		}
		*/
	}
	executeSessionTest(test)
}
