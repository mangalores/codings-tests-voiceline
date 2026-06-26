package repository

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mangalores/case-studies-voiceline/internal/app"
)

var ErrNotFound = errors.New("recording not found")

type RecordingRepository struct {
	rootPath string
}

func NewRecordingRepository(rootPath string) *RecordingRepository {
	return &RecordingRepository{rootPath: rootPath}
}

func (r *RecordingRepository) StoreRecording(file []byte) (string, error) {
	if err := os.MkdirAll(r.rootPath, 0o755); err != nil {
		return "", err
	}

	entries, err := os.ReadDir(r.rootPath)
	if err != nil {
		return "", err
	}

	id := strconv.Itoa(countDirectories(entries) + 1)
	folderPath := filepath.Join(r.rootPath, id)
	if err := os.Mkdir(folderPath, 0o755); err != nil {
		return "", err
	}

	if err := os.WriteFile(filepath.Join(folderPath, "recording.mp3"), file, 0o644); err != nil {
		_ = os.RemoveAll(folderPath)
		return "", err
	}

	return id, nil
}

func (r *RecordingRepository) SaveTranscription(id string, transcription []byte) error {
	folderPath, err := r.recordingFolder(id)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(folderPath, "transcription.md"), transcription, 0o644)
}

func (r *RecordingRepository) OpenRecording(id string) (io.ReadCloser, error) {
	folderPath, err := r.recordingFolder(id)
	if err != nil {
		return nil, err
	}

	return os.Open(filepath.Join(folderPath, "recording.mp3"))
}

func (r *RecordingRepository) SaveExtraction(id string, extraction []byte) error {
	folderPath, err := r.recordingFolder(id)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(folderPath, "extracted.json"), extraction, 0o644)
}

func (r *RecordingRepository) SaveError(id string, worker string, message string) error {
	folderPath, err := r.recordingFolder(id)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(struct {
		Worker  string `json:"worker"`
		Message string `json:"message"`
	}{
		Worker:  worker,
		Message: message,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(folderPath, "error.json"), payload, 0o644)
}

func (r *RecordingRepository) GetTranscription(id string) (string, error) {
	folderPath, err := r.recordingFolder(id)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(filepath.Join(folderPath, "transcription.md"))
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (r *RecordingRepository) GetExtraction(id string) (app.ExtractedData, error) {
	folderPath, err := r.recordingFolder(id)
	if err != nil {
		return app.ExtractedData{}, err
	}

	content, err := os.ReadFile(filepath.Join(folderPath, "extracted.json"))
	if err != nil {
		return app.ExtractedData{}, err
	}

	var extracted app.ExtractedData
	if err := json.Unmarshal(content, &extracted); err != nil {
		return app.ExtractedData{}, err
	}

	return extracted, nil
}

func (r *RecordingRepository) recordingFolder(id string) (string, error) {
	folderPath := filepath.Join(r.rootPath, id)
	info, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrNotFound
		}

		return "", err
	}

	if !info.IsDir() {
		return "", ErrNotFound
	}

	return folderPath, nil
}

func countDirectories(entries []os.DirEntry) int {
	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			count++
		}
	}

	return count
}
