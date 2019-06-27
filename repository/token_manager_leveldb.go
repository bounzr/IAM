package repository

import (
	"../oauth2"
	"bytes"
	"encoding/gob"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
	"time"
)

type TokenManagerLeveldb struct {
	tokens     *leveldb.DB
	tokensPath string
	codes      *leveldb.DB
	codesPath  string
}

func (r *TokenManagerLeveldb) init() {
	if len(r.codesPath) == 0 {
		r.codesPath = "./rep/code"
	}
	if len(r.tokensPath) == 0 {
		r.tokensPath = "./rep/token"
	}
	dbt, err := leveldb.OpenFile(r.tokensPath, nil)
	if err != nil {
		log.Error("can not init token repository", zap.Error(err))
	}
	r.tokens = dbt
	dbc, err := leveldb.OpenFile(r.codesPath, nil)
	if err != nil {
		log.Error("can not init codes repository", zap.Error(err))
	}
	r.codes = dbc
}

func (r *TokenManagerLeveldb) close() {
	defer r.tokens.Close()
	defer r.codes.Close()
}

func (r *TokenManagerLeveldb) deleteAccessToken(tokenHint *oauth2.AccessTokenHint) {
	err := r.tokens.Delete([]byte(tokenHint.Token), nil)
	if err != nil {
		log.Error("can not delete token", zap.Error(err))
	} else {
		log.Debug("deleted token")
	}
}

func (r *TokenManagerLeveldb) setAccessToken(accessToken *oauth2.AccessToken) error {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(accessToken)
	if err != nil {
		log.Error("can not encode access token", zap.String("user ID", accessToken.OwnerID.String()), zap.String("client ID", accessToken.ClientID.String()), zap.Error(err))
		return err
	}
	err = r.tokens.Put(accessToken.AccessToken, data.Bytes(), nil)
	if err != nil {
		log.Error("can not add token", zap.String("user ID", accessToken.OwnerID.String()), zap.String("client ID", accessToken.ClientID.String()), zap.Error(err))
		return err
	} else {
		log.Debug("token added", zap.String("user ID", accessToken.OwnerID.String()), zap.String("client ID", accessToken.ClientID.String()))
		return nil
	}
}

func (r *TokenManagerLeveldb) setAuthorizationCode(code *oauth2.AuthorizationCode) error {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(code)
	if err != nil {
		log.Error("can not encode authorization code", zap.String("user ID", code.OwnerID.String()), zap.String("client ID", code.ClientID.String()), zap.Error(err))
		return err
	}
	err = r.codes.Put([]byte(code.Code), data.Bytes(), nil)
	if err != nil {
		log.Error("can not add authorization code", zap.String("user ID", code.OwnerID.String()), zap.String("client ID", code.ClientID.String()), zap.Error(err))
		return err
	} else {
		log.Debug("authorization code added", zap.String("user ID", code.OwnerID.String()), zap.String("client ID", code.ClientID.String()))
		return nil
	}
}

func (r *TokenManagerLeveldb) validateAccessToken(tokenHint *oauth2.AccessTokenHint) (token *oauth2.AccessToken, ok bool) {
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
	ok = aToken.GetExpirationTime() > time.Now().Unix()
	if !ok{
		r.deleteAccessToken(tokenHint)
		return nil, ok
	}
	return &aToken, ok
}

func (r *TokenManagerLeveldb) validateAuthorizationCode(request *oauth2.AuthorizationCodeAccessTokenRequest) (code *oauth2.AuthorizationCode, ok bool) {
	dataBytes, err := r.codes.Get([]byte(request.Code), nil)
	if err != nil {
		log.Error("can not find code", zap.String("client ID", code.ClientID.String()), zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var rCode oauth2.AuthorizationCode
	err = dec.Decode(&rCode)
	if err != nil {
		log.Error("can not decode code data", zap.String("client ID", code.ClientID.String()), zap.Error(err))
		return nil, false
	}
	return &rCode, true
}

