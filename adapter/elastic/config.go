package elastic

type Config struct {
	Address string `json:"address" yaml:"address"`
}

func (conf *Config) String() string {
	// Success
	return conf.Address
}
