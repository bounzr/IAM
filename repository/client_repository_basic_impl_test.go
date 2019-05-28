package repository

import(
	"../oauth2"
	"testing"
	"github.com/gofrs/uuid"
)

func TestGetClient(t *testing.T) {
	uid, err := uuid.NewV4()
	if err != nil {
		t.Error(err)
	}
	cli := &oauth2.Client{
		ID:uid,
	}
	rep := NewClientRepositoryBasic()
	rep.updateClient(cli)
	client, found :=  rep.getClient(uid)
	if !found{
		t.Error("client not found")
	}
	if client.ID != uid{
		t.Error("client id is not matching the previously assigned client id")
	}
}
