package transaction

import (
	"fmt"
	"net/http"

	"Ecommerce-basic/infra/gin"
	"Ecommerce-basic/infra/response"
	"github.com/gin-gonic/gin"
)

type handler struct {
	svc service
}

func newHandler(svc service) handler {
	return handler{
		svc: svc,
	}
}

func (h handler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequestPayload

	// Bind JSON request body ke struct
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusBadRequest),
			infragin.WithMessage(err.Error()),
			infragin.WithError(response.ErrorBadRequest),
		)
		resp.Send(c)
		return
	}

	// Ambil userPublicId dari context (setelah middleware CheckAuth)
	userPublicId, exists := c.Get("PUBLIC_ID")
	if !exists {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusUnauthorized),
			infragin.WithMessage("user not authenticated"),
			infragin.WithError(response.ErrorUnauthorized),
		)
		resp.Send(c)
		return
	}

	// Set userPublicId ke request payload
	req.UserPublicId = fmt.Sprintf("%v", userPublicId)

	// Panggil service untuk membuat transaksi
	if err := h.svc.CreateTransaction(c.Request.Context(), req); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}

		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(c)
		return
	}

	// Kirim response sukses
	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusCreated),
		infragin.WithMessage("create transactions success"),
	)
	resp.Send(c)
}

func (h handler) GetTransactionByUser(c *gin.Context) {
	// Ambil userPublicId dari context (setelah middleware CheckAuth)
	userPublicId, exists := c.Get("PUBLIC_ID")
	if !exists {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusUnauthorized),
			infragin.WithMessage("user not authenticated"),
			infragin.WithError(response.ErrorUnauthorized),
		)
		resp.Send(c)
		return
	}

	// Panggil service untuk mendapatkan riwayat transaksi
	trxs, err := h.svc.TransactionHistories(c.Request.Context(), fmt.Sprintf("%v", userPublicId))
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}

		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(c)
		return
	}

	// Map data transaksi ke response
	var response = []TransactionHisotryResponse{}
	for _, trx := range trxs {
		response = append(response, trx.ToTransactionHistoryResponse())
	}

	// Kirim response sukses
	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusOK), // Ubah dari Created ke OK karena ini GET request
		infragin.WithPayload(response),
		infragin.WithMessage("get transaction histories success"),
	)
	resp.Send(c)
}

// untuk mengupdate status transaksi:
func (h handler) UpdateTransactionStatus(c *gin.Context) {
	var req struct {
		TrxId     int               `json:"trx_id" binding:"required"`
		NewStatus TransactionStatus `json:"new_status" binding:"required"`
	}

	// Bind JSON request body ke struct
	if err := c.ShouldBindJSON(&req); err != nil {
		resp := infragin.NewResponse(
			infragin.WithHttpCode(http.StatusBadRequest),
			infragin.WithMessage(err.Error()),
			infragin.WithError(response.ErrorBadRequest),
		)
		resp.Send(c)
		return
	}

	// Panggil service untuk mengupdate status transaksi
	if err := h.svc.UpdateTransactionStatus(c.Request.Context(), req.TrxId, req.NewStatus); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}

		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(c)
		return
	}

	// Kirim response sukses
	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusOK),
		infragin.WithMessage("update transaction status success"),
	)
	resp.Send(c)
}

// mendapatkan riwayat transaksi berdasarkan product_sku
func (h handler) GetTransactionHistoriesByProduct(c *gin.Context) {
	productSKU := c.Param("sku")

	// Panggil service untuk mendapatkan riwayat transaksi
	trxs, err := h.svc.GetTransactionHistoriesByProduct(c.Request.Context(), productSKU)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}

		resp := infragin.NewResponse(
			infragin.WithMessage(err.Error()),
			infragin.WithError(myErr),
		)
		resp.Send(c)
		return
	}

	// Map data transaksi ke response
	var response = []TransactionHisotryResponse{}
	for _, trx := range trxs {
		response = append(response, trx.ToTransactionHistoryResponse())
	}

	// Kirim response sukses
	resp := infragin.NewResponse(
		infragin.WithHttpCode(http.StatusOK),
		infragin.WithPayload(response),
		infragin.WithMessage("get transaction histories by product success"),
	)
	resp.Send(c)
}
