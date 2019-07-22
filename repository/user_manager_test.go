package repository

import (
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"os"
	"testing"
)

type UserDataProvider struct {
	manager  UserManager
	username string
	password string
	id       uuid.UUID
}

var (
	basicUMTest      = &UserManagerBasic{name: "basic"}
	leveldbUMTest    = &UserManagerLeveldb{cfgPath: "../test/user_cfg", namePath: "../test/user"}
	userDataProvider = []UserDataProvider{
		{basicUMTest, "testusername", "testuserpwd", uuid.FromStringOrNil("2490c31d-3005-47b4-9bc0-45952a2e505e")},
		{leveldbUMTest, "otherusername", "otheruserpwd", uuid.FromStringOrNil("68d0dffb-3dbf-4086-965f-33dd5d012995")},
	}
)

func executeUserTest(test func(provider UserDataProvider)) {
	os.RemoveAll("../test/")
	log, _ = zap.NewDevelopment()
	for _, provider := range userDataProvider {
		user := &User{
			UserName: provider.username,
			ID:       provider.id,
			Password: provider.password,
		}
		provider.manager.init()
		provider.manager.setUser(user)
		test(provider)
		provider.manager.close()
	}
}

func TestGetUser(t *testing.T) {
	test := func(provider UserDataProvider) {
		user, found := provider.manager.getUser(provider.username)
		if !found {
			t.Errorf("user %s not found", provider.username)
		}
		if user.ID != provider.id {
			t.Errorf("want %v got %v", provider.id, user.ID)
		}
	}
	executeUserTest(test)
}

func TestValidateUser(t *testing.T) {
	test := func(provider UserDataProvider) {
		err := provider.manager.validateUser(provider.username, provider.password)
		if err != nil {
			t.Errorf("want no error, got %s", err.Error())
		}
		err = provider.manager.validateUser(provider.username, "wrongpassword")
		if err == nil {
			t.Errorf("want login error, got no error")
		}
	}
	executeUserTest(test)
}

func TestDeleteUser(t *testing.T) {
	test := func(provider UserDataProvider) {
		err := provider.manager.validateUser(provider.username, provider.password)
		if err != nil {
			t.Errorf("want no error, got %s", err.Error())
		}
		err = provider.manager.validateUser(provider.username, "wrongpassword")
		if err == nil {
			t.Errorf("want login error, got no error")
		}
	}
	executeUserTest(test)
}
