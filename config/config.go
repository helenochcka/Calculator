package config

type Config struct {
	ExecutorServer struct {
		Address string
		Port    string
	} `yaml:"executor_server"`

	ArithmeticServer struct {
		Address string
		Port    string
	} `yaml:"arithmetic_server"`

	RabbitMQBroker struct {
		URI         string `yaml:"uri"`
		ContentType string `yaml:"content_type"`
	} `yaml:"rabbitmq_broker"`
}
