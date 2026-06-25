package extractor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/mangalores/case-studies-voiceline/internal/app"
)

const openAIResponsesURL = "https://api.openai.com/v1/responses"

type OpenAIExtractor struct {
	APIKey string
}

func NewOpenAIExtractor(apiKey string) *OpenAIExtractor {
	return &OpenAIExtractor{APIKey: apiKey}
}

func (e *OpenAIExtractor) Extract(ctx context.Context, transcription string) (app.ExtractedData, error) {
	if strings.TrimSpace(e.APIKey) == "" {
		return app.ExtractedData{}, fmt.Errorf("openai api key is required")
	}

	payload, err := buildAnalysisRequest(transcription)
	if err != nil {
		return app.ExtractedData{}, err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return app.ExtractedData{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		openAIResponsesURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return app.ExtractedData{}, err
	}

	req.Header.Set("Authorization", "Bearer "+e.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return app.ExtractedData{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return app.ExtractedData{}, err
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		return app.ExtractedData{}, fmt.Errorf("openai error: %s", strings.TrimSpace(string(respBody)))
	}

	var raw struct {
		Output []struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}

	if err := json.Unmarshal(respBody, &raw); err != nil {
		return app.ExtractedData{}, err
	}

	if len(raw.Output) == 0 || len(raw.Output[0].Content) == 0 {
		return app.ExtractedData{}, fmt.Errorf("empty OpenAI response")
	}

	var result app.ExtractedData
	if err := json.Unmarshal([]byte(raw.Output[0].Content[0].Text), &result); err != nil {
		return app.ExtractedData{}, err
	}

	slog.Info("OpenAIExtractor extracted data", "participants", len(result.Participants), "actionItems", len(result.ActionItems))

	return result, nil

}

func buildAnalysisRequest(transcript string) (map[string]any, error) {
	return map[string]any{
		"model": "gpt-4.1-mini",
		"input": []map[string]any{
			{
				"role": "system",
				"content": `You analyze meeting transcripts.
							Tasks:
							1. Create a concise summary.
							2. Extract all participants explicitly mentioned.
							3. Extract decisions made during the meeting.
							4. Extract action items from the transcript concerning future tasks (e.g. meetings). Action items are tasks when a person is responsible for doing something, even if the sentence does not use the words "action item".

							Rules:
							- Do not invent names.
							- If a due date is not mentioned, leave it empty.
							- If a participant is not mentioned, leave it empty.
							`,
			},
			{
				"role":    "user",
				"content": "Create a concise summary, extract participant names and action items from this transcript:\n\n" + transcript,
			},
		},
		"text": map[string]any{
			"format": map[string]any{
				"type":   "json_schema",
				"name":   "meeting_analysis",
				"strict": true,
				"schema": map[string]any{
					"type":                 "object",
					"additionalProperties": false,
					"required":             []string{"summary", "participants", "decisions", "action_items"},
					"properties": map[string]any{
						"summary": map[string]any{
							"type": "string",
						},
						"participants": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "string",
							},
						},
						"decisions": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "string",
							},
						},
						"action_items": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type":                 "object",
								"additionalProperties": false,
								"required": []string{
									"owner",
									"task",
									"due",
								},
								"properties": map[string]any{
									"owner": map[string]any{
										"type": "string",
									},
									"task": map[string]any{
										"type": "string",
									},
									"due": map[string]any{
										"type": "string",
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}
