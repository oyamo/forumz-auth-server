package user

import (
	"github.com/google/uuid"
	"github.com/oyamo/forumz-auth-server/internal/pkg"
	"time"
)

type RegistrationRequest struct {
	FirstName    string   `json:"firstName" validate:"required"`
	LastName     string   `json:"lastName" validate:"required"`
	EmailAddress string   `json:"emailAddress" validate:"required,email"`
	Username     string   `json:"username" validate:"required"`
	Password     string   `json:"password" validate:"required,min=8,max=16"`
	Dob          pkg.Date `json:"dob" validate:"required"`
}

type UpdateInfoRequest struct {
	Id        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Dob       pkg.Date  `json:"dob" validate:"required"`
}

type UpdateInfoResponse struct {
	Id              uuid.UUID `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	EmailAddress    string    `json:"emailAddress"`
	Username        string    `json:"username"`
	Dob             time.Time `json:"dob"`
	DatetimeCreated time.Time `json:"datetimeCreated"`
	LastModified    time.Time `json:"lastModified"`
}

type UserInfo struct {
	Id              uuid.UUID `json:"id"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	EmailAddress    string    `json:"emailAddress"`
	Username        string    `json:"username"`
	Dob             time.Time `json:"dob"`
	DatetimeCreated time.Time `json:"datetimeCreated"`
	LastModified    time.Time `json:"lastModified"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type Token struct {
	AccessToken string    `json:"accessToken"`
	ExpiresIn   int64     `json:"expiresIn"`
	Sub         uuid.UUID `json:"sub"`
}
