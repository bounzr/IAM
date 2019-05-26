package repository

import (
	"../utils"
	"github.com/gofrs/uuid"
)

type ConsentToken struct {
	ClientID uuid.UUID
	Token    string
}

func newConsentToken(clientID uuid.UUID) *ConsentToken {
	token := utils.GetRandomString(10)
	return &ConsentToken{clientID, token}
}
