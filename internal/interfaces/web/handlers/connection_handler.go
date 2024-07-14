package handlers

import (
	"auth/internal/domain/connections"
	"auth/internal/domain/user"
	"auth/internal/interfaces/web/dto"
	"auth/internal/pkg"
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ConnectionHandler struct {
	useCase       *connections.UseCase
	personUseCase *user.UseCase
	logger        *zap.SugaredLogger
	jsonSender    *pkg.JSONSender
	v             *validator.Validate
}

func (handler *ConnectionHandler) Connect(c *gin.Context) {
	var responseDto dto.ResponseDto
	requestIdCtx, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("cannot find id from context")
		return
	}

	requestId, ok := requestIdCtx.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("requestId is not uuid type")
		return
	}

	responseDto.RequestId = requestId
	var req connections.CreateConnectionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		responseDto.Description = "invalid request body"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	if err := handler.v.Struct(req); err != nil {
		responseDto.Description = "invalid request body"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	initiator, exists := c.Get("initiator")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("cannot find id from context")
		return
	}

	userId, isUUID := initiator.(uuid.UUID)
	if !isUUID {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("initiator is not uuid")
		return
	}

	if userId == req.ConnectionTo {
		responseDto.Description = "Operation not allowed"
		c.JSON(http.StatusUnprocessableEntity, responseDto)
		return
	}

	req.UserId = userId
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "id", responseDto.RequestId)

	connectionTo, err := handler.personUseCase.UserInfo(req.ConnectionTo, ctx)
	if err != nil {
		responseDto.Description = "something went wrong"
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, sql.ErrNoRows):
			responseDto.Description = "Other person not found"
			status = http.StatusNotFound
		}
		handler.logger.Errorw("cannot get user info", "error", err)
		c.JSON(status, responseDto)
		return
	}

	err = handler.useCase.Connect(ctx, &req)
	if err != nil {
		responseDto.Description = "Something went wrong"
		c.JSON(http.StatusInternalServerError, responseDto)
		handler.logger.Errorw("error on connect", "error", err)
		return
	}

	responseDto.Description = "Connection successfully created."
	c.JSON(http.StatusOK, responseDto)

	err = handler.jsonSender.Send("Put-Connection-v1", req)
	if err != nil {
		handler.logger.Error(err)
		return
	}

	notification := map[string]interface{}{
		"datetimeCreated": time.Now(),
		"recipient":       userId,
		"type":            "Connection",
		"additionalInfo": map[string]interface{}{
			"connectionTo":       req.ConnectionTo,
			"connectionFrom":     req.UserId,
			"connectionFromName": connectionTo.FirstName,
		},
	}

	err = handler.jsonSender.Send("Put-Notification-v1", notification)
	if err != nil {
		handler.logger.Error(err)
		return
	}

}

func (handler *ConnectionHandler) Disconnect(c *gin.Context) {
	var responseDto dto.ResponseDto
	requestIdCtx, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("cannot find id from context")
		return
	}

	requestId, ok := requestIdCtx.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("requestId is not uuid type")
		return
	}

	responseDto.RequestId = requestId

	var req connections.CreateConnectionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		responseDto.Description = "invalid request body"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	initiator, exists := c.Get("initiator")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("cannot find id from context")
		return
	}

	userId, isUUID := initiator.(uuid.UUID)
	if !isUUID {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("initiator is not uuid")
		return
	}

	req.UserId = userId
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "id", responseDto.RequestId)

	err := handler.useCase.Disconnect(ctx, &req)
	if err != nil {
		responseDto.Description = "Something went wrong"
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	responseDto.Description = "Connection successfully deleted."
	c.JSON(http.StatusOK, responseDto)

	err = handler.jsonSender.Send("Delete-Connection-v1", req)
	if err != nil {
		handler.logger.Error(err)
		return
	}
}

func (handler *ConnectionHandler) ListConnections(c *gin.Context) {
	var responseDto dto.ResponseDto
	requestIdCtx, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("cannot find id from context")
		return
	}

	requestId, ok := requestIdCtx.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("requestId is not uuid type")
		return
	}

	responseDto.RequestId = requestId
	initiator, exists := c.Get("initiator")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("cannot find id from context")
		return
	}

	userId, isUUID := initiator.(uuid.UUID)
	if !isUUID {
		c.JSON(http.StatusInternalServerError, gin.H{})
		handler.logger.Errorw("initiator is not uuid")
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "id", responseDto.RequestId)
	conns, err := handler.useCase.ListConnections(ctx, userId)
	if err != nil {
		responseDto.Description = "Something went wrong"
		c.JSON(http.StatusInternalServerError, responseDto)
		handler.logger.Error(err)
		return
	}

	responseDto.Data = conns
	responseDto.Description = "Success"
	c.JSON(http.StatusOK, responseDto)
}

func NewConnectionHandler(useCase *connections.UseCase, personUseCase *user.UseCase, logger *zap.SugaredLogger, sender *pkg.JSONSender) *ConnectionHandler {
	return &ConnectionHandler{
		useCase:       useCase,
		logger:        logger,
		v:             validator.New(),
		jsonSender:    sender,
		personUseCase: personUseCase,
	}
}
