# UVIGO
# Universal Voice Input Agent (Bahasa Indonesia)

Universal Voice Input Agent berbasis Whisper.cpp yang memungkinkan konversi ucapan ke teks secara lokal (offline) dengan dukungan distribusi untuk Linux dan Windows.

* Transkripsi suara ke teks secara lokal (offline)
* Menggunakan mesin pengenalan suara Whisper.cpp
* Berjalan di system tray
* Dukungan Linux dan Windows
* Distribusi ringan tanpa ketergantungan cloud

# Kompilasi (Build)

## A. Linux

### 1. Install Dependencies

```bash
sudo apt-get update

sudo apt-get install \
gcc \
g++ \
libx11-dev \
xorg-dev \
libxtst-dev \
libpng-dev \
libasound2-dev \
libxkbcommon-x11-dev \
libayatana-appindicator3-dev \
libgtk-3-dev
```

### 2. Build Library Whisper.cpp

Library whisper.cpp/bindings/go mengharuskan file libwhisper.a sudah tersedia.

Ikuti dokumentasi resmi:

https://github.com/ggml-org/whisper.cpp/tree/master/bindings/go

### 3. Build Binary Linux

```bash
CGO_ENABLED=1 \
CGO_LDFLAGS="-L/path/ke/whisper.cpp -lwhisper" \
CGO_CFLAGS="-I/path/ke/whisper.cpp" \
go build -o voice-agent main.go
```

Ganti `/path/ke/whisper.cpp` sesuai lokasi hasil kompilasi library Whisper.cpp.

### 4. Membuat Paket Debian

Struktur direktori:

```bash
mkdir -p voice-agent_1.0.0_amd64/usr/local/bin

mkdir -p voice-agent_1.0.0_amd64/usr/local/share/voice-agent/model

mkdir -p voice-agent_1.0.0_amd64/DEBIAN
```

Salin file aplikasi:

```bash
cp voice-agent \
voice-agent_1.0.0_amd64/usr/local/bin/

cp model/ggml-base.bin \
voice-agent_1.0.0_amd64/usr/local/share/voice-agent/model/
```

Buat file:

```text
voice-agent_1.0.0_amd64/DEBIAN/control
```

Isi file:

```text
Package: voice-agent-id

Version: 1.0.0

Architecture: amd64

Maintainer: Nama Anda

Description: Universal Voice Input Agent untuk Linux.
```

Build paket Debian:

```bash
dpkg-deb --build voice-agent_1.0.0_amd64
```

Hasil:

```text
voice-agent_1.0.0_amd64.deb
```

Instalasi:

```bash
sudo dpkg -i voice-agent_1.0.0_amd64.deb
```

## B. Windows (Cross-Compilation dari Linux)

### 1. Install MinGW

```bash
sudo apt-get install mingw-w64
```

### 2. Build Binary Windows

```bash
CC=x86_64-w64-mingw32-gcc \
CXX=x86_64-w64-mingw32-g++ \
CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
CGO_LDFLAGS="-L/path/ke/whisper.cpp -lwhisper" \
CGO_CFLAGS="-I/path/ke/whisper.cpp" \
go build \
-ldflags "-H=windowsgui" \
-o voice-agent.exe \
main.go
```

### Catatan

Flag berikut:

```bash
-ldflags "-H=windowsgui"
```

digunakan agar aplikasi berjalan di System Tray tanpa menampilkan jendela konsol (Command Prompt).

Saat distribusi, pastikan struktur file seperti berikut:

```text
voice-agent.exe

model/
└── ggml-base.bin
```

Didistribusikan dalam format:

```text
voice-agent-windows.zip
```

# Distribusi

## Linux

Instal menggunakan:

```bash
sudo dpkg -i voice-agent_1.0.0_amd64.deb
```

## Windows

Distribusikan sebagai file ZIP yang berisi:

```text
voice-agent.exe

model/
```

# Whisper.cpp

Penggunaan Whisper.cpp sebagai mesin Speech-to-Text mengikuti ketentuan masing-masing dependency yang digunakan.
