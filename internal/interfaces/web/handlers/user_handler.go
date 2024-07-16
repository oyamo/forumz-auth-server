package handlers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/oyamo/forumz-auth-server/internal/domain/user"
	"github.com/oyamo/forumz-auth-server/internal/interfaces/web/dto"
	"github.com/oyamo/forumz-auth-server/internal/pkg"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
)

type UserHandler struct {
	useCase    *user.UseCase
	logger     *zap.SugaredLogger
	validator  *validator.Validate
	jsonSender *pkg.JSONSender
	tracer     trace.Tracer
}

func (h *UserHandler) Register(c *gin.Context) {
	ctx, span := h.tracer.Start(c, "UserHandler.Register")
	defer span.End()

	var responseDto dto.ResponseDto
	requestIdCtx, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		h.logger.Errorw("cannot find id from context")
		return
	}

	requestId, ok := requestIdCtx.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{})
		h.logger.Errorw("requestId is not uuid type")
		return
	}

	responseDto.RequestId = requestId

	var personDTO user.RegistrationRequest
	if err := c.ShouldBindJSON(&personDTO); err != nil {
		responseDto.Description = "Invalid request body"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	if err := h.validator.Struct(&personDTO); err != nil {
		responseDto.Description = "request has missing/invalid fields"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	ctx = context.WithValue(ctx, "id", responseDto.RequestId)

	ret, err := h.useCase.Register(ctx, &personDTO)
	if err == nil {
		responseDto.Description = "Success"
		responseDto.Data = ret
		c.JSON(http.StatusOK, responseDto)
		return
	}

	switch {
	case errors.Is(err, user.ErrUserAlreadyExists):
		responseDto.Description = err.Error()
		c.JSON(http.StatusBadRequest, responseDto)
		return
	default:
		c.JSON(http.StatusInternalServerError, responseDto)
	}

	err = h.jsonSender.Send("Put-Person-v1", ret)
	if err != nil {
		h.logger.Error(err)
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	ctx, span := h.tracer.Start(c, "UserHandler.Login")
	defer span.End()

	var responseDto dto.ResponseDto
	requestIdCtx, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		h.logger.Errorw("cannot find id from context")
		return
	}

	requestId, ok := requestIdCtx.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{})
		h.logger.Errorw("requestId is not uuid type")
		return
	}

	responseDto.RequestId = requestId

	var personDTO user.LoginRequest
	if err := c.ShouldBindJSON(&personDTO); err != nil {
		responseDto.Description = "Invalid request body"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	if err := h.validator.Struct(&personDTO); err != nil {
		responseDto.Description = "invalid request body"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	ctx = context.WithValue(ctx, "id", responseDto.RequestId)

	ret, err := h.useCase.Login(ctx, &personDTO)
	if err == nil {
		responseDto.Description = "Success"
		responseDto.Data = ret
		c.JSON(http.StatusOK, responseDto)
		return
	}

	switch {
	case errors.Is(err, user.ErrUserNotFound):
		responseDto.Description = err.Error()
		c.JSON(http.StatusBadRequest, responseDto)
	case errors.Is(err, user.ErrIncorrectCredentials):
		responseDto.Description = err.Error()
		c.JSON(http.StatusUnauthorized, responseDto)
	}
}

func (h *UserHandler) Update(c *gin.Context) {
	var responseDto dto.ResponseDto
	requestIdCtx, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		h.logger.Errorw("cannot find id from context")
		return
	}

	requestId, ok := requestIdCtx.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{})
		h.logger.Errorw("requestId is not uuid type")
		return
	}

	responseDto.RequestId = requestId

	initiator, exists := c.Get("initiator")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{})
		h.logger.Errorw("cannot find initiator from context")
		return
	}

	personId, err := uuid.Parse(c.Param("personId"))
	if err != nil {
		responseDto.Description = "invalid request body"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	initiatorId := initiator.(uuid.UUID)
	if initiatorId != personId {
		responseDto.Description = "Unauthorized"
		c.JSON(http.StatusUnauthorized, responseDto)
		return
	}

	var request user.UpdateInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		responseDto.Description = "Invalid request body"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	if err := h.validator.Struct(&request); err != nil {
		responseDto.Description = "invalid request body"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	request.Id = initiatorId
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "id", responseDto.RequestId)

	ret, err := h.useCase.Update(ctx, &request)
	if err != nil {
		responseDto.Description = "Something went wrong"
		c.JSON(http.StatusInternalServerError, responseDto)
		return
	}

	responseDto.Description = "Success"
	responseDto.Data = ret
	c.JSON(http.StatusOK, responseDto)

	err = h.jsonSender.Send("Put-Person-v1", ret)
	if err != nil {
		h.logger.Error(err)
		return
	}
}

func (h *UserHandler) UserInfo(c *gin.Context) {
	ctx, span := h.tracer.Start(c, "UserHandler.userInfo")
	defer span.End()

	var responseDto dto.ResponseDto
	requestIdCtx, exists := c.Get("id")
	if !exists {
		err := errors.New("cannot find id from context")
		c.JSON(http.StatusInternalServerError, gin.H{})
		h.logger.Errorw("error while getting user info from context", "error", err)
		return
	}

	requestId, ok := requestIdCtx.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{})
		err := errors.New("requestId is not uuid type")
		h.logger.Errorw("requestId is not uuid type", "error", err)
		return
	}

	responseDto.RequestId = requestId

	personId, err := uuid.Parse(c.Param("personId"))
	if err != nil {
		responseDto.Description = "invalid personId"
		c.JSON(http.StatusBadRequest, responseDto)
		return
	}

	ctx = context.WithValue(ctx, "id", responseDto.RequestId)

	info, err := h.useCase.UserInfo(personId, ctx)
	if err == nil {
		responseDto.Description = "Success"
		responseDto.Data = info
		c.JSON(http.StatusOK, responseDto)
		return
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		responseDto.Description = err.Error()
		c.JSON(http.StatusBadRequest, responseDto)
	default:
		responseDto.Description = "Something went wrong"
		c.JSON(http.StatusInternalServerError, responseDto)
	}

}

func NewUserHandler(useCase *user.UseCase, logger *zap.SugaredLogger, sender *pkg.JSONSender, tracer trace.Tracer) *UserHandler {
	return &UserHandler{
		useCase:    useCase,
		logger:     logger,
		validator:  validator.New(),
		jsonSender: sender,
		tracer:     tracer,
	}
}
