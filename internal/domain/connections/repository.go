package connections

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, connection *Connection) error
	Find(ctx context.Context, id uuid.UUID) ([]ConnectionsItem, error)
	Delete(ctx context.Context, id, connectedTo uuid.UUID) error
}
