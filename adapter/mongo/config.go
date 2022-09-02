package mongo

import "fmt"

type (
	Config struct {
		Address string     `json:"address" yaml:"address"`
		Auth    AuthConfig `json:"auth" yaml:"auth"`
	}

	AuthConfig struct {
		Enable   bool   `json:"enable" yaml:"enable"`
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
		AuthDB   string `json:"auth_db" yaml:"auth_db"`
	}
)

var (
	DefaultBatchSize int32 = 1000
)

func (conf *Config) String() string {
	// Success
	return fmt.Sprintf("mongodb://%s", conf.Address)
}
