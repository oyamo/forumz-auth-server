package postgres

import (
	"auth/internal/domain/user"
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type psqlPersonRespository struct {
	db *sql.DB
}

func (repo psqlPersonRespository) Upsert(ctx context.Context, person *user.Person) error {
	stmt, err := repo.db.Prepare(`insert into person (id,first_name, last_name, email_address, password_hash, username, dob) 
    values ($1, $2, $3, $4, $5, $6, $7) 
    on conflict(id) do update set first_name = $2, last_name = $3, dob = $7`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		person.Id,
		person.FirstName,
		person.LastName,
		person.EmailAddress,
		person.PasswordHash,
		person.Username,
		person.Dob)

	if err != nil {
		return err
	}

	return nil
}

func (repo psqlPersonRespository) UpdatePassword(ctx context.Context, u uuid.UUID, s string) error {
	stmt, err := repo.db.Prepare(`update person set password_hash = $1 where id = $2`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, s, u)
	if err != nil {
		return err
	}

	return nil
}

func (repo psqlPersonRespository) Find(ctx context.Context, u uuid.UUID) (*user.Person, error) {
	stmt, err := repo.db.Prepare(`select id, first_name, last_name, email_address, password_hash, username, status, dob, datetime_created, last_modified from person where id = $1`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	var person user.Person

	row := stmt.QueryRowContext(ctx, u)
	err = row.Scan(
		&person.Id,
		&person.FirstName,
		&person.LastName,
		&person.EmailAddress,
		&person.PasswordHash,
		&person.Username,
		&person.Status,
		&person.Dob,
		&person.DatetimeCreated,
		&person.LastModified,
	)

	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (repo psqlPersonRespository) FindByUsername(ctx context.Context, username string) (*user.Person, error) {
	stmt, err := repo.db.Prepare(`select id, first_name, last_name, email_address, password_hash, username, status, dob, datetime_created, last_modified from person where username = $1 or email_address = $1 limit 1`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	var person user.Person
	row := stmt.QueryRowContext(ctx, username)
	err = row.Scan(
		&person.Id,
		&person.FirstName,
		&person.LastName,
		&person.EmailAddress,
		&person.PasswordHash,
		&person.Username,
		&person.Status,
		&person.Dob,
		&person.DatetimeCreated,
		&person.LastModified,
	)

	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (repo psqlPersonRespository) Exists(ctx context.Context, u uuid.UUID) (bool, error) {
	stmt, err := repo.db.Prepare(`select exists(select 1 from person where id = $1)`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var exists bool
	err = stmt.QueryRowContext(ctx, u).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (repo psqlPersonRespository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	stmt, err := repo.db.Prepare(`select exists(select 1 from person where email_address = $1)`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var exists bool
	err = stmt.QueryRowContext(ctx, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (repo psqlPersonRespository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	stmt, err := repo.db.Prepare(`select exists(select 1 from person where username = $1)`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRowContext(ctx, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func NewPersonRepository(db *sql.DB) user.PersonRepository {
	return &psqlPersonRespository{
		db: db,
	}
}
