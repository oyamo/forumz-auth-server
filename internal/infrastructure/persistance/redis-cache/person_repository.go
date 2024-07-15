package redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/oyamo/forumz-auth-server/internal/domain/user"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisPersonRepository struct {
	client *redis.Client
}

var (
	personTTL = time.Minute * 5
)

func (r redisPersonRepository) Upsert(ctx context.Context, person *user.Person) error {
	person.PasswordHash = "-"

	b, err := json.Marshal(person)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("person-%s", person.Id)
	_, err = r.client.Set(ctx, key, b, personTTL).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r redisPersonRepository) UpdatePassword(ctx context.Context, uuid uuid.UUID, s string) error {
	return errors.New("not implemented")
}

func (r redisPersonRepository) Find(ctx context.Context, uuid uuid.UUID) (*user.Person, error) {
	key := fmt.Sprintf("person-%s", uuid)
	res := r.client.Get(ctx, key)
	if res.Err() != nil {
		return nil, res.Err()
	}
	var person user.Person
	err := json.Unmarshal([]byte(res.Val()), &person)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (r redisPersonRepository) FindByUsername(ctx context.Context, s string) (*user.Person, error) {
	return nil, errors.New("not implemented")
}

func (r redisPersonRepository) Exists(ctx context.Context, uuid uuid.UUID) (bool, error) {
	key := fmt.Sprintf("person-%s", uuid)
	res := r.client.Get(ctx, key)
	if res.Err() != nil {
		return false, res.Err()
	}
	var person user.Person
	err := json.Unmarshal([]byte(res.Val()), &person)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r redisPersonRepository) ExistsByEmail(ctx context.Context, s string) (bool, error) {
	return false, errors.New("not implemented")
}

func (r redisPersonRepository) ExistsByUsername(ctx context.Context, s string) (bool, error) {
	return false, errors.New("not implemented")
}

func NewRedisPersonRepository(client *redis.Client) user.PersonRepository {
	return &redisPersonRepository{
		client: client,
	}
}
