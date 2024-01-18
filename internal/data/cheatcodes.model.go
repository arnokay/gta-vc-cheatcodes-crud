package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type CheatcodeModel struct {
	DB *sql.DB
}

func (m CheatcodeModel) Insert(cheatcode *Cheatcode) error {
	query := `
    INSERT INTO cheatcodes (code, description, tags)
    VALUES ($1, $2, $3)
    RETURNING id, created_at, version
  `

	args := []any{cheatcode.Code, cheatcode.Description, pq.Array(cheatcode.Tags)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&cheatcode.ID, &cheatcode.CreatedAt, &cheatcode.Version)
}

func (m CheatcodeModel) Get(id int64) (*Cheatcode, error) {
	query := `
    SELECT id, created_at, code, description, tags, version
    FROM cheatcodes
    WHERE id = $1
  `

	var cheatcode Cheatcode

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&cheatcode.ID,
		&cheatcode.CreatedAt,
		&cheatcode.Code,
		&cheatcode.Description,
		pq.Array(&cheatcode.Tags),
		&cheatcode.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &cheatcode, nil
}

func (m CheatcodeModel) Update(cheatcode *Cheatcode) error {
	query := `
    UPDATE cheatcodes
    SET code = $1, description = $2, tags = $3, version = version + 1
    WHERE id = $4 AND version = $5
    RETURNING version
  `

	args := []any{
		cheatcode.Code,
		cheatcode.Description,
		pq.Array(cheatcode.Tags),
		cheatcode.ID,
		cheatcode.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&cheatcode.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m CheatcodeModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
    DELETE FROM cheatcodes
    WHERE id = $1
  `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
