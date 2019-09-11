package repository

import (
	"../config"
	"../scim2"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type ClientManager interface {
	close()
	deleteClient(id uuid.UUID)
	findClients() ([]Client, error)
	getClient(id uuid.UUID) (client *Client, found bool)
	init()
	setClient(cli *Client)
}

func DeleteClient(clientID uuid.UUID) {
	clientManager.deleteClient(clientID)
}

func FindClients() []scim2.Client {
	var clients []scim2.Client
	repClients, err := clientManager.findClients()
	if err != nil {
		log.Error("can not add clients from repository", zap.Error(err))
	}
	for _, client := range repClients {
		clients = append(clients, *client.GetScim())
	}
	return clients
}

func GetClient(id interface{}) (client *Client, found bool) {
	if id == nil {
		return nil, false
	}
	return clientManager.getClient(id.(uuid.UUID))
}

func initClients() {
	implementation := config.IAM.Clients.Implementation
	switch implementation {
	case "leveldb":
		clientManager = &ClientManagerLeveldb{path: "./rep/client"}
	default:
		clientManager = &ClientManagerBasic{}
	}
	clientManager.init()
}

func ReplaceClientByScim(clientID uuid.UUID, scim *scim2.Client) error {
	client, found := GetClient(clientID)
	if !found {
		log.Debug("client not found", zap.String("id", clientID.String()))
		return ErrInvalidRequest
	}
	client.SetScim(scim)
	SetClient(client)
	return nil
}

func SetClient(cli *Client) {
	clientManager.setClient(cli)
	AddGroupResource(privateGroups["Clients"], cli.GetResourceTag())
}
