package tray

import (
	"github.com/getlantern/systray"
)

type UI struct {
	MenuPanduan *systray.MenuItem
	MenuKeluar  *systray.MenuItem
}

// SetupUI dipanggil saat systray siap. Mengembalikan referensi menu untuk diikat dengan event handler
func SetupUI() *UI {
	systray.SetTitle("Siap")
	systray.SetTooltip("Universal Voice Input Agent")
	
	// Secara ideal Anda akan menaruh byte array icon di systray.SetIcon(iconBytes)
	// Untuk saat ini kita gunakan text title saja jika icon tidak tersedia.

	panduan := systray.AddMenuItem("Panduan", "Cara menggunakan aplikasi")
	keluar := systray.AddMenuItem("Keluar", "Tutup aplikasi")

	return &UI{
		MenuPanduan: panduan,
		MenuKeluar:  keluar,
	}
}

// SetStateStandby mengubah tampilan tray menjadi Siap
func SetStateStandby() {
	systray.SetTitle("Siap")
	systray.SetTooltip("Tekan F4 untuk bicara")
}

// SetStateListening mengubah tampilan tray menjadi Mendengarkan
func SetStateListening() {
	systray.SetTitle("Mendengarkan")
	systray.SetTooltip("Sedang mendengarkan suara Anda...")
}

// SetStateProcessing mengubah tampilan tray menjadi Memproses
func SetStateProcessing() {
	systray.SetTitle("Memproses")
	systray.SetTooltip("Sedang mentranskripsi suara ke teks...")
}
