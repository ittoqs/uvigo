package keyboard

import (
	"log"

	hook "github.com/robotn/gohook"
)

type HookCallbacks struct {
	OnF4Down func()
	OnF4Up   func()
}

// StartF4Hook memulai pemantauan global key secara blocking.
// Panggil menggunakan goroutine jika perlu berjalan di background.
func StartF4Hook(callbacks HookCallbacks) {
	log.Println("Global keyboard hook started...")
	
	evChan := hook.Start()
	defer hook.End()

	isF4Pressed := false

	for ev := range evChan {
		if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyUp {
			// Keycode untuk F4 di gohook adalah 62
			// Anda dapat memeriksa keycode spesifik dengan mencetak ev.Rawcode
			if ev.Rawcode == 62 || ev.Keychar == 62 {
				if ev.Kind == hook.KeyDown {
					if !isF4Pressed {
						isF4Pressed = true
						if callbacks.OnF4Down != nil {
							callbacks.OnF4Down()
						}
					}
				} else if ev.Kind == hook.KeyUp {
					if isF4Pressed {
						isF4Pressed = false
						if callbacks.OnF4Up != nil {
							callbacks.OnF4Up()
						}
					}
				}
			}
		}
	}
}
