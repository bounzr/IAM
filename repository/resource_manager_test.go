package repository

import (
	"encoding/gob"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"os"
	"testing"
)

type ResourceDataProvider struct {
	Manager     ResourceManager
	TagProvider ResourceTagProvider
	ID          uuid.UUID
}

var (
	basicRMTest   = &ResourceManagerBasic{}
	leveldbRMTest = &ResourceManagerLeveldb{resourcesPath: "../test/resource"}
	userIDTest    = uuid.FromStringOrNil("2490c31d-3005-47b4-9bc0-45952a2e505e")
	userMDPTest   = NewUser(userIDTest, "testusername", "testpwd", "main")
	groupIDTest   = uuid.FromStringOrNil("3501d42e-4116-58c5-0cd1-56063b3a616a")
	groupMDPTest  = NewGroup(groupIDTest, "somegroup")
	clientIDTest  = uuid.FromStringOrNil("4601d42e-4517-59c4-1ba1-46063b3a617b")
	clientMDTest  = &Client{ID: clientIDTest}

	resourceDataProvider = []ResourceDataProvider{
		{basicRMTest, userMDPTest, userIDTest},
		{leveldbRMTest, userMDPTest, userIDTest},
		{basicRMTest, groupMDPTest, groupIDTest},
		{leveldbRMTest, groupMDPTest, groupIDTest},
		{basicRMTest, groupMDPTest, groupIDTest},
		{leveldbRMTest, groupMDPTest, groupIDTest},
	}
)

func executeResourceTest(test func(provider ResourceDataProvider)) {
	os.RemoveAll("../test/")
	log, _ = zap.NewDevelopment()
	for _, provider := range resourceDataProvider {
		provider.Manager.init()
		metadata := provider.TagProvider.GetResourceTag()
		log.Info("set resource metadata", zap.String("name", metadata.GetName()), zap.String("ID", metadata.GetUUID().String()))
		provider.Manager.setResourceTag(metadata)
		test(provider)
		provider.Manager.close()
	}
}

func TestLoadResourceMetadata(t *testing.T) {
	gob.Register(&ResourceTag{})
	test := func(provider ResourceDataProvider) {
		tag := &ResourceTag{}
		err := provider.Manager.loadResourceTag(provider.ID, tag)
		log.Debug("tag metadata returned", zap.String("name", tag.GetName()), zap.String("id", tag.GetUUID().String()), zap.Error(err))
		if err != nil {
			t.Errorf("tag %v not found", provider.ID)
		}
		if tag.GetUUID() != provider.ID {
			log.Debug("tag", zap.String("id", tag.GetUUID().String()))
			t.Errorf("want %v got %v", provider.ID, tag.GetUUID())
		}
	}
	executeResourceTest(test)
}
