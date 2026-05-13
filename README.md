# Madlab Toolkits API (Backend)

Halo! Selamat datang di repo backend **Madlab Toolkits**. Singkatnya, ini adalah API serbaguna yang dibangun murni menggunakan Golang. Gak ada framework web tebal kayak Gin atau Fiber di sini, semuanya jalan di atas _Standard Library_ (`net/http`) bawaan Go.

Tujuannya simpel: bikin backend yang enteng, kenceng, minim _dependency_, dan gampang di-maintain buat ngurusin manipulasi gambar dan download media.

---

## Spesifikasi Mesin

- **Bahasa**: Golang versi `1.25.0`
- **Routing**: Pakai bawaan Go `net/http` (sudah pakai fitur routing _method_ HTTP yang ada sejak Go 1.22)
- **Struktur Kode**: Dibikin rapi dan modular (`handlers`, `middlewares`, `services`, dll) biar gampang kalau mau nambah fitur baru ke depannya.

---

## Library yang Dipakai

Walaupun nggak pakai framework utama, project ini tetep butuh beberapa alat bantu khusus buat ngerjain tugas-tugas berat:

- [`disintegration/imaging`](https://github.com/disintegration/imaging): Buat ngurusin urusan potong-memotong, _resize_, dan manipulasi resolusi gambar.
- [`chai2010/webp`](https://github.com/chai2010/webp): Biar backend kita jago nge-_encode_ dan _decode_ gambar berformat WebP.
- [`joho/godotenv`](https://github.com/joho/godotenv): Buat baca konfigurasi dari file `.env` biar aman dan gampang diatur tiap pindah _environment_.
- [`golang.org/x/time`](https://pkg.go.dev/golang.org/x/time): Library bawaan tambahan dari Go buat bikin fitur _Rate Limiting_ (biar API kita gak gampang di-spam sama orang).

---

## Tech Stack & Deployment

- **Deployment**: Bisa jalan di VPS biasa dan Render.
- **Container**: Udah dibungkus pakai Docker (`Dockerfile` udah rapi), jadi tinggal _build and run_ di mana aja tanpa ribet ngurus _environment_.
- **Tools Tambahan**: Fitur downloader di API ini bergantung banget sama **`yt-dlp`**. Jadi pastikan `yt-dlp` udah ter-install di server/Docker-nya buat nyedot video/audio dari sosmed kayak YouTube atau IG.

---

## Fitur yang Udah Dibikin

Ada dua kelompok fitur utama yang bisa dipakai di API ini:

### 1. Image Toolkit (Tukang Gambar)

Kumpulan endpoint buat ngoprek gambar. Semuanya butuh dikirim lewat format form upload (`multipart/form-data`).

- **Convert Image** (`POST /api/v1/image/convert`): Tinggal upload gambar, terus pilih mau diubah formatnya jadi `JPG`, `PNG`, atau `WebP`.
- **Compress Image** (`POST /api/v1/image/compress-image`): Buat ngecilin ukuran (_size_) file gambar. Kualitasnya (`quality`) bisa diatur manual dari angka 1 sampai 100.
- **Resize Image** (`POST /api/v1/image/resize-image`): Buat ngubah dimensi/resolusi gambar. Tinggal masukin aja `width` (lebar) dan `height` (tinggi) yang dimau.

### 2. Media Downloader (Tukang Sedot Media)

Endpoint buat narik video/audio dari platform luar.

- **Get Info** (`POST /api/v1/downloader/info`): Buat ngintip detail video dari URL (bakal dapet info judul, _thumbnail_, sama daftar resolusi yang tersedia) memanfaatkan _engine_ `yt-dlp`.
- **Download File** (`POST /api/v1/downloader/download`): Buat langsung nyedot file dari URL dengan masukin `format_id`-nya. Nanti filenya bakal otomatis ke-download ke perangkat pengguna.

---

## Extra Keamanan (Custom Middleware)

Karena nggak pake framework web, semua _security layer_ dan _middleware_ dibikin manual pakai Go murni (dan ini seru banget!):

- **CORS Middleware**: Biar frontend tetep bisa narik data dengan aman tanpa kena blokir sama browser.
- **Rate Limiter**: IP yang terlalu barbar bakal dibatesin jumlah _request_-nya biar server nggak _down_ gara-gara terlalu sibuk.
- **Timeout Middleware**: Kalau ada proses yang nyangkut dan kelamaan (misal download gagal putus di tengah jalan), koneksinya otomatis digagalkan.
- **Logger**: Semua API yang di-hit bakal otomatis dicatet di terminal beserta waktu prosesnya.
- **Auto Cleanup**: Punya fungsi _job_ kecil di _background_ yang rajin ngehapusin file-file sampah hasil sisa kompres/download biar _storage_ server nggak gampang penuh.
