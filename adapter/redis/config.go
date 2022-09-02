package redis

type Config struct {
	Address  string `json:"address" yaml:"address"`
	Password string `json:"password" yaml:"password"`
	Db       int    `json:"db" yaml:"db"`
}
