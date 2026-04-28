# Changelog

Semua perubahan penting di project ini didokumentasikan di file ini.

Format berdasarkan [Keep a Changelog](https://keepachangelog.com/id-ID/1.0.0/).

---

## [1.1.0] - 2026-04-29

### Added

- **Social Media Downloader** – Fitur untuk mengunduh media dari URL (YouTube, TikTok, Instagram) menggunakan `yt-dlp`.
  - Endpoint `POST /api/v1/downloader` untuk mengunduh video/audio berdasarkan `format_id`.
  - Endpoint `POST /api/v1/downloader/info` untuk mengambil metadata/info (resolusi, ekstensi, dsb.) sebelum mengunduh.
- **Sistem Caching** (`DownloadCache`): Cache in-memory untuk menyimpan path file hasil unduhan, mempercepat request berikutnya untuk video dan `format_id` yang sama.
- Support file cookies (`youtube_cookies.txt`, `ig-cookies.txt`) pada eksekusi yt-dlp untuk pengunduhan yang optimal.
- Batas konkurensi (semaphore) pada fitur downloader untuk menjaga kestabilan server.

### Changed

- Parameter form pengubahan gambar diubah dari `targetFormat` menjadi `format` untuk kejelasan (`image_handler.go`).
- Perbaikan pengelolaan `.env` (dihapus dari repo, ditambahkan ke `.gitignore`) dan default port di `main.go`.
- Versi Go diubah dari 1.25.5 ke 1.24.0 di `go.mod`.

---

## [1.0.0] - 2025-02-14

### Added

- **Image Convert** – `POST /api/v1/image/convert`: konversi format gambar (JPG, PNG, WebP).
- **Image Compress** – `POST /api/v1/image/compress-image`: kompresi dengan parameter quality (1–100).
- **Image Resize** – `POST /api/v1/image/resize-image`: resize dengan width/height (maks. 4096 px).
- Utils `LoadImage` / `SaveImage` untuk decode & encode gambar (DRY di semua fitur image).
- Storage cleanup job: hapus file lama di `temp/uploads`, `temp/processed`, `temp/compressed`, `temp/resized` (usia > 24 jam, jalan tiap 5 detik).
- CORS konfigurasi dari env dengan default aman (ALLOWED_ORIGINS, ALLOWED_METHODS, ALLOWED_HEADERS).
- Validasi: ukuran file max 5MB, quality 1–100, width/height max 4096 px.
- Batas konkurensi: maks. 5 proses image bersamaan; request tambahan dapat 503 "Server sibuk".

### Changed

- Cleanup file upload: file di `temp/uploads` dihapus setelah proses (sukses atau gagal) di Convert, Compress, dan Resize.
- Storage cleaner memakai `filepath.WalkDir` (menggantikan `filepath.Walk`) untuk efisiensi.

### Fixed

- CORS: default origin tanpa trailing slash; methods/headers punya default jika env kosong.
- Pembersihan folder hasil: `temp/compressed` dan `temp/resized` ikut dibersihkan oleh job (sebelumnya hanya processed & uploads).
