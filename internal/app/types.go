package app

import (
	"context"
	"io"
)

type ExtractedData struct {
	Summary      string        `json:"summary,omitempty"`
	Participants []string      `json:"participants,omitempty"`
	Decisions    []string      `json:"decisions,omitempty"`
	ActionItems  []ActionItems `json:"actionItems,omitempty"`
}

type ActionItems struct {
	Owner string `json:"owner,omitempty"`
	Task  string `json:"task,omitempty"`
	Due   string `json:"due,omitempty"`
}

type TranscribeCommand struct {
	ID string
}

type ExtractCommand struct {
	ID string
}

type ExportCommand struct {
	ID string
}

type RecordingStorer interface {
	StoreRecording(file []byte) (string, error)
}

type RecordingTranscriptionStore interface {
	OpenRecording(id string) (io.ReadCloser, error)
	SaveTranscription(id string, transcription []byte) error
	SaveError(id string, worker string, message string) error
}

type RecordingExtractionStore interface {
	GetTranscription(id string) (string, error)
	SaveExtraction(id string, extraction []byte) error
	SaveError(id string, worker string, message string) error
}

type RecordingExportStore interface {
	GetExtraction(id string) (ExtractedData, error)
	SaveError(id string, worker string, message string) error
}

type Transcriber interface {
	Transcribe(ctx context.Context, audio io.Reader) (string, error)
}

type Extractor interface {
	Extract(ctx context.Context, transcription string) (ExtractedData, error)
}

type Exporter interface {
	Export(ctx context.Context, data ExtractedData) error
}
