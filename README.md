# Audio Transcription API (mp3)

A simple audio upload API which in the background transcribes the uploaded mp3 file,
extracts useful information into a normalized data structure and pushes to an export target.

## Environment

The service loads its configuration from environment variables.

If no `.env` file is present, the application falls back to its default configuration, which uses mock adapters and the test assets located in `assets/`. This allows the complete processing pipeline to be executed without external dependencies or API keys.

A `.env.dist` file is provided as a template for configuring external services.

## Configuration

For local development, copy the example file:

```bash
cp .env.dist .env
```

using default env is configued with mock adapters. This should work out of the box


## Build

Requires Go 1.24 or newer.

To build execute from the project root:

```bash
    go mod vendor
    go build -o api ./cmd/api
```

## Usage

After build start web server:

```bash
    ./api
```

### Using the Default Mock Configuration


```bash
curl -X POST http://localhost:8080/recordings \
  -F "file=@assets/test_recording.mp3"
```

### With OpenAI Transcriber,Extractor & Webhook Export

Configure the OpenAI-based adapters and point the webhook exporter to your endpoint (e.g. webhook.site)::

```text
TRANSCRIBER=openai
EXTRACTOR=openai
EXPORTER=webhook
OPENAI_API_KEY=<YOUR_API_KEY>
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

- googlesheet (exists as option, but implementation was not finished due to time constraints)

## File Storage
Generated artifacts are written to:

```text
storage/<recording-id>/
    recording.mp3
    transcription.md
    extracted.json
``` 

## Design Decisions

### 1. File-based persistence

Recording artifacts are stored on the local filesystem. This avoids database setup while keeping the processing pipeline observable. In production, object storage and a database-backed metadata store would be a better fit.

### 2. In-memory channels

Background workers communicate through Go channels to model asynchronous processing. In production, durable queues such as SQS, Pub/Sub, or RabbitMQ would provide persistence, retries, and recovery. 

This choice keeps the API responsive while transcription and extraction are in progress as the LLM api calls + export may take time to process. The tradeoff is limited throughput and only viable in this proof of concept.

### 3. Adapter-based integrations

Speech-to-text, extraction, and export are abstracted behind interfaces. Mock implementations provide deterministic behaviour, while real adapters such as OpenAI and Webhook can be enabled through configuration.

### 4. Mock fixtures

Mock adapters load deterministic responses from the assets directory. This allows the application flow to be exercised without external dependencies or API keys. It allowed ensuring complete processing flow independent of the implementation states.

### 5. Error handling of workers

Asynchronous processing cannot return errors to the initial upload request. Failed processing steps create an `error.json` artifact and are also reported in the logs. This should be later served via a status endpoint to allow polling for process updates.
