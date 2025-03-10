package transaction

import (
	"Ecommerce-basic/external/database"
	"Ecommerce-basic/internal"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var svc service

func init() {
	filename := "../../cmd/api/config.yaml"
	err := config.LoadConfig(filename)
	if err != nil {
		panic(err)
	}

	db, err := database.ConnectPostgres(config.Cfg.DB)
	if err != nil {
		panic(err)
	}
	repo := newRepository(db)
	svc = newService(repo)
}

func TestCreateTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := CreateTransactionRequestPayload{
			ProductSKU:   "a98dcf06-7b4b-4f33-a6d2-20738bb8081b",
			Amount:       2,
			UserPublicId: "5c534133-f81f-4df4-977e-38669242eb48",
		}

		err := svc.CreateTransaction(context.Background(), req)
		require.Nil(t, err)
	})
}

func TestUpdateTransactionStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		trxId := 1
		newStatus := TransactionStatus_Progress

		err := svc.UpdateTransactionStatus(context.Background(), trxId, newStatus)
		require.Nil(t, err)
	})
}

func TestGetTransactionHistoriesByProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		productSKU := "a98dcf06-7b4b-4f33-a6d2-20738bb8081b"

		trxs, err := svc.GetTransactionHistoriesByProduct(context.Background(), productSKU)
		require.Nil(t, err)
		require.NotEmpty(t, trxs)
	})
}
