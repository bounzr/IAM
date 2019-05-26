package repository

import (
	"../oauth2"
	"bytes"
	"encoding/gob"
	"github.com/gofrs/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

type ClientRepositoryLeveldb struct {
	db *leveldb.DB
}

func NewClientRepositoryLeveldb() ClientRepository {
	rep := &ClientRepositoryLeveldb{}
	rep.init()
	return rep
}

func (r *ClientRepositoryLeveldb) init() {
	db, err := leveldb.OpenFile("./rep/client", nil)
	if err != nil {
		log.Error("can not init client repository", zap.Error(err))
	}
	r.db = db
	//todo defer db.Close()
}

func (r *ClientRepositoryLeveldb) getClients() map[uuid.UUID]*oauth2.Client {
	clients := make(map[uuid.UUID]*oauth2.Client)
	iter := r.db.NewIterator(nil, nil)
	for iter.Next() {
		data := bytes.NewBuffer(iter.Value())
		dec := gob.NewDecoder(data)
		var client oauth2.Client
		err := dec.Decode(&client)
		if err == nil {
			clients[uuid.FromBytesOrNil(iter.Key())] = &client
		}
	}
	return clients
}

func (r *ClientRepositoryLeveldb) deleteClient(clientID uuid.UUID) {
	err := r.db.Delete(clientID.Bytes(), nil)
	if err != nil {
		log.Error("can not delete client", zap.Error(err))
	} else {
		log.Debug("deleted client", zap.String("client id", clientID.String()))
	}
}

func (r *ClientRepositoryLeveldb) getClient(clientID uuid.UUID) (*oauth2.Client, bool) {
	dataBytes, err := r.db.Get(clientID.Bytes(), nil)
	if err != nil {
		log.Error("can not get client", zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var client oauth2.Client
	err = dec.Decode(&client)
	if err != nil {
		log.Error("can not decode client", zap.Error(err))
		return nil, false
	}
	log.Debug("retrieved client", zap.String("client id", clientID.String()))
	return &client, true
}

func (r *ClientRepositoryLeveldb) updateClient(cli *oauth2.Client) {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(cli)
	if err != nil {
		log.Error("can not encode client", zap.Error(err))
	}
	err = r.db.Put(cli.GetClientID().Bytes(), data.Bytes(), nil)
	if err != nil {
		log.Error("can not update client", zap.Error(err))
	} else {
		log.Debug("updated client", zap.String("client id", cli.GetClientID().String()))
	}
}
