package infragin

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"Ecommerce-basic/infra/response"
	"Ecommerce-basic/internal"
	infraLog "Ecommerce-basic/internal/log"
	"Ecommerce-basic/utility"
	"github.com/NooBeeID/go-logging/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Custom ResponseWriter untuk capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *strings.Builder
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b) // Simpan response body
	return rw.ResponseWriter.Write(b)
}

// Middleware Trace
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		now := time.Now()
		traceId := uuid.New().String()
		c.Set("X-Trace-ID", traceId)

		// Data log awal
		data := map[logger.LogKey]interface{}{
			logger.TRACER_ID: traceId,
			logger.METHOD:    c.Request.Method,
			logger.PATH:      c.Request.URL.Path,
		}

		// Set log data di context
		ctx = context.WithValue(ctx, logger.DATA, data)
		infraLog.Log.Infof(ctx, "incoming request")

		// Update context di request
		c.Request = c.Request.WithContext(ctx)

		// Bungkus response writer
		rw := &responseWriter{ResponseWriter: c.Writer, body: &strings.Builder{}}
		c.Writer = rw

		// Lanjutkan request
		c.Next()

		// Tambahkan response time
		data[logger.RESPONSE_TIME] = time.Since(now).Milliseconds()
		data[logger.RESPONSE_TYPE] = "ms"

		httpStatusCode := c.Writer.Status()
		if httpStatusCode >= 200 && httpStatusCode <= 299 {
			infraLog.Log.Infof(ctx, "success")
		} else {
			data["response_body"] = rw.body.String() // Ambil response body
			infraLog.Log.Errorf(ctx, "failed")
		}
	}
}

// CheckAuth Middleware
func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":   false,
				"error":     response.ErrorUnauthorized.Message,
				"errorCode": response.ErrorUnauthorized.Code,
			})
			return
		}

		// Extract token from "Bearer <token>"
		bearer := strings.Split(authorization, "Bearer ")
		if len(bearer) != 2 {
			log.Println("token invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":   false,
				"error":     response.ErrorUnauthorized.Message,
				"errorCode": response.ErrorUnauthorized.Code,
			})
			return
		}

		token := bearer[1]

		// Validate token
		publicId, role, err := utility.ValidateToken(token, config.Cfg.App.Encryption.JWTSecret)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":   false,
				"error":     response.ErrorUnauthorized.Message,
				"errorCode": response.ErrorUnauthorized.Code,
			})
			return
		}

		// Set role and public ID in context
		c.Set("ROLE", role)
		c.Set("PUBLIC_ID", publicId)

		c.Next()
	}
}

// CheckRoles Middleware
func CheckRoles(authorizedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := fmt.Sprintf("%v", c.GetString("ROLE"))

		// Check if role is authorized
		isExists := false
		for _, authorizedRole := range authorizedRoles {
			if role == authorizedRole {
				isExists = true
				break
			}
		}

		if !isExists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success":   false,
				"error":     response.ErrorForbiddenAccess.Message,
				"errorCode": response.ErrorForbiddenAccess.Code,
			})
			return
		}

		c.Next()
	}
}
