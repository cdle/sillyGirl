package core

type Yaml struct {
	EnableRedis   bool   `yaml:"enable_redis"`
	RedisAddr     string `yaml:"redis_addr"`
	RedisPassword string `yaml:"redis_password"`
	SlaveMode     bool   `yaml:"slave_mode"`
	HttpPort      string `yaml:"http_port"`
}

var Config Yaml
