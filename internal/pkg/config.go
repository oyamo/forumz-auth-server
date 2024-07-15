package pkg

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseDSN         string
	P12Certificate      string
	PublicKey           string
	CertPassword        string
	KafkaConsumerServer string
	KafkaProducerServer string
	RedisServer         string
}

const (
	envDatabaseDSN    = "AUTH_SERVICE_DATABASE_DSN"
	envP12Certificate = "AUTH_SERVICE_P12_CERTIFICATE"
	envPublicKey      = "AUTH_SERVICE_PUBLIC_KEY"
	envCertPassword   = "AUTH_SERVICE_CERT_PASSWORD"
	envKafkaConsumer  = "AUTH_SERVICE_KAFKA_CONSUMER"
	envKafkaProducer  = "AUTH_SERVICE_KAFKA_PRODUCER"
	envRedisServer    = "AUTH_SERVICE_REDIS_SERVER"
)

func EnvNotSetError(env string) error {
	return fmt.Errorf("%s environment variable not set", env)
}

func NewConfig() (*Config, error) {
	dbDSN, ok := os.LookupEnv(envDatabaseDSN)
	if !ok {
		return nil, EnvNotSetError(envDatabaseDSN)
	}

	p12Cert, ok := os.LookupEnv(envP12Certificate)
	if !ok {
		return nil, EnvNotSetError(envP12Certificate)
	}

	publicKey, ok := os.LookupEnv(envPublicKey)
	if !ok {
		return nil, EnvNotSetError(envPublicKey)
	}

	certPassword, ok := os.LookupEnv(envCertPassword)
	if !ok {
		return nil, EnvNotSetError(envCertPassword)
	}

	kafkaProducer, ok := os.LookupEnv(envKafkaProducer)
	if !ok {
		return nil, EnvNotSetError(envKafkaProducer)
	}

	kafkaConsumer, ok := os.LookupEnv(envKafkaConsumer)
	if !ok {
		return nil, EnvNotSetError(envKafkaConsumer)
	}

	redisServer, ok := os.LookupEnv(envRedisServer)
	if !ok {
		return nil, EnvNotSetError(envRedisServer)
	}

	return &Config{
		DatabaseDSN:         dbDSN,
		P12Certificate:      p12Cert,
		PublicKey:           publicKey,
		CertPassword:        certPassword,
		KafkaConsumerServer: kafkaConsumer,
		KafkaProducerServer: kafkaProducer,
		RedisServer:         redisServer,
	}, nil
}
