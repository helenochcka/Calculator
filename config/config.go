package config

type Config struct {
	ExecutorGinServer struct {
		Address string
		Port    string
	} `yaml:"executor_gin_server"`

	ExecutorGRPCServer struct {
		Address string
		Port    string
	} `yaml:"executor_grpc_server"`

	ArithmeticServer struct {
		Address string
		Port    string
	} `yaml:"arithmetic_server"`

	RabbitMQBroker struct {
		URI         string `yaml:"uri"`
		ContentType string `yaml:"content_type"`
	} `yaml:"rabbitmq_broker"`
}
