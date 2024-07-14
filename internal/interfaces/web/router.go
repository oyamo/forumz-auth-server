package web

import (
	"auth/internal/domain/connections"
	"auth/internal/domain/user"
	"auth/internal/interfaces/web/handlers"
	"auth/internal/pkg"
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router struct {
	logger       *zap.SugaredLogger
	connectionUC *connections.UseCase
	userUC       *user.UseCase
	jsonSender   *pkg.JSONSender
	pub          *rsa.PublicKey
}

func (router *Router) Setup() *gin.Engine {

	userHandler := handlers.NewUserHandler(router.userUC, router.logger, router.jsonSender)
	connectionHandler := handlers.NewConnectionHandler(router.connectionUC, router.userUC, router.logger, router.jsonSender)
	middlewareHandler := handlers.NewMiddlewareHandler(router.pub, router.logger)

	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.Use(middlewareHandler.AddRequestID)

	auth := v1.Group("/auth")
	conn := v1.Group("/connections")

	auth.POST("/login", userHandler.Login)
	auth.POST("/register", userHandler.Register)
	auth.GET("/persons/:personId", middlewareHandler.AuthenticateRequest, userHandler.UserInfo)
	auth.PATCH("/persons/:personId", middlewareHandler.AuthenticateRequest, userHandler.Update)

	conn.Use(middlewareHandler.AuthenticateRequest)
	conn.POST("/", connectionHandler.Connect)
	conn.GET("/", connectionHandler.ListConnections)
	conn.DELETE("/", connectionHandler.Disconnect)

	return r
}

func NewRouter(logger *zap.SugaredLogger, sender *pkg.JSONSender, connectionUC *connections.UseCase, userUC *user.UseCase, pub *rsa.PublicKey) *Router {
	return &Router{
		logger:       logger,
		connectionUC: connectionUC,
		userUC:       userUC,
		jsonSender:   sender,
		pub:          pub,
	}
}
