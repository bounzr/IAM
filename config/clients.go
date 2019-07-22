package config

import "time"

type Clients struct {
	Implementation string `yaml:"implementation"`
	SecretDuration string `yaml:"secretDuration"`
}

func (c *Clients) GetSecretDuration() time.Duration {
	dur, err := time.ParseDuration(c.SecretDuration)
	if err == nil {
		return dur
	} else {
		return time.Hour
	}
}
