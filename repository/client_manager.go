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

func AddClient(cli *Client) {
	clientManager.setClient(cli)
	AddResource(cli)
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

func GetClient(id uuid.UUID) (client *Client, found bool) {
	return clientManager.getClient(id)
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
