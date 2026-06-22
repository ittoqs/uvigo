package keyboard

import (
	"log"
	"runtime"
	"time"

	"github.com/go-vgo/robotgo"
	"golang.design/x/clipboard"
)

func InitClipboard() error {
	return clipboard.Init()
}

// InjectText menyalin teks ke clipboard lalu menyuntikkan perintah Ctrl+V atau Cmd+V
func InjectText(text string) {
	if text == "" {
		return
	}

	// Backup isi clipboard sebelumnya (Teks maupun Gambar)
	originalText := clipboard.Read(clipboard.FmtText)
	originalImage := clipboard.Read(clipboard.FmtImage)

	// 1. Tulis ke clipboard
	clipboard.Write(clipboard.FmtText, []byte(text))

	// Memberi waktu sejenak agar clipboard system selesai melakukan update
	time.Sleep(50 * time.Millisecond)

	// 2. Simulasi pencet tombol paste
	log.Printf("Injecting text via keyboard: %s", text)
	
	if runtime.GOOS == "darwin" {
		robotgo.KeyTap("v", "command")
	} else {
		robotgo.KeyTap("v", "control")
	}

	// Memberi waktu untuk aplikasi menerima paste sebelum mengembalikan clipboard asli
	go func() {
		time.Sleep(150 * time.Millisecond)
		if len(originalImage) > 0 {
			clipboard.Write(clipboard.FmtImage, originalImage)
		} else if len(originalText) > 0 {
			clipboard.Write(clipboard.FmtText, originalText)
		} else {
			clipboard.Write(clipboard.FmtText, []byte("")) // clear
		}
	}()
}
