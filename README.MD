# Ecommerce Basic API

Ecommerce Basic API adalah aplikasi backend sederhana yang dibangun menggunakan Golang dan framework Gin. Aplikasi ini dirancang untuk mengelola operasi dasar e-commerce seperti autentikasi pengguna, manajemen produk, dan transaksi. Aplikasi ini menerapkan prinsip Domain-Driven Design (DDD) untuk memisahkan logika bisnis dari infrastruktur.

## Fitur
### Autentikasi Pengguna
- Registrasi pengguna baru.
- Login pengguna dengan token JWT.
- Middleware untuk memeriksa autentikasi dan role pengguna.

### Manajemen Produk
- Menambahkan produk baru (hanya admin).
- Mendapatkan daftar produk dengan paginasi.
- Mendapatkan detail produk berdasarkan SKU.

### Transaksi
- Checkout produk.
- Melihat riwayat transaksi pengguna.

### Infrastruktur
- Koneksi ke database PostgreSQL.
- Logging dengan trace ID untuk setiap request.
- Middleware untuk tracing, autentikasi, dan otorisasi.

## Teknologi yang Digunakan
- **Bahasa Pemrograman**: Golang
- **Framework Web**: Gin
- **Database**: PostgreSQL
- **ORM**: sqlx
- **JWT**: golang-jwt
- **Logging**: NooBeeID/go-logging
- **Konfigurasi**: YAML
- **Lainnya**: UUID, bcrypt, dan middleware custom

## Struktur Proyek
```
Ecommerce-basic/
├── apps/
│   ├── auth/           # Modul autentikasi
│   ├── product/        # Modul manajemen produk
│   └── transaction/    # Modul transaksi
├── cmd/
│   └── api/            # Entry point aplikasi
├── external/
│   └── database/       # Koneksi dan operasi database
├── infra/
│   ├── gin/            # Middleware dan response handler untuk Gin
│   └── response/       # Custom error response
├── internal/
│   ├── config/         # Konfigurasi aplikasi
│   └── log/            # Logging
├── utility/            # Utility functions (e.g., JWT, UUID)
├── go.mod              # Dependencies
├── go.sum              # Dependencies checksum
└── README.md           # Dokumentasi proyek
```

## Cara Menjalankan Proyek
### Prasyarat
- **Go**: Pastikan Go sudah terinstall (minimal versi 1.20).
- **PostgreSQL**: Pastikan PostgreSQL sudah terinstall dan berjalan.
- **Environment Variables**: Buat file `config.yaml` di folder `cmd/api` dengan konfigurasi berikut:

```yaml
app:
  name: Ecommerce-basic
  port: ":4000"
  encryption:
    salt: 10
    jwt_secret: "your_jwt_secret_key"

db:
  host: "localhost"
  port: "5432"
  name: "ecommerce_db"
  user: "postgres"
  password: "your_password"
  connection_pool:
    max_idle_connection: 10
    max_open_connection: 30
    max_lifetime_connection: 60
    max_idletime_connection: 60
```

### Langkah-langkah
1. **Clone Repository**
```bash
git clone https://github.com/username/Ecommerce-basic.git
cd Ecommerce-basic
```

2. **Install Dependencies**
```bash
go mod download
```

3. **Jalankan Migrasi Database (Jika Diperlukan)**
    - Buat database `ecommerce_db` di PostgreSQL.
    - Jalankan skrip SQL untuk membuat tabel yang diperlukan.

4. **Jalankan Aplikasi**
```bash
go run cmd/api/main.go
```

5. **Akses API**
    - Server berjalan di `http://localhost:4000`
    - Gunakan Postman atau curl untuk mengakses endpoint API.

## Endpoint API
### Autentikasi
#### Registrasi Pengguna
- **Method**: POST
- **Endpoint**: `/auth/register`
- **Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Login Pengguna
- **Method**: POST
- **Endpoint**: `/auth/login`
- **Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Produk
#### Mendapatkan Daftar Produk
- **Method**: GET
- **Endpoint**: `/products`
- **Query Parameters**:
    - `cursor`: ID cursor untuk paginasi (default: 0).
    - `size`: Jumlah item per halaman (default: 10).

#### Mendapatkan Detail Produk
- **Method**: GET
- **Endpoint**: `/products/sku/:sku`

#### Menambahkan Produk Baru (Admin Only)
- **Method**: POST
- **Endpoint**: `/products`
- **Headers**:
    - `Authorization`: Bearer <token>
- **Body**:
```json
{
  "name": "Product Name",
  "stock": 10,
  "price": 100000
}
```

### Transaksi
#### Checkout Produk
- **Method**: POST
- **Endpoint**: `/transactions/checkout`
- **Headers**:
    - `Authorization`: Bearer <token>
- **Body**:
```json
{
  "product_sku": "product-sku-123",
  "amount": 2
}
```

#### Melihat Riwayat Transaksi
- **Method**: GET
- **Endpoint**: `/transactions/user/histories`
- **Headers**:
    - `Authorization`: Bearer <token>

## Middleware
### Trace
- Menambahkan trace ID ke setiap request untuk logging.
- Mencatat waktu respons dan status HTTP.

### CheckAuth
- Memeriksa token JWT di header Authorization.
- Menyimpan role dan public ID pengguna di context.

### CheckRoles
- Memeriksa apakah role pengguna diizinkan untuk mengakses endpoint tertentu.

#### Contoh Penggunaan Middleware
```go
productRouter := router.Group("products")
{
    productRouter.GET("", handler.GetListProducts)
    productRouter.GET("/sku/:sku", handler.GetProductDetail)
    
    // Hanya admin yang bisa membuat produk
    productRouter.POST("",
        infragin.CheckAuth(),
        infragin.CheckRoles([]string{"admin"}),
        handler.CreateProduct,
    )
}
```

## Logging
Aplikasi ini menggunakan **NooBeeID/go-logging** untuk logging. Setiap request memiliki **trace ID** unik untuk memudahkan pelacakan.

#### Contoh log:
```
[INFO] method=GET path=/products status=200 trace_id=550e8400-e29b-41d4-a716-446655440000 duration=10ms
```

## Kontribusi
Jika ingin berkontribusi, silakan ikuti langkah berikut:
1. Fork repository ini.
2. Buat branch baru (`git checkout -b fitur-baru`).
3. Commit perubahan (`git commit -am 'Menambahkan fitur baru'`).
4. Push ke branch (`git push origin fitur-baru`).
5. Buat Pull Request.

## Lisensi
Proyek ini dilisensikan di bawah **MIT License**. Lihat file LICENSE untuk detail lebih lanjut.

## Kontak
- **Nama**: [Fathur]
- **Email**: [fathurrahman123@gmail.com]
- **GitHub**: [fahernkhan]

Selamat Menggunakan **Ecommerce Basic API**! 🚀

