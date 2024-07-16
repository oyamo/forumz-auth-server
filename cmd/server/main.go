package main

import (
	"context"
	"errors"
	"github.com/oyamo/forumz-auth-server/internal/domain/connections"
	"github.com/oyamo/forumz-auth-server/internal/domain/user"
	"github.com/oyamo/forumz-auth-server/internal/infrastructure/persistance/postgres"
	redis_cache "github.com/oyamo/forumz-auth-server/internal/infrastructure/persistance/redis-cache"
	"github.com/oyamo/forumz-auth-server/internal/interfaces/web"
	"github.com/oyamo/forumz-auth-server/internal/pkg"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"log"
)

func main() {
	logConf, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	logger := logConf.Sugar()

	conf, err := pkg.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	// Set up OpenTelemetry.
	otelShutdown, err := pkg.SetupOTelSDK(context.Background())
	if err != nil {
		return
	}

	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	privateKey, err := getPrivateKeyFromP12(conf.P12Certificate, conf.CertPassword)
	if err != nil {
		logger.Fatal(err)
	}

	publicKey, err := getPublicKeyFromFile(conf.PublicKey)
	if err != nil {
		logger.Fatal(err)
	}

	conn, err := pkg.NewPostgresClient(conf.DatabaseDSN)
	defer conn.Close()
	if err != nil {
		logger.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: conf.RedisServer,
	})

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		logger.Fatalf("error dialing redis: %s\n", err)
	}

	jsonSender := pkg.NewJSONSender(conf)

	personRepo := postgres.NewPersonRepository(conn)
	redisPersonRepo := redis_cache.NewRedisPersonRepository(redisClient)
	connectionRepo := postgres.NewConnectionRepository(conn)
	connectionsUC := connections.NewUseCase(connectionRepo)
	personsUC := user.NewUseCase(personRepo, redisPersonRepo, logger, conf, privateKey, publicKey)

	tracer := otel.Tracer("forumz-auth-server")
	router := web.NewRouter(logger, jsonSender, connectionsUC, personsUC, publicKey, tracer)
	ginEngine := router.Setup()

	logger.Fatal(ginEngine.Run(":3000"))
}
