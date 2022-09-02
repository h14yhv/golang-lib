package rabbit

import "fmt"

type Config struct {
	Secure   bool   `json:"secure" yaml:"secure"`
	Address  string `json:"address" yaml:"address"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func (conf *Config) String() string {
	protocol := "amqp"
	if conf.Secure {
		protocol = "amqps"
	}
	// Success
	return fmt.Sprintf("%s://%s:%s@%s", protocol, conf.Username, conf.Password, conf.Address)
}
