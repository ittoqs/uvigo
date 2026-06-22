package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"voice-agent-id/pkg/audio"
	"voice-agent-id/pkg/keyboard"
	"voice-agent-id/pkg/stt"
	"voice-agent-id/pkg/tray"

	"github.com/getlantern/systray"
)

var (
	recorder     *audio.Recorder
	model        *stt.WhisperModel
	isProcessing bool
	procMutex    sync.Mutex
)

func main() {
	log.Println("Starting Universal Voice Input Agent...")

	// 1. Inisialisasi Clipboard
	if err := keyboard.InitClipboard(); err != nil {
		log.Fatalf("Failed to initialize clipboard: %v", err)
	}

	// 2. Inisialisasi Audio Recorder
	var err error
	recorder, err = audio.NewRecorder()
	if err != nil {
		log.Fatalf("Failed to init audio recorder: %v", err)
	}
	defer recorder.Close()

	// 3. Load Model STT
	// Dapatkan path lokasi asli dari executable ini berjalan
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	baseDir := filepath.Dir(exePath)

	// Cek model file
	modelPath := filepath.Join(baseDir, "model", "ggml-base.bin")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		modelPath = filepath.Join(baseDir, "model", "ggml-tiny.bin")
	}

	log.Printf("Loading STT model from %s...", modelPath)
	model, err = stt.NewWhisperModel(modelPath)
	if err != nil {
		log.Fatalf("Failed to load Whisper model: %v", err)
	}
	defer model.Close()

	// 4. Mulai Global Keyhook
	go keyboard.StartF4Hook(keyboard.HookCallbacks{
		OnF4Down: func() {
			procMutex.Lock()
			if isProcessing {
				procMutex.Unlock()
				return
			}
			procMutex.Unlock()

			log.Println("F4 Ditekan. Memulai perekaman...")
			tray.SetStateListening()
			
			// Mainkan Beep
			exePath, err := os.Executable()
			if err == nil {
				baseDir := filepath.Dir(exePath)
				beepPath := filepath.Join(baseDir, "model", "beep.wav")
				audio.PlayBeep(beepPath)
			}

			if err := recorder.Start(); err != nil {
				log.Printf("Failed to start recording: %v", err)
			}
		},
		OnF4Up: func() {
			procMutex.Lock()
			if isProcessing {
				procMutex.Unlock()
				return
			}
			procMutex.Unlock()

			log.Println("F4 Dilepas. Memproses suara...")
			
			audioData, err := recorder.Stop()
			if err != nil {
				log.Printf("Failed to stop recording: %v", err)
				tray.SetStateStandby()
				return
			}

			if len(audioData) == 0 {
				log.Println("No audio recorded")
				tray.SetStateStandby()
				return
			}

			procMutex.Lock()
			isProcessing = true
			procMutex.Unlock()
			tray.SetStateProcessing()

			// Proses audio di background agar tidak memblokir global keyboard hook
			go func(data []float32) {
				defer func() {
					procMutex.Lock()
					isProcessing = false
					procMutex.Unlock()
					tray.SetStateStandby()
				}()

				// Proses audio menjadi teks
				text, err := model.ProcessAudio(data)
				if err != nil {
					log.Printf("STT Processing error: %v", err)
					return
				}

				log.Printf("Transcribed Text: %s", text)
				
				// Inject ke form
				keyboard.InjectText(text)
			}(audioData)
		},
	})

	// 5. Jalankan UI Systray (Blocking)
	systray.Run(onReady, onExit)
}

func onReady() {
	ui := tray.SetupUI()

	go func() {
		for {
			select {
			case <-ui.MenuPanduan.ClickedCh:
				log.Println("Buka Panduan (TODO: open workflow.md in browser or text editor)")
			case <-ui.MenuKeluar.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func onExit() {
	log.Println("Exiting application...")
}
