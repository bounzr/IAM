package repository

import (
	"../oauth2"
	"github.com/gofrs/uuid"
)

type ClientRepository interface {
	init()
	updateClient(cli *oauth2.Client)
	deleteClient(id uuid.UUID)
	getClient(id uuid.UUID) (client *oauth2.Client, found bool)
	getClients() map[uuid.UUID]*oauth2.Client
}

func AddClient(cli *oauth2.Client) {
	clientManager.updateClient(cli)
}

func GetClient(id uuid.UUID) (client *oauth2.Client, found bool) {
	return clientManager.getClient(id)
}

func GetClients() []*oauth2.ClientSummary {
	var clientsList []*oauth2.ClientSummary
	clients := clientManager.getClients()
	clientsList = make([]*oauth2.ClientSummary, len(clients))

	i := 0
	for _, cli := range clients {
		newCS := &oauth2.ClientSummary{
			Groups: cli.GetGroups(),
			ID:     cli.GetClientID().String(),
			Name:   cli.GetClientName(),
			Owner:  cli.GetOwnerID().String(),
		}
		clientsList[i] = newCS
		i++
	}
	return clientsList
}
