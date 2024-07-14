package user

import (
	"github.com/google/uuid"
	"time"
)

type Person struct {
	Id              uuid.UUID
	FirstName       string
	LastName        string
	EmailAddress    string
	Username        string
	Status          string
	PasswordHash    string
	Dob             time.Time
	DatetimeCreated time.Time
	LastModified    time.Time
}
