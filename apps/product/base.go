package product

import (
	"Ecommerce-basic/apps/auth"
	"Ecommerce-basic/infra/gin"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Init(router *gin.Engine, db *sqlx.DB) {
	repo := newRepository(db)
	svc := newService(repo)
	handler := newHandler(svc)

	productRoute := router.Group("/products")
	{
		productRoute.GET("", handler.GetListProducts)
		productRoute.GET("/sku/:sku", handler.GetProductDetail)

		// Authorization middleware
		authRequired := productRoute.Group("")
		authRequired.Use(infragin.CheckAuth(), infragin.CheckRoles([]string{string(auth.ROLE_Admin)}))
		{
			authRequired.POST("", handler.CreateProduct)
		}
	}
}
