package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Init(router *gin.Engine, db *sqlx.DB) {
	repo := newRepository(db)
	svc := newService(repo)
	handler := newHandler(svc)

	authRouter := router.Group("auth")
	{
		authRouter.POST("register", handler.register)
		authRouter.POST("login", handler.login)
	}
}
