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
		Host    string
	} `yaml:"arithmetic_server"`

	RabbitMQBroker struct {
		URI string `yaml:"uri"`
	} `yaml:"rabbitmq_broker"`
}
