package repository

import (
	"bytes"
	"encoding/gob"
	"github.com/gofrs/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

type ClientManagerLeveldb struct {
	db   *leveldb.DB
	path string
}

func (r *ClientManagerLeveldb) close() {
	defer r.db.Close()
}

func (r *ClientManagerLeveldb) init() {
	if len(r.path) == 0 {
		r.path = "./rep/client"
	}
	db, err := leveldb.OpenFile(r.path, nil)
	if err != nil {
		log.Error("can not init client repository", zap.Error(err))
	}
	r.db = db
}

/*
func (r *ClientManagerLeveldb) getClients() map[uuid.UUID]*oauth2.Client {
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
*/
func (r *ClientManagerLeveldb) deleteClient(clientID uuid.UUID) {
	err := r.db.Delete(clientID.Bytes(), nil)
	if err != nil {
		log.Error("can not delete client", zap.Error(err))
	} else {
		log.Debug("deleted client", zap.String("client ID", clientID.String()))
	}
}

func (r *ClientManagerLeveldb) getClient(ID uuid.UUID) (*Client, bool) {
	dataBytes, err := r.db.Get(ID.Bytes(), nil)
	if err != nil {
		log.Error("can not get client", zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var client Client
	err = dec.Decode(&client)
	if err != nil {
		log.Error("can not decode client", zap.Error(err))
		return nil, false
	}
	log.Debug("retrieved client", zap.String("client ID", ID.String()))
	return &client, true
}

//todo findClients(search parameters)
func (r *ClientManagerLeveldb) findClients() ([]Client, error) {
	var clients []Client
	iter := r.db.NewIterator(nil, nil)
	for iter.Next() {
		dataBytes := iter.Value()
		data := bytes.NewBuffer(dataBytes)
		dec := gob.NewDecoder(data)
		var client Client
		err := dec.Decode(&client)
		if err != nil {
			log.Error("can not decode client", zap.ByteString("id", iter.Key()), zap.Error(err))
		}
		log.Debug("client decoded", zap.ByteString("id", iter.Key()))
		clients = append(clients, client)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *ClientManagerLeveldb) setClient(cli *Client) {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(cli)
	if err != nil {
		log.Error("can not encode client", zap.Error(err))
	}
	err = r.db.Put(cli.ID.Bytes(), data.Bytes(), nil)
	if err != nil {
		log.Error("can not update client", zap.Error(err))
	} else {
		log.Debug("updated client", zap.String("client ID", cli.ID.String()))
	}
}
