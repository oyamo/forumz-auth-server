package connections

import (
	"github.com/google/uuid"
	"time"
)

type Connection struct {
	UserId          uuid.UUID
	ConnectedTo     uuid.UUID
	DatetimeCreated time.Time
}
