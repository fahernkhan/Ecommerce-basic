package transaction

import (
	"Ecommerce-basic/infra/response"
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	TransactionDBRepository
	TransactionRepository
	ProductRepository
}

type TransactionDBRepository interface {
	Begin(ctx context.Context) (tx *sqlx.Tx, err error)
	Rollback(ctx context.Context, tx *sqlx.Tx) (err error)
	Commit(ctx context.Context, tx *sqlx.Tx) (err error)
}

type TransactionRepository interface {
	CreateTransactionWithTx(ctx context.Context, tx *sqlx.Tx, trx Transaction) (err error)
	GetTransactionsByUserPublicId(ctx context.Context, userPublicId string) (trxs []Transaction, err error)
	GetTransactionById(ctx context.Context, trxId int) (trx Transaction, err error)                     // Method baru
	UpdateTransactionStatusWithTx(ctx context.Context, tx *sqlx.Tx, trx Transaction) (err error)        // Method baru
	GetTransactionsByProductSku(ctx context.Context, productSKU string) (trxs []Transaction, err error) // Method baru
}
type ProductRepository interface {
	GetProductBySku(ctx context.Context, productSKU string) (product Product, err error)
	UpdateProductStockWithTx(ctx context.Context, tx *sqlx.Tx, product Product) (err error)
}

type service struct {
	repo Repository
}

func newService(repo Repository) service {
	return service{
		repo: repo,
	}
}

func (s service) CreateTransaction(ctx context.Context, req CreateTransactionRequestPayload) (err error) {
	myProduct, err := s.repo.GetProductBySku(ctx, req.ProductSKU)
	if err != nil {
		return
	}

	if !myProduct.IsExists() {
		err = response.ErrNotFound
		return
	}

	trx := NewTransactionFromCreateRequest(req)
	trx.FromProduct(myProduct).
		SetPlatformFee(1_000).
		SetGrandTotal()

	if err = trx.Validate(); err != nil {
		return
	}

	if err = trx.ValidateStock(uint8(myProduct.Stock)); err != nil {
		return
	}

	// start transaction database
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return
	}

	// defer rollback if any error or after commit
	defer s.repo.Rollback(ctx, tx)

	if err = s.repo.CreateTransactionWithTx(ctx, tx, trx); err != nil {
		return
	}

	// update current stock
	if err = myProduct.UpdateStockProduct(trx.Amount); err != nil {
		return
	}

	// update into database
	if err = s.repo.UpdateProductStockWithTx(ctx, tx, myProduct); err != nil {
		return
	}

	// commit to end the transactions
	if err = s.repo.Commit(ctx, tx); err != nil {
		return
	}
	return

}

func (s service) TransactionHistories(ctx context.Context, userPublicId string) (trxs []Transaction, err error) {
	trxs, err = s.repo.GetTransactionsByUserPublicId(ctx, userPublicId)
	if err != nil {
		if err == response.ErrNotFound {
			trxs = []Transaction{}
			return trxs, nil
		}

		return
	}

	if len(trxs) == 0 {
		trxs = []Transaction{}
		return trxs, nil
	}
	return
}

// method untuk mengupdate status transaksi:
func (s service) UpdateTransactionStatus(ctx context.Context, trxId int, newStatus TransactionStatus) (err error) {
	// Mulai transaksi database
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return
	}

	defer s.repo.Rollback(ctx, tx)

	// Dapatkan transaksi berdasarkan ID
	trx, err := s.repo.GetTransactionById(ctx, trxId)
	if err != nil {
		return
	}

	// Update status transaksi
	trx.UpdateStatus(newStatus)

	// Simpan perubahan ke database
	if err = s.repo.UpdateTransactionStatusWithTx(ctx, tx, trx); err != nil {
		return
	}

	// Commit transaksi
	if err = s.repo.Commit(ctx, tx); err != nil {
		return
	}

	return
}

// method untuk mendapatkan riwayat transaksi
func (s service) GetTransactionHistoriesByProduct(ctx context.Context, productSKU string) (trxs []Transaction, err error) {
	trxs, err = s.repo.GetTransactionsByProductSku(ctx, productSKU)
	if err != nil {
		if err == response.ErrNotFound {
			trxs = []Transaction{}
			return trxs, nil
		}
		return
	}

	if len(trxs) == 0 {
		trxs = []Transaction{}
		return trxs, nil
	}
	return
}
