package stt

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
)

type WhisperModel struct {
	model   whisper.Model
	context whisper.Context
}

// NewWhisperModel menginisialisasi model dari path yang diberikan
func NewWhisperModel(modelPath string) (*WhisperModel, error) {
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model file not found at %s", modelPath)
	}

	model, err := whisper.New(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	context, err := model.NewContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create model context: %w", err)
	}

	// Set up bahasa Indonesia
	if err := context.SetLanguage("id"); err != nil {
		log.Printf("Warning: failed to set language to 'id', fallback to auto: %v", err)
	}

	return &WhisperModel{
		model:   model,
		context: context,
	}, nil
}

// ProcessAudio memproses buffer audio 16kHz float32 menjadi string
func (w *WhisperModel) ProcessAudio(samples []float32) (string, error) {
	if len(samples) == 0 {
		return "", fmt.Errorf("empty audio samples")
	}

	// Lakukan inferensi
	if err := w.context.Process(samples, nil, nil, nil); err != nil {
		return "", fmt.Errorf("failed to process audio: %w", err)
	}

	var sb strings.Builder
	var nextSegment string
	for {
		segment, err := w.context.NextSegment()
		if err != nil {
			break
		}
		nextSegment = segment.Text
		if nextSegment != "" {
			sb.WriteString(nextSegment)
		}
	}

	return sb.String(), nil
}

// Close melepaskan resource model
func (w *WhisperModel) Close() {
	if w.context != nil {
		// handle if they have Close/Free method
	}
	if w.model != nil {
		w.model.Close()
	}
}
