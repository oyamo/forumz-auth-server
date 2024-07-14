package postgres

import (
	"auth/internal/domain/connections"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type psqlConnectionRepository struct {
	db *sql.DB
}

func (p psqlConnectionRepository) Save(ctx context.Context, connection *connections.Connection) error {
	stmt, err := p.db.Prepare(`insert into connection(user_id, connected_to)
		values ($1, $2) on conflict (user_id, connected_to) do nothing`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.ExecContext(
		ctx,
		connection.UserId,
		connection.ConnectedTo,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p psqlConnectionRepository) Find(ctx context.Context, id uuid.UUID) ([]connections.ConnectionsItem, error) {
	stmt, err := p.db.Prepare(`select connected_to, person.first_name, last_name,  connection.datetime_created FROM connection
                inner join person on person.id = connection.connected_to WHERE user_id = $1`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var conns []connections.ConnectionsItem
	for rows.Next() {
		var connection connections.ConnectionsItem
		err = rows.Scan(
			&connection.Id,
			&connection.FirstName,
			&connection.LastName,
			&connection.DatetimeCreated)
		if err != nil {
			return nil, err
		}
		conns = append(conns, connection)
	}

	return conns, nil
}

func (p psqlConnectionRepository) Delete(ctx context.Context, id, connectedTo uuid.UUID) error {
	stmt, err := p.db.Prepare(`delete from connection where user_id = $1 and connected_to = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id, connectedTo)
	if err != nil {
		return err
	}

	return nil
}

func NewConnectionRepository(db *sql.DB) connections.Repository {
	return &psqlConnectionRepository{
		db: db,
	}
}
