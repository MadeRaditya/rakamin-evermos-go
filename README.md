# Rakamin Evermos Virtual Internship - Backend Go

Repository ini adalah hasil pengerjaan tugas akhir (Final Task) untuk **Rakamin Evermos Virtual Internship**. Project ini merupakan backend REST API untuk simulasi platform e-commerce sederhana yang dibangun menggunakan **Golang**, **Fiber**, dan **GORM**.

---

## 📋 Daftar Isi
- [Tech Stack](#-tech-stack)
- [Fitur](#-fitur)
- [Instalasi & Cara Menjalankan](#-instalasi--cara-menjalankan)
- [Struktur Folder](#-struktur-folder)
- [API Documentation](#-api-documentation)
- [Deskripsi Soal & Ketentuan](#-deskripsi-soal--ketentuan-project)
- [Sumber Daya & Referensi](#-sumber-daya--referensi)

---

## 🛠 Tech Stack
* **Language:** [Go (Golang)](https://go.dev/)
* **Framework:** [Go Fiber v2](https://gofiber.io/)
* **ORM:** [GORM](https://gorm.io/)
* **Database:** MySQL / MariaDB
* **Auth:** JWT (JSON Web Token) & Bcrypt
* **External API:** Emsifa (Wilayah Indonesia)

---

## ✨ Fitur
1.  **Authentication:** Register (Auto create Toko) & Login (JWT).
2.  **User Management:** Update Profile, Kelola Alamat (CRUD).
3.  **Toko Management:** Update informasi toko & Upload Foto Toko.
4.  **Category:** CRUD Kategori (Hanya Admin).
5.  **Product:** CRUD Produk dengan upload multiple foto, filter pencarian, dan validasi kepemilikan.
6.  **Transaction (Trx):** Pembuatan transaksi dengan validasi stok, snapshot data produk (Log Produk), dan kalkulasi total harga.
7.  **Wilayah:** Integrasi API eksternal untuk data Provinsi dan Kota.

---

## 🚀 Instalasi & Cara Menjalankan

### Prasyarat
* Go (versi 1.18 ke atas)
* MySQL Server (via XAMPP atau Docker)
* Git

### Langkah-langkah
1.  **Clone Repository**
    ```bash
    git clone [github.com/MadeRaditya/rakamin-evermos-go.git](github.com/MadeRaditya/rakamin-evermos-go.git)
    cd rakamin-evermos-go
    ```

2.  **Install Dependencies**
    ```bash
    go mod tidy
    ```

3.  **Setup Database**
    * Buat database baru di MySQL dengan nama: `rakamin_evermos_go`.
    * Pastikan konfigurasi username/password database sudah sesuai (lihat langkah selanjutnya).
    * *Catatan:* Kamu **tidak perlu** meng-import tabel secara manual. Aplikasi ini menggunakan **Auto Migration**, sehingga tabel akan otomatis dibuat saat server dijalankan pertama kali.

4.  **Setup Environment Variables**
    Duplikat file `.env.example` menjadi `.env`, lalu sesuaikan isinya dengan konfigurasi lokal kamu (seperti DB Password atau JWT Secret).
    
    **Linux/Mac:**
    ```bash
    cp .env.example .env
    ```
    
    **Windows (CMD/Powershell):**
    ```cmd
    copy .env.example .env
    ```
    
    *Pastikan isi `.env` sudah benar, contoh:*
    ```env
    APP_PORT=8000

    DB_DRIVER=mysql
    DB_HOST=127.0.0.1
    DB_PORT=3306
    DB_USER=root
    DB_PASSWORD=
    DB_NAME=rakamin_evermos_go

    DB_CHARSET=utf8mb4
    DB_PARSE_TIME=true
    DB_LOC=Local

    JWT_SECRET=rahasia_negara_bos
    
    ```

5.  **Jalankan Server**
    ```bash
    go run main.go
    ```
    Server akan berjalan di `http://localhost:8000`.

---

## 📂 Struktur Folder
Menerapkan pendekatan *Clean Architecture* sederhana:

```text
├── controllers/    # Handler logic untuk setiap endpoint
├── database/       # Konfigurasi koneksi database & migration
├── middlewares/    # Middleware (JWT Auth)
├── models/         # Struct entity database (GORM)
├── routes/         # Definisi URL endpoint
├── public/         # Static files (Uploads foto)
├── main.go         # Entry point aplikasi
└── go.mod          # Dependency manager

```

---

## 📖 API Documentation (Postman)

Sesuai ketentuan, routing API mengikuti standar Postman Collection dari Rakamin.

* **Base URL:** `http://localhost:8000/api/v1`
* **Postman Collection:** [Download JSON Collection](https://www.google.com/search?q=https://github.com/Fajar-Islami/go-example-cruid/blob/master/Rakamin%2520Evermos%2520Virtual%2520Internship.postman_collection.json)

---

## 📝 Deskripsi Soal & Ketentuan Project

*Bagian ini disalin langsung dari dokumen tugas.*

Soal 

1. Buatlah service **Login** dan **Register**.


2. Ketika berhasil register, **toko otomatis terbuat**.


3. Service untuk mengelola akun (**User**).


4. Buatlah service **Toko**.


5. Buatlah service **Alamat**.


6. Buatlah service **Kategori**. Kategori hanya dapat dikelola oleh **Admin**, untuk itu kalian harus mengubah status Admin di database langsung.


7. Buatlah service **Produk**.


8. Buatlah service **Transaksi**.



Ketentuan 

1. Harus memiliki routing seperti collection berikut: `Rakamin Evermos Virtual Internship.postman collection.json`.


2. Boleh menambahkan dari API yang sudah tapi tidak boleh dikurangi.


3. Email dan no telepon user tidak boleh ada yang sama (**Unique**).


4. Menggunakan **JWT**.


5. Harus terdapat API yang **meng-upload file**.


6. Toko otomatis terbuat ketika user mendaftar.


7. Alamat diperlukan untuk alamat kirim produk.


8. Yang dapat mengelola kategori hanyalah **admin**.


9. Menerapkan **pagination** seperti di postman.


10. Menerapkan **filtering data**.


11. User tidak dapat mendapatkan data user lain atau meng-update user lain.


12. User tidak dapat mengelola alamat data user lain.


13. User tidak dapat mengelola data toko dari data user lain.


14. User tidak dapat mengelola data product dari data user lain.


15. User tidak dapat mengelola data transaksi dari data user lain.


16. Tabel log product diisi ketika melakukan transaksi.


17. Tabel log produk digunakan untuk menyimpan data produk yang ada di transaksi (**Snapshot**).


18. Menerapkan **clean architecture**.



---

🔗 Sumber Daya & Referensi 

Berikut adalah sumber daya yang disediakan:

1. **API Wilayah Indonesia**
Untuk mendapatkan data wilayah:


[https://www.emsifa.com/api-wilayah-indonesia/](https://www.emsifa.com/api-wilayah-indonesia/) 


2. **Konfigurasi MySQL (Referensi GORM & Mux)**
Untuk konfigurasi MySQL dapat dilihat di:


[Build a REST API using Go, MySQL, GORM](https://levelup.gitconnected.com/build-a-rest-api-using-go-mysql-gorm-and-mux-a02e9a2865ee) 


3. **Desain Database (ERD)**
Berikut design database:


[Link Google Drive](https://www.google.com/search?q=https://drive.google.com/file/d/1L7pzFNjMNfUU-f3tsrZPKTjMpblGaodW/view%3Fusp%3Dsharing) 



---

*Dibuat sebagai bagian dari Final Task Rakamin Evermos Virtual Internship.*

```
