package user

import (
	"auth/internal/pkg"
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type UseCase struct {
	personRepository PersonRepository
	logger           *zap.SugaredLogger
	conf             *pkg.Config
	privateKey       *rsa.PrivateKey
	publicKey        *rsa.PublicKey
}

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrIncorrectCredentials = errors.New("incorrect credentials")
)

func (u *UseCase) Register(ctx context.Context, dto *RegistrationRequest) (*UpdateInfoResponse, error) {
	exists, err := u.personRepository.ExistsByUsername(ctx, dto.Username)
	if err != nil {
		u.logger.Errorw("error while checking username exists", "error", err)
		return nil, err
	}

	if exists {
		err = errors.New("username already exists")
		err = errors.Join(ErrUserAlreadyExists, err)
		return nil, err
	}

	exists, err = u.personRepository.ExistsByEmail(ctx, dto.EmailAddress)
	if err != nil {
		u.logger.Errorw("error while checking if email exists", "error", err)
		return nil, err
	}

	if exists {
		err = errors.New("email already exists")
		err = errors.Join(ErrUserAlreadyExists, err)
		return nil, err
	}

	// hash password
	hashedPassword, err := pkg.HashPassword(dto.Password)
	if err != nil {
		u.logger.Errorw("error while hashing password", "error", err)
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		u.logger.Errorw("error while generating uuid", "error", err)
		return nil, err
	}

	person := Person{
		Id:           id,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		EmailAddress: dto.EmailAddress,
		Username:     dto.Username,
		PasswordHash: hashedPassword,
		Dob:          time.Time(dto.Dob),
	}

	err = u.personRepository.Upsert(ctx, &person)
	if err != nil {
		u.logger.Errorw("error while upserting user", "error", err)
		return nil, err
	}

	ret := &UpdateInfoResponse{
		Id:              person.Id,
		FirstName:       person.FirstName,
		LastName:        person.LastName,
		EmailAddress:    person.EmailAddress,
		Username:        person.Username,
		Dob:             time.Time(dto.Dob),
		DatetimeCreated: time.Now(),
		LastModified:    time.Now(),
	}

	return ret, nil
}

func (u *UseCase) Login(ctx context.Context, dto *LoginRequest) (*Token, error) {
	user, err := u.personRepository.FindByUsername(ctx, dto.Username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUserNotFound
		}
		u.logger.Errorw("error while checking if user exists", "error", err)
		return nil, err
	}

	correct, err := pkg.ComparePasswordAndHash(dto.Password, user.PasswordHash)
	if err != nil {
		u.logger.Errorw("error while comparing password", "error", err)
		return nil, err
	}

	if !correct {
		return nil, ErrIncorrectCredentials
	}

	expiry := time.Now().Add(time.Hour * 24)
	iat := time.Now().Unix()
	nbf := time.Now().Unix()

	claims := jwt.MapClaims{
		"sub": user.Id,
		"exp": expiry.Unix(),
		"iat": iat,
		"nbf": nbf,
		"iss": "http://localhost",
		"jti": uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	encoded, err := token.SignedString(u.privateKey)
	if err != nil {
		u.logger.Errorw("error while signing token", "error", err)
		return nil, fmt.Errorf("create: sign token: %w", err)
	}

	expiresIn := expiry.Sub(time.Now())
	ret := Token{
		AccessToken: encoded,
		ExpiresIn:   int64(expiresIn.Seconds()),
		Sub:         user.Id,
	}

	return &ret, nil
}

func (u *UseCase) Update(ctx context.Context, dto *UpdateInfoRequest) (*UpdateInfoResponse, error) {
	exists, err := u.personRepository.Exists(ctx, dto.Id)
	if err != nil {
		u.logger.Errorw("error while checking if user exists", "error", err)
		return nil, err
	}

	if !exists {
		return nil, ErrUserNotFound
	}

	info := Person{
		Id:        dto.Id,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Dob:       time.Time(dto.Dob),
	}

	err = u.personRepository.Upsert(ctx, &info)
	if err != nil {
		u.logger.Errorw("error while updating user", "error", err)
		return nil, err
	}

	// return user
	person, err := u.personRepository.Find(ctx, dto.Id)
	if err != nil {
		u.logger.Errorw("error while checking if user exists", "error", err)
		return nil, err
	}

	updateInfo := &UpdateInfoResponse{
		Id:              person.Id,
		FirstName:       person.FirstName,
		LastName:        person.LastName,
		EmailAddress:    person.EmailAddress,
		Username:        person.Username,
		Dob:             time.Time(dto.Dob),
		DatetimeCreated: person.DatetimeCreated,
		LastModified:    person.LastModified,
	}

	return updateInfo, nil
}

func (u *UseCase) UserInfo(id uuid.UUID, ctx context.Context) (*UserInfo, error) {
	info, err := u.personRepository.Find(ctx, id)
	if err != nil {
		u.logger.Errorw("error while finding user", "error", err)
		return nil, err
	}

	ret := UserInfo{
		Id:              info.Id,
		FirstName:       info.FirstName,
		LastName:        info.LastName,
		EmailAddress:    info.EmailAddress,
		Username:        info.Username,
		Dob:             info.Dob,
		DatetimeCreated: info.DatetimeCreated,
		LastModified:    info.LastModified,
	}

	return &ret, nil
}

func NewUseCase(personRepository PersonRepository, logger *zap.SugaredLogger, conf *pkg.Config, privatekey *rsa.PrivateKey, publickey *rsa.PublicKey) *UseCase {
	return &UseCase{
		personRepository: personRepository,
		logger:           logger,
		conf:             conf,
		privateKey:       privatekey,
		publicKey:        publickey,
	}
}
