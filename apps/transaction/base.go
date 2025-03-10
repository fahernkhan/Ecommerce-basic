package transaction

import (
	"Ecommerce-basic/infra/gin"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Init(router *gin.Engine, db *sqlx.DB) {
	repo := newRepository(db)
	svc := newService(repo)
	handler := newHandler(svc)

	trxRoute := router.Group("transactions")
	{
		// menggunakan middleware
		trxRoute.Use(infragin.CheckAuth())

		// route dibawahnya akan menggunakan middleware tersebut
		trxRoute.POST("/checkout", handler.CreateTransaction)
		trxRoute.GET("/user/histories", handler.GetTransactionByUser)
		trxRoute.PUT("/status", handler.UpdateTransactionStatus)
		trxRoute.GET("/product/:sku/histories", handler.GetTransactionHistoriesByProduct)
	}
}
