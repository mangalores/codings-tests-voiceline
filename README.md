# Transcription API

## Usage

```text
curl -X POST http://localhost:8080/recordings \
  -F "file=@assets/test.mp3"
```




## Design Decisions

1. File-based persistence

Recording artifacts are stored on the local filesystem. This avoids database setup while keeping the processing pipeline observable. In production, object storage and a database-backed metadata store would be used.

2. In-memory channels

Background workers communicate through Go channels to model asynchronous processing. In production, durable queues (e.g. SQS, Pub/Sub, RabbitMQ) would provide persistence, retries and recovery.

3. Adapter-based integrations

Speech-to-text, extraction and export are abstracted behind interfaces. Mock implementations provide deterministic behaviour, while real adapters (OpenAI, Webhook, Google Sheets) can be enabled through configuration.

4. Mock fixtures

Mock adapters load deterministic responses from the assets directory. This allows the application flow to be tested without external dependencies or API keys.

5. Error handling & status

Asynchronous processing cannot return errors to the initial upload request. Failed processing steps create an error.json artifact, while the status endpoint derives the current processing state from available artifacts.

## Stretchgoal

