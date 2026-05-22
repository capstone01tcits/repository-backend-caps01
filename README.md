# Sevima AI Video Gen - Backend

Aplikasi backend untuk platform pembuatan video pembelajaran otomatis berbasis AI. Proyek ini dibangun menggunakan **Go (Fiber)** dan **PostgreSQL (GORM)**.

---

## 🚀 Prasyarat

Pastikan sudah terinstall:
- Go >= 1.20
- PostgreSQL
- Python 3.10+ (Untuk AI Service microservice)

---

## 📦 Instalasi dan Konfigurasi

1. **Clone repository:**
   ```bash
   git clone <URL_REPO_BE>
   cd Capstone-01-Backend
   ```

2. **Konfigurasi Database & Environment:**
   Buat file `.env` di root proyek:
   ```env
   APP_PORT=5000
   APP_ENV=development
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=admin
   DB_NAME=sevima_video_db
   JWT_SECRET=super_secret_jwt_key
   AI_SERVICE_URL=http://localhost:8000
   ```

3. **Jalankan Aplikasi:**
   ```bash
   go mod tidy
   go run cmd/main.go
   ```

4. **Jalankan AI Service (Python):**
   ```bash
   cd ai-service
   pip install -r requirements.txt
   python main.py
   ```

---

## 🏗️ Fitur Terkini

- **Autentikasi JWT**: Registrasi, Login, Refresh Token, dan Get Profile.
- **Role-Based Access Control**: Role `admin` dan `user`.
- **Manajemen Proyek & Storyboard**: Pembuatan proyek secara atomik, soft-delete, dan restore.
- **Integrasi Video Generation**: Sinkronisasi dengan AI Python Microservice (Google Veo 3.1 Lite & LTX).
- **Background Worker**: Proses antrean pembuatan video berjalan di latar belakang (goroutines).
- **Sistem Kredit**: Fitur top-up kredit khusus Admin (`/api/admin/credits`) dan lihat daftar pengguna (`/api/admin/users`).

Dokumentasi API lengkap tersedia di folder `docs/`.
