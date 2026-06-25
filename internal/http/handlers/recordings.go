package handlers

import (
	"errors"
	"io"

	"log/slog"

	"github.com/gin-gonic/gin"
)

type RecordingHandler struct {
	uploader recordingUploader
}

type recordingUploader interface {
	UploadRecording(file []byte) (string, error)
}

func NewRecordingHandler(uploader recordingUploader) *RecordingHandler {
	return &RecordingHandler{uploader: uploader}
}

func (h *RecordingHandler) Create(c *gin.Context) {
	file, err := firstMultipartFile(c)
	if err != nil {
		slog.Error("upload recording", "worker", "API", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		slog.Error("upload recording", "worker", "API", "error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, err := h.uploader.UploadRecording(body)
	if err != nil {
		slog.Error("upload recording", "worker", "API", "error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("upload recording", "worker", "API", "id", id, "status", "success")

	c.JSON(201, gin.H{"id": id})
}

func firstMultipartFile(c *gin.Context) (multipartFile, error) {
	if file, _, err := c.Request.FormFile("file"); err == nil {
		return file, nil
	}

	return nil, errors.New("multipart file is required")
}

type multipartFile interface {
	io.Reader
	io.Closer
}
