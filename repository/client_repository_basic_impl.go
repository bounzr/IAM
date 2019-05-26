package repository

import (
	"../oauth2"
	"github.com/gofrs/uuid"
)

type ClientRepositoryBasic struct {
	clients map[uuid.UUID]*oauth2.Client
}

func NewClientRepositoryBasic() ClientRepository {
	mcr := &ClientRepositoryBasic{}
	mcr.init()
	return mcr
}

func (r *ClientRepositoryBasic) getClients() map[uuid.UUID]*oauth2.Client {
	return r.clients
}

func (r *ClientRepositoryBasic) init() {
	r.clients = make(map[uuid.UUID]*oauth2.Client)
}

func (r *ClientRepositoryBasic) updateClient(cli *oauth2.Client) {
	clientId := cli.GetClientID()
	r.clients[clientId] = cli
}

func (r *ClientRepositoryBasic) deleteClient(id uuid.UUID) {
	delete(r.clients, id)
}

func (r *ClientRepositoryBasic) getClient(id uuid.UUID) (client *oauth2.Client, found bool) {
	client, found = r.clients[id]
	return
}
