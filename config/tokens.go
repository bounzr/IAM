package config

import (
	"time"
)

type Tokens struct{
	Implementation		string	`yaml:"implementation"`
	AccessDuration		string	`yaml:"accessDuration"`
	RefreshDuration		string	`yaml:"refreshDuration"`
}

func (t *Tokens) GetAccessDuration() time.Duration{
	dur, err := time.ParseDuration(t.AccessDuration)
	if err == nil{
		return dur
	}else{
		return time.Hour
	}
}

func (t *Tokens) GetRefreshDuration() time.Duration{
	dur, err := time.ParseDuration(t.RefreshDuration)
	if err == nil{
		return dur
	}else{
		return time.Hour * 24
	}
}