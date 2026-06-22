package audio

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"unsafe"

	"github.com/gen2brain/malgo"
)

type Recorder struct {
	ctx           *malgo.AllocatedContext
	device        *malgo.Device
	deviceConfig  malgo.DeviceConfig
	capturedAudio []float32
	isRecording   bool
	mu            sync.Mutex
}

func NewRecorder() (*Recorder, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		log.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		return nil, fmt.Errorf("malgo init error: %v", err)
	}

	r := &Recorder{
		ctx:           ctx,
		capturedAudio: make([]float32, 0, 16000*10), // preallocate 10 seconds of 16kHz audio
	}

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatF32
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = 16000
	deviceConfig.Alsa.NoMMap = 1

	r.deviceConfig = deviceConfig
	return r, nil
}

func (r *Recorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRecording {
		return nil
	}

	// Reset buffer
	r.capturedAudio = r.capturedAudio[:0]

	onRecvFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		r.mu.Lock()
		defer r.mu.Unlock()
		if !r.isRecording {
			return
		}

		// Convert bytes to float32
		// pInputSamples contains float32 samples
		sampleCount := len(pInputSamples) / 4
		samples := unsafe.Slice((*float32)(unsafe.Pointer(&pInputSamples[0])), sampleCount)
		r.capturedAudio = append(r.capturedAudio, samples...)
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}

	device, err := malgo.InitDevice(r.ctx.Context, r.deviceConfig, captureCallbacks)
	if err != nil {
		return fmt.Errorf("malgo init device error: %v", err)
	}
	r.device = device

	if err := r.device.Start(); err != nil {
		r.device.Uninit()
		r.device = nil
		return fmt.Errorf("failed to start device: %v", err)
	}

	r.isRecording = true
	return nil
}

func (r *Recorder) Stop() ([]float32, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isRecording {
		return nil, fmt.Errorf("not recording")
	}

	if r.device != nil {
		r.device.Uninit()
		r.device = nil
	}
	r.isRecording = false

	// Return a copy of the captured audio
	result := make([]float32, len(r.capturedAudio))
	copy(result, r.capturedAudio)
	return result, nil
}

func (r *Recorder) Close() {
	if r.ctx != nil {
		r.ctx.Free()
	}
}

// PlayBeep memainkan file wav lokal (misal beep.wav). 
// Karena keterbatasan malgo untuk WAV decoding manual, 
// kita hanya akan men-log bahwa beep dimainkan,
// atau di production bisa menggunakan package seperti `github.com/faiface/beep`
// Namun untuk zero-dependency, kita abaikan sementara atau jalankan file audio secara asinkron.
func PlayBeep(filepath string) {
	// Jika file ada, coba mainkan lewat OS call ringan
	if _, err := os.Stat(filepath); err == nil {
		if runtime.GOOS == "linux" {
			go func() { exec.Command("aplay", "-q", filepath).Run() }()
		} else if runtime.GOOS == "windows" {
			go func() {
				cmd := exec.Command("powershell", "-c", `(New-Object System.Media.SoundPlayer $env:BEEP_PATH).PlaySync()`)
				cmd.Env = append(os.Environ(), "BEEP_PATH="+filepath)
				cmd.Run()
			}()
		}
	} else {
		// Implementasi sederhana untuk notifikasi log
		log.Printf("BEEP! (Simulated playing %s, file not found)", filepath)
	}
}
