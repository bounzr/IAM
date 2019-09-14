package repository

import (
	"bounzr/iam/oauth2"
	"bytes"
	"encoding/gob"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
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

func (r *TokenManagerLeveldb) setTokenUnit(accessToken *oauth2.TokenUnit) error {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(accessToken)
	if err != nil {
		log.Error("can not encode access token", zap.String("user ID", accessToken.OwnerID.String()), zap.String("client ID", accessToken.ClientID.String()), zap.Error(err))
		return err
	}
	err = r.tokens.Put(accessToken.Token, data.Bytes(), nil)
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

func (r *TokenManagerLeveldb) getTokenUnit(tokenHint *oauth2.AccessTokenHint) (token *oauth2.TokenUnit, ok bool) {
	if tokenHint == nil {
		log.Error("token hint is nil")
		return nil, false
	}
	dataBytes, err := r.tokens.Get([]byte(tokenHint.Token), nil)
	if err != nil {
		log.Error("can not find token", zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var aToken oauth2.TokenUnit
	err = dec.Decode(&aToken)
	if err != nil {
		log.Error("can not decode token data", zap.Error(err))
		return nil, false
	}
	return &aToken, true
}

func (r *TokenManagerLeveldb) validateAuthorizationCode(request *oauth2.AuthorizationCodeAccessTokenRequest) (code *oauth2.AuthorizationCode, ok bool) {
	dataBytes, err := r.codes.Get([]byte(request.Code), nil)
	if err != nil {
		log.Error("can not find code", zap.String("client ID", request.ClientID), zap.Error(err))
		return nil, false
	}
	data := bytes.NewBuffer(dataBytes)
	dec := gob.NewDecoder(data)
	var rCode oauth2.AuthorizationCode
	err = dec.Decode(&rCode)
	if err != nil {
		log.Error("can not decode code data", zap.String("client ID", request.ClientID), zap.Error(err))
		return nil, false
	}
	//remove code code after use
	err = r.codes.Delete([]byte(request.Code), nil)
	if err != nil {
		log.Error("can not removed consumed code", zap.Error(err))
	}
	return &rCode, true
}
