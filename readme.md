# 🟢 upgak

CLI sederhana untuk mengecek apakah suatu website *up* atau *down*.  
Dibangun dengan Go, ringan dan mudah digunakan.

## ✨ Fitur

- ✅ **Cek Status Website**: Menentukan apakah website *up* (HTTP 200–399) atau *down* (HTTP ≥400, timeout, atau error koneksi)
- ✅ **Tampilkan Info Detail**: Menampilkan kode response HTTP dan waktu respons dalam milliseconds
- ✅ **Mode Eksekusi**: Pilihan menggunakan goroutine (concurrent) untuk performa lebih cepat atau mode serial

## 🚀 Penggunaan

### Instalasi
```bash
git clone https://github.com/akbarhlubis/upgak.git
cd upgak
go build ./cmd/upgak
```

### Contoh Penggunaan

**Cek satu website:**
```bash
./upgak -urls="https://github.com"
```

**Cek beberapa website secara serial:**
```bash
./upgak -urls="https://google.com,https://github.com,https://stackoverflow.com"
```

**Cek beberapa website secara concurrent (lebih cepat):**
```bash
./upgak -urls="https://google.com,https://github.com,https://stackoverflow.com" -concurrent
```

**Dengan timeout custom:**
```bash
./upgak -urls="https://github.com" -timeout=5
```

### Opsi Command Line

- `-urls string`: Daftar URL yang dipisahkan koma (wajib)
- `-concurrent`: Gunakan goroutines untuk checking parallel (default: false)
- `-timeout int`: Timeout dalam detik untuk setiap request (default: 10)
- `-help`: Tampilkan pesan bantuan

### Contoh Output

```
🔍 Checking 3 website(s) with concurrent mode (timeout: 10s)

✅ https://github.com - UP - HTTP 200 - Response time: 203 ms
✅ https://google.com - UP - HTTP 200 - Response time: 156 ms
❌ https://httpstat.us/500 - DOWN - HTTP 500 - Response time: 842 ms

📊 Summary: 2/3 websites are UP
⏱️  Total checking time: 845ms
```

## ⚡ Performa

Mode concurrent memberikan peningkatan performa signifikan saat mengecek banyak website:

- **Serial mode**: ~918ms untuk 5 URL (satu per satu)
- **Concurrent mode**: ~419ms untuk 5 URL (parallel) - **54% lebih cepat**

## 🗂️ Struktur Projek

```
upgak/
├── cmd/upgak           # Entry point CLI (main.go)
├── internal/checker    # Logika pengecekan website
├── go.mod / go.sum     # Module dan dependency
├── README.md           # Dokumentasi ini
└── LICENSE             # Lisensi project
```

## 🧪 Testing

Jalankan test:
```bash
go test ./...
```

Test mencakup:
- Pengecekan website valid/invalid
- Logika penentuan status berdasarkan HTTP code
- Kedua mode eksekusi (serial vs concurrent)

## Lisensi
MIT 2025