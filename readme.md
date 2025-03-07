# News API - Golang

Simple News API adalah aplikasi berbasis Golang yang menggunakan database MySQL untuk menyimpan dan mengelola berita.

---

## 1. Persiapan Database

### 1.1 Login ke MySQL
```sh
mysql -u root -p
```

### 1.2 Start MySQL Database
```sh
# Untuk Linux
sudo systemctl start mysql  

# Untuk Windows
net start MySQL             
```

### 1.3 Buat Database
```sql
CREATE DATABASE newsdb;
```

### 1.4 Cek Database yang Ada
```sql
SHOW DATABASES;
```

### 1.5 Gunakan Database `newsdb`
```sql
USE newsdb;
```

### 1.6 Buat Tabel `news`
```sql
CREATE TABLE news (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

## 2. Menjalankan Redis

### 2.1 Start Redis
```sh
# Untuk Linux
sudo systemctl start redis  

# Untuk Mac & Windows
redis-server                
```

### 2.2 Login ke Redis
```sh
redis-cli
```

---

## 3. Menjalankan Aplikasi Golang

### 3.1 Clone Repository
```sh
git clone https://github.com/kusnadi8605/news.git
cd news
```

### 3.2 Inisialisasi Modul Golang
```sh
go mod init github.com/kusnadi8605/news
go mod tidy
go mod vendor
```

### 3.3 Konfigurasi Koneksi Database dan Redis
Buat file `.env` dan sesuaikan konfigurasi berikut:
```
# MySQL Configuration
DB_USER=root
DB_PASSWORD=root
DB_HOST=localhost
DB_PORT=3306
DB_NAME=newsdb

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
```

### 3.4 Jalankan Aplikasi
```sh
go run main.go
```

---

## 4. API Endpoint

### 4.1 Get All News
```sh
curl --location 'http://localhost:8080/news'
```

### 4.2 Get News by ID
```sh
curl --location 'http://localhost:8080/news/1'
```

### 4.3 Create News
```sh
curl --location --request POST 'http://localhost:8080/news' \
--header 'Content-Type: application/json' \
--data '{
    "title": "Hello Golang",
    "content": "Golang adalah bahasa yang hebat!"
}'
```

### 4.4 Update News
```sh
curl --location --request PUT 'http://localhost:8080/news/1' \
--header 'Content-Type: application/json' \
--data '{
    "title": "Hello Golang & Java",
    "content": "Golang & Java adalah bahasa yang hebat!"
}'
```

### 4.5 Delete News
```sh
curl --location --request DELETE 'http://localhost:8080/news/1'
```

---

## 5. Teknologi yang Digunakan
- **Golang** sebagai bahasa pemrograman utama
- **Echo** sebagai framework web
- **MySQL** sebagai database
- **Redis** untuk caching


