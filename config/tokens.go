package config

import (
	"time"
)

type Tokens struct {
	Implementation  string `yaml:"implementation"`
	AccessDuration  string `yaml:"accessDuration"`
	RefreshDuration string `yaml:"refreshDuration"`
}

func (t *Tokens) GetAccessDuration() time.Duration {
	dur, err := time.ParseDuration(t.AccessDuration)
	if err == nil {
		return dur
	} else {
		log.Error("can not parse access_token duration from config. 1 hour will be used")
		return time.Hour
	}
}

func (t *Tokens) GetRefreshDuration() time.Duration {
	dur, err := time.ParseDuration(t.RefreshDuration)
	if err == nil {
		return dur
	} else {
		log.Error("can not parse refresh_token duration from config. 24 hours will be used")
		return time.Hour * 24
	}
}
