package handlers

import (
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"

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
	file, filename, err := firstMultipartFile(c)
	if err != nil {
		slog.Error("upload recording", "worker", "API", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.ToLower(filepath.Ext(filename)) != ".mp3" {
		err := errors.New("only mp3 files are supported")
		slog.Error("upload recording", "worker", "API", "error", err)
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "only mp3 files are supported"})
		return
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		slog.Error("upload recording", "worker", "API", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := h.uploader.UploadRecording(body)
	if err != nil {
		slog.Error("upload recording", "worker", "API", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("upload recording", "worker", "API", "id", id, "status", "success")

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func firstMultipartFile(c *gin.Context) (multipartFile, string, error) {
	if file, header, err := c.Request.FormFile("file"); err == nil {
		return file, header.Filename, nil
	}

	return nil, "", errors.New("multipart file is required")
}

type multipartFile interface {
	io.Reader
	io.Closer
}
