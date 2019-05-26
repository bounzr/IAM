package repository

import (
	"../oauth2"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
	"bytes"
	"encoding/gob"
)

type TokenRepositoryLeveldb struct {
	tokens *leveldb.DB
	codes  *leveldb.DB
}

func NewTokenRepositoryLeveldb() TokenRepository {
	repo := &TokenRepositoryLeveldb{}
	repo.init()
	return repo
}

func (r *TokenRepositoryLeveldb) init() {
	dbt, err := leveldb.OpenFile("./rep/token", nil)
	if err != nil {
		log.Error("can not init tokens repository", zap.Error(err))
	}
	r.tokens = dbt
	dbc, err := leveldb.OpenFile("./rep/code", nil)
	if err != nil {
		log.Error("can not init codes repository", zap.Error(err))
	}
	r.codes = dbc
	//todo defer db.Close()
}

func (r *TokenRepositoryLeveldb) deleteAccessToken(tokenHint *oauth2.AccessTokenHint) {
	err := r.tokens.Delete([]byte(tokenHint.Token), nil)
	if err != nil {
		log.Error("can not delete token", zap.Error(err))
	} else {
		log.Debug("deleted token")
	}
}

func (r *TokenRepositoryLeveldb) setAccessToken(accessToken *oauth2.AccessToken) error {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(accessToken)
	if err != nil {
		log.Error("can not encode access token", zap.String("user id", accessToken.OwnerID.String()), zap.String("client id", accessToken.ClientID.String()), zap.Error(err))
		return err
	}
	err = r.tokens.Put(accessToken.AccessToken, data.Bytes(), nil)
	if err != nil {
		log.Error("can not add token", zap.String("user id", accessToken.OwnerID.String()), zap.String("client id", accessToken.ClientID.String()), zap.Error(err))
		return err
	} else {
		log.Debug("token added", zap.String("user id", accessToken.OwnerID.String()), zap.String("client id", accessToken.ClientID.String()))
		return nil
	}
}

func (r *TokenRepositoryLeveldb) setAuthorizationCode(code *oauth2.AuthorizationCode) error {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(code)
	if err != nil {
		log.Error("can not encode authorization code", zap.String("user id", code.OwnerID.String()), zap.String("client id", code.ClientID.String()), zap.Error(err))
		return err
	}
	err = r.codes.Put([]byte(code.Code), data.Bytes(), nil)
	if err != nil {
		log.Error("can not add authorization code", zap.String("user id", code.OwnerID.String()), zap.String("client id", code.ClientID.String()), zap.Error(err))
		return err
	} else {
		log.Debug("authorization code added", zap.String("user id", code.OwnerID.String()), zap.String("client id", code.ClientID.String()))
		return nil
	}
}

func (r *TokenRepositoryLeveldb) validateAccessToken(tokenHint *oauth2.AccessTokenHint) (token *oauth2.AccessToken, ok bool) {
	dataBytes, err := r.tokens.Get([]byte(tokenHint.Token), nil)
	if err != nil {
		log.Error("can not find token", zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var aToken oauth2.AccessToken
	err = dec.Decode(&aToken)
	if err != nil {
		log.Error("can not decode token data", zap.Error(err))
		return nil, false
	}
	return &aToken, true
}

func (r *TokenRepositoryLeveldb) validateAuthorizationCode(request *oauth2.AuthorizationCodeAccessTokenRequest) (code *oauth2.AuthorizationCode, ok bool) {
	dataBytes, err := r.codes.Get([]byte(request.Code), nil)
	if err != nil {
		log.Error("can not find code", zap.String("client id", code.ClientID.String()), zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var rCode oauth2.AuthorizationCode
	err = dec.Decode(&rCode)
	if err != nil {
		log.Error("can not decode code data", zap.String("client id", code.ClientID.String()), zap.Error(err))
		return nil, false
	}
	return &rCode, true
}
