package connections

import (
	"context"
	"github.com/google/uuid"
)

type UseCase struct {
	connectionRepository Repository
}

func (uc *UseCase) Connect(ctx context.Context, dto *CreateConnectionDTO) error {
	connection := &Connection{
		UserId:      dto.UserId,
		ConnectedTo: dto.ConnectionTo,
	}
	err := uc.connectionRepository.Save(ctx, connection)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) Disconnect(ctx context.Context, dto *CreateConnectionDTO) error {
	err := uc.connectionRepository.Delete(ctx, dto.UserId, dto.ConnectionTo)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) ListConnections(ctx context.Context, id uuid.UUID) ([]ConnectionsItem, error) {
	connections, err := uc.connectionRepository.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	return connections, nil
}

func NewUseCase(connectionRepository Repository) *UseCase {
	return &UseCase{connectionRepository: connectionRepository}
}
