package auth

import (
	"Ecommerce-basic/infra/response"
	"net/http"

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

func (h handler) register(c *gin.Context) {
	var req RegisterRequestPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":   false,
			"message":   "register fail",
			"error":     err.Error(),
			"errorCode": response.ErrorBadRequest.Code,
		})
		return
	}

	if err := h.svc.register(c.Request.Context(), req); err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}

		c.JSON(myErr.HttpCode, gin.H{
			"success":   false,
			"message":   err.Error(),
			"error":     myErr.Message,
			"errorCode": myErr.Code,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "register success",
	})
}

func (h handler) login(c *gin.Context) {
	var req LoginRequestPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":   false,
			"message":   "login fail",
			"error":     err.Error(),
			"errorCode": response.ErrorBadRequest.Code,
		})
		return
	}

	token, err := h.svc.login(c.Request.Context(), req)
	if err != nil {
		myErr, ok := response.ErrorMapping[err.Error()]
		if !ok {
			myErr = response.ErrorGeneral
		}

		c.JSON(myErr.HttpCode, gin.H{
			"success":   false,
			"message":   err.Error(),
			"error":     myErr.Message,
			"errorCode": myErr.Code,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":      true,
		"message":      "login success",
		"access_token": token,
	})
}
