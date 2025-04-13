# 
# 🚀 BankAPI - Money Transfer REST API Service : A Portfolio Project

BankAPI adalah RESTful API service untuk pengelolaan user, cabang, akun, mata uang, dan biaya transaksi.
- BankAPI menerapkan Dependency Injection, Domain Driven Design, dan Clean Architecture.
- Manajemen error dilakukan dengan membuat custom error type beserta fungsi untuk menentukan kode http status dari error tersebut.

---

## 📌 Auth Endpoints (Public)

| Method | Endpoint         | Deskripsi                   |
|--------|------------------|-----------------------------|
| POST   | `/auth/login`    | Login dan mendapatkan token |
| POST   | `/auth/refresh`  | Refresh JWT token           |

---

## 👤 User Endpoints (Private - JWT Required)

| Method | Endpoint             | Deskripsi                 |
|--------|----------------------|---------------------------|
| GET    | `/users/search`      | Cari user berdasarkan NIK |
| POST   | `/users`             | Tambah user baru          |
| GET    | `/users/:nik`        | Ambil user berdasarkan NIK|

---

## 🏢 Branch Endpoints (Private - JWT Required)

| Method | Endpoint             | Deskripsi                     |
|--------|----------------------|-------------------------------|
| GET    | `/branches`          | Ambil semua cabang            |
| GET    | `/branches/:id`      | Ambil cabang berdasarkan ID   |
| POST   | `/branches`          | Tambah cabang baru            |
| PUT    | `/branches/:id`      | Update data cabang            |
| DELETE | `/branches/:id`      | 💤 *Masih di-comment*          |

---

## 🏦 Account Type Endpoints (Private - JWT Required)

| Method | Endpoint                   | Deskripsi                      |
|--------|----------------------------|--------------------------------|
| GET    | `/account-types`           | Ambil semua jenis akun         |
| GET    | `/account-types/:id`       | Ambil jenis akun berdasarkan ID|
| POST   | `/account-types`           | Tambah jenis akun baru         |
| PUT    | `/account-types/:id`       | Update jenis akun              |
| DELETE | `/account-types/:id`       | 💤 *Masih di-comment*           |

---

## 💱 Currency Endpoints (Private - JWT Required)

| Method | Endpoint           | Deskripsi                        |
|--------|--------------------|----------------------------------|
| GET    | `/currencies`      | Ambil semua mata uang            |
| GET    | `/currencies/:id`  | Ambil mata uang berdasarkan ID   |
| POST   | `/currencies`      | Tambah mata uang baru            |
| PUT    | `/currencies/:id`  | Update data mata uang            |
| DELETE | `/currencies/:id`  | 💤 *Masih di-comment*             |

---

## 📘 Account Endpoints (Private - JWT Required)

| Method | Endpoint                                       | Deskripsi                                      |
|--------|------------------------------------------------|------------------------------------------------|
| GET    | `/accounts`                                    | Ambil semua akun                               |
| GET    | `/accounts/:accnumber`                         | Ambil akun berdasarkan nomor akun              |
| POST   | `/accounts/:branch/:acctype/:nik`              | Tambah akun berdasarkan branch, type, dan NIK  |
| PUT    | `/accounts/:accnumber`                         | Update informasi akun                          |
| DELETE | `/accounts/:id`                                | 💤 *Masih di-comment*                           |

---

## 💸 Transaction Fee Endpoints (Private - JWT Required)

| Method | Endpoint                        | Deskripsi                                |
|--------|----------------------------------|------------------------------------------|
| GET    | `/transaction-fees`             | Ambil semua biaya transaksi              |
| GET    | `/transaction-fees/:id`         | Ambil biaya transaksi berdasarkan ID     |
| POST   | `/transaction-fees`             | Tambah biaya transaksi baru              |
| PUT    | `/transaction-fees/:id`         | Update biaya transaksi                   |
| DELETE | `/transaction-fees/:id`         | Hapus biaya transaksi                    |

---

## 🛠 Cara Menjalankan Project

1. Copy `.env.example` ke `.env`:
   ```bash
   cp .env.example .env
