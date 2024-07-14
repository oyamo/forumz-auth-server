package connections

import (
	"github.com/google/uuid"
	"time"
)

type CreateConnectionDTO struct {
	UserId       uuid.UUID `json:"userId"`
	ConnectionTo uuid.UUID `json:"connectionTo" validate:"required"`
}

type ConnectionsItem struct {
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	Id              string    `json:"id"`
	DatetimeCreated time.Time `json:"datetimeCreated"`
}
