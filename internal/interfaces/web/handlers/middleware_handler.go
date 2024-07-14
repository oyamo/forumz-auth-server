package handlers

import (
	"auth/internal/interfaces/web/dto"
	"crypto/rsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type MiddlewareHandler struct {
	publicKey *rsa.PublicKey
	logger    *zap.SugaredLogger
}

func (mw *MiddlewareHandler) AddRequestID(c *gin.Context) {
	requestId, err := uuid.NewV7()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		mw.logger.Warnw("failed to generate request id", zap.Error(err))
		return
	}
	c.Set("id", requestId)
	c.Next()
}

func (mw *MiddlewareHandler) AuthenticateRequest(c *gin.Context) {
	var res dto.ResponseDto
	requestIdCtx, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		c.Abort()
		mw.logger.Errorw("cannot find id from context")
		return
	}

	requestId, ok := requestIdCtx.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{})
		c.Abort()
		mw.logger.Errorw("requestId is not uuid type")
		return
	}

	res.RequestId = requestId
	headers := c.Request.Header
	token := strings.TrimPrefix(headers.Get("Authorization"), "Bearer ")
	if token == "" {
		res.Description = "Missing token"
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		mw.logger.Errorw("token is empty")
		return
	}

	claims := &jwt.MapClaims{}
	_, err := jwt.NewParser().ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mw.publicKey, nil
	})

	if err != nil {
		res.Description = "Invalid/expired token"
		c.JSON(http.StatusUnauthorized, res)
		mw.logger.Errorw("invalid token", zap.Error(err), zap.String("token", token))
		c.Abort()
		return
	}

	sub, err := claims.GetSubject()
	if err != nil {
		res.Description = "Invalid token"
		c.JSON(http.StatusUnauthorized, res)
		return
	}
	identity, err := uuid.Parse(sub)
	if err != nil {
		res.Description = "Invalid token"
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	c.Set("initiator", identity)
	c.Next()
}

func NewMiddlewareHandler(publicKey *rsa.PublicKey, logger *zap.SugaredLogger) *MiddlewareHandler {
	return &MiddlewareHandler{
		publicKey: publicKey,
		logger:    logger,
	}
}
