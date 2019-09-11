package config

type Server struct {
	Hostname    string `yaml:"hostname"`
	Port        string `yaml:"port"`
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"privateKey"`
}
