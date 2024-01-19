package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (m CheatcodeModel) GetAll(code string, description string, tags []string, filters Filters) ([]*Cheatcode, Metadata, error) {
	query := fmt.Sprintf(`
    SELECT count(*) OVER(), id, created_at, code, description, tags, version
    FROM cheatcodes
    WHERE (to_tsvector('simple', code) @@ plainto_tsquery('simple', $1) OR $1 = '')
    AND (to_tsvector('simple', description) @@ plainto_tsquery('simple', $2) OR $2 = '')
    AND (tags @> $3 OR $3 = '{}')
    ORDER BY %s %s, id ASC
    LIMIT $4 OFFSET $5
  `, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{code, description, pq.Array(tags), filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	cheatcodes := []*Cheatcode{}

	for rows.Next() {
		var cheatcode Cheatcode

		err := rows.Scan(
			&totalRecords,
			&cheatcode.ID,
			&cheatcode.CreatedAt,
			&cheatcode.Code,
			&cheatcode.Description,
			pq.Array(&cheatcode.Tags),
			&cheatcode.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		cheatcodes = append(cheatcodes, &cheatcode)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return cheatcodes, metadata, nil
}
