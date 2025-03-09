package main

import (
	"Ecommerce-basic/apps/auth"
	"Ecommerce-basic/apps/transaction"

	"Ecommerce-basic/apps/product"
	"Ecommerce-basic/external/database"
	"Ecommerce-basic/infra/gin"
	"Ecommerce-basic/internal"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load konfigurasi aplikasi
	filename := "cmd/api/config.yaml"
	if err := config.LoadConfig(filename); err != nil {
		panic(err)
	}

	// Koneksi ke database
	db, err := database.ConnectPostgres(config.Cfg.DB)
	if err != nil {
		panic(err)
	}

	if db != nil {
		log.Println("db connected")
	}

	// Gunakan semua core CPU yang tersedia untuk multi-threading
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Buat instance Gin
	router := gin.Default()

	// Middleware tracing (menggantikan infrafiber.Trace())
	router.Use(infragin.Trace())

	// Inisialisasi modul aplikasi
	auth.Init(router, db)
	product.Init(router, db)
	transaction.Init(router, db)

	// Jalankan server
	port := config.Cfg.App.Port
	log.Printf("Starting server on %s", port)
	router.Run(port)
}
