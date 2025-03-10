package product

import (
	"Ecommerce-basic/infra/response"
	"context"
	"database/sql"
	"time"

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

func (r repository) CreateProduct(ctx context.Context, model Product) (err error) {
	query := `
        INSERT INTO products (
            sku, name, stock, price, created_at, updated_at
        ) VALUES (
            :sku, :name, :stock, :price, :created_at, :updated_at
        )
    `
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, model)
	return
}

func (r repository) GetAllProductsWithPaginationCursor(ctx context.Context, model ProductPagination) (products []Product, err error) {
	query := `
        SELECT 
            id, sku, name, stock, price, created_at, updated_at, deleted_at
        FROM products
        WHERE id > $1 AND deleted_at IS NULL
        ORDER BY id ASC
        LIMIT $2
    `
	err = r.db.SelectContext(ctx, &products, query, model.Cursor, model.Size)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.ErrNotFound
		}
		return
	}
	return
}

func (r repository) GetProductBySKU(ctx context.Context, sku string) (product Product, err error) {
	query := `
        SELECT 
            id, sku, name, stock, price, created_at, updated_at, deleted_at
        FROM products
        WHERE sku = $1 AND deleted_at IS NULL
    `
	err = r.db.GetContext(ctx, &product, query, sku)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
		}
		return
	}
	return
}

func (r repository) GetProductByID(ctx context.Context, id int) (product Product, err error) {
	query := `
		SELECT 
			id, sku, name, stock, price, created_at, updated_at, deleted_at
		FROM products
		WHERE id=$1 AND deleted_at IS NULL
	`

	err = r.db.GetContext(ctx, &product, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
		}
		return
	}
	return
}

func (r repository) UpdateProduct(ctx context.Context, model Product) (err error) {
	query := `
		UPDATE products
		SET name=:name, stock=:stock, price=:price, updated_at=:updated_at
		WHERE id=:id AND deleted_at IS NULL
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, model)
	return
}

func (r repository) SoftDeleteProduct(ctx context.Context, id int) (err error) {
	query := `
		UPDATE products
		SET deleted_at=:deleted_at
		WHERE id=:id AND deleted_at IS NULL
	`

	deletedAt := time.Now()
	_, err = r.db.ExecContext(ctx, query, map[string]interface{}{
		"id":         id,
		"deleted_at": deletedAt,
	})
	return
}

func (r repository) SearchProducts(ctx context.Context, keyword string, pagination ProductPagination) (products []Product, err error) {
	query := `
		SELECT 
			id, sku, name, stock, price, created_at, updated_at, deleted_at
		FROM products
		WHERE (name ILIKE $1 OR sku ILIKE $1) AND deleted_at IS NULL
		ORDER BY id ASC
		LIMIT $2 OFFSET $3
	`

	keyword = "%" + keyword + "%"
	err = r.db.SelectContext(ctx, &products, query, keyword, pagination.Size, pagination.Cursor)
	return
}

func (r repository) FilterProducts(ctx context.Context, minPrice, maxPrice int, minStock, maxStock int16, pagination ProductPagination) (products []Product, err error) {
	query := `
        SELECT 
            id, sku, name, stock, price, created_at, updated_at, deleted_at
        FROM products
        WHERE (price BETWEEN $1 AND $2) 
          AND (stock BETWEEN $3 AND $4) 
          AND deleted_at IS NULL
        ORDER BY id ASC
        LIMIT $5 OFFSET $6
    `
	err = r.db.SelectContext(ctx, &products, query, minPrice, maxPrice, minStock, maxStock, pagination.Size, pagination.Cursor)
	return
}

// untuk validate unique
func (r repository) GetProductByName(ctx context.Context, name string) (product Product, err error) {
	query := `
       SELECT 
          id, sku, name, stock, price, created_at, updated_at, deleted_at
       FROM products
       WHERE name=$1 AND deleted_at IS NULL
    `

	err = r.db.GetContext(ctx, &product, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
		}
		return
	}
	return
}
