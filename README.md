# Audio Transcription API (mp3)

A simple audio upload API which in the background transcribes the uploaded mp3 file,
extracts useful information into a normalized data structure and pushes to an export target.

## Environment

The service uses .env to load environment variables. If no .env file is present, the application falls back to default configuration using the mock adapters and the test assets located in assets/. This allows the complete processing pipeline to be executed without external dependencies or API keys.

A `.env.dist` is provided as template for adding required API Keys or settings.

## File Storage
Generated artifacts are written to:

storage/<recording-id>/
    recording.mp3
    transcription.md
    extracted.json

## Build

Go 1.26.4 was used for implementation
From the project root:

```bash
    go mod vendor
    go build -o api ./cmd/api
```

## Usage

After build start web server:

```bash
    ./api
```

### With Mock Adapters

using default env with mock adapters. Should work out of the box

```bash
curl -X POST http://localhost:8080/recordings \
  -F "file=@assets/test_recording.mp3"
```

### With OpenAI Transcriber,Extractor & Webhook Export

Set the exporter to webhook and point it at your endpoint:

```text
TRANSCRIBER=openai
EXTRACTOR=openai
EXPORTER=webhook
OPENAI_API_KEY=<
WEBHOOK_EXPORT_URL=<YOUR_EXPORT_WEBHOOK_URL>
```

The webhook receives a `POST` request with the extracted data as JSON:

```json
{
  "summary": "...",
  "participants": ["..."],
  "decisions": ["..."],
  "actionItems": [
    {
      "owner": "...",
      "task": "...",
      "due": "..."
    }
  ]
}
```


## Supported Adapters

The ability to have various adapters for the different processing steps is represented and the concrete implementations are:

### Transcriber
- mock
- openai

### Extractor
- mock
- openai

### Exporter
- mock
- webhook

- googlesheet (exists as option, but implementation was not finished due to time constrains)

## Design Decisions

### 1. File-based persistence

Recording artifacts are stored on the local filesystem. This avoids database setup while keeping the processing pipeline observable. In production, object storage and a database-backed metadata store would be used.

### 2. In-memory channels

Background workers communicate through Go channels to model asynchronous processing. In production, durable queues (e.g. SQS, Pub/Sub, RabbitMQ) would provide persistence, retries and recovery. This was chosen to keep the API response fast given delays in actual transcription and extraction calls.

### 3. Adapter-based integrations

Speech-to-text, extraction and export are abstracted behind interfaces. Mock implementations provide deterministic behaviour, while real adapters (OpenAI, Webhook) can be enabled through configuration.

### 4. Mock fixtures

Mock adapters load deterministic responses from the assets directory. This allows the application flow to be tested without external dependencies or API keys.

### 5. Error handling of workers

Asynchronous processing cannot return errors to the initial upload request. Failed processing steps create an error.json artifact and will report error in logging output
