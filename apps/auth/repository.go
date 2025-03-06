package auth

import (
	"Ecommerce-basic/infra/response"
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func newRepository(db *sqlx.DB) repository {
	return repository{
		db: db,
	}
}

func (r repository) CreateAuth(ctx context.Context, model AuthEntity) (err error) {
	query := `
		INSERT INTO auth (
			email, password, role, created_at, updated_at, public_id
		) VALUES (
			:email, :password, :role, :created_at, :updated_at, :public_id
		)
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	defer func() {
		if cerr := stmt.Close(); cerr != nil {
			log.Println("failed to close statement:", cerr)
		}
	}()

	_, err = stmt.ExecContext(ctx, model)

	return
}

// GetAuthByEmail implements Repository.
func (r repository) GetAuthByEmail(ctx context.Context, email string) (model AuthEntity, err error) {
	query := `
		SELECT 
			id, email, password, role, created_at, updated_at, public_id
		FROM auth
		WHERE email=$1
	`

	err = r.db.GetContext(ctx, &model, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return
		}
		return
	}

	return
}
