package app

type RecordingStore interface {
	StoreRecording(file []byte) (string, error)
}

type UploadService struct {
	recordings         RecordingStorer
	transcribeCommands chan<- TranscribeCommand
}

func NewUploadService(recordings RecordingStorer, transcribeCommands chan<- TranscribeCommand) *UploadService {
	return &UploadService{recordings: recordings, transcribeCommands: transcribeCommands}
}

func (s *UploadService) UploadRecording(file []byte) (string, error) {
	id, err := s.recordings.StoreRecording(file)
	if err != nil {
		return "", err
	}

	if s.transcribeCommands != nil {
		s.transcribeCommands <- TranscribeCommand{ID: id}
	}

	return id, nil
}
