package user

import (
	"context"
	"github.com/google/uuid"
)

type PersonRepository interface {
	Upsert(context.Context, *Person) error
	UpdatePassword(context.Context, uuid.UUID, string) error
	Find(context.Context, uuid.UUID) (*Person, error)
	FindByUsername(context.Context, string) (*Person, error)
	Exists(context.Context, uuid.UUID) (bool, error)
	ExistsByEmail(context.Context, string) (bool, error)
	ExistsByUsername(context.Context, string) (bool, error)
}
