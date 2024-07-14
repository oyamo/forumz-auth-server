package main

import (
	"auth/internal/domain/connections"
	"auth/internal/domain/user"
	"auth/internal/infrastructure/persistance/postgres"
	"auth/internal/interfaces/web"
	"auth/internal/pkg"
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

	jsonSender := pkg.NewJSONSender(conf)

	personRepo := postgres.NewPersonRepository(conn)
	connectionRepo := postgres.NewConnectionRepository(conn)
	connectionsUC := connections.NewUseCase(connectionRepo)
	personsUC := user.NewUseCase(personRepo, logger, conf, privateKey, publicKey)

	router := web.NewRouter(logger, jsonSender, connectionsUC, personsUC, publicKey)
	ginEngine := router.Setup()

	logger.Fatal(ginEngine.Run(":3000"))
}
