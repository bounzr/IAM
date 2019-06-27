package repository

import (
	"github.com/gofrs/uuid"
)

type ClientManagerBasic struct {
	clients map[uuid.UUID]*Client
}

func (r *ClientManagerBasic) close() {
	//nothing
}

func (r *ClientManagerBasic) deleteClient(id uuid.UUID) {
	delete(r.clients, id)
}

func (r *ClientManagerBasic) findClients() ([]Client, error) {
	clients := make([]Client, len(r.clients))
	idx := 0
	for _, client := range r.clients {
		clients[idx] = *client
		idx++
	}
	return clients, nil
}

func (r *ClientManagerBasic) init() {
	r.clients = make(map[uuid.UUID]*Client)
}

func (r *ClientManagerBasic) setClient(cli *Client) {
	r.clients[cli.ID] = cli
}

func (r *ClientManagerBasic) getClient(id uuid.UUID) (client *Client, found bool) {
	client, found = r.clients[id]
	return
}
