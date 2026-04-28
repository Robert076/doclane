package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const modelID = "eu.amazon.nova-lite-v1:0"

type IBedrockService interface {
	InterpretDocument(ctx context.Context, extractedText string, documentTitle string) (string, error)
}

type BedrockService struct {
	client *bedrockruntime.Client
	logger *slog.Logger
}

func NewBedrockService(client *bedrockruntime.Client, logger *slog.Logger) *BedrockService {
	return &BedrockService{
		client: client,
		logger: logger,
	}
}

type novaResponse struct {
	Output struct {
		Message struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"message"`
	} `json:"output"`
}

func (s *BedrockService) InterpretDocument(ctx context.Context, extractedText string, documentTitle string) (string, error) {
	prompt := fmt.Sprintf(`Ești un asistent pentru funcționari publici români. Ai primit textul extras dintr-un document de tip "%s".

Textul documentului:
---
%s
---

Oferă:
1. Un rezumat scurt și clar al documentului (2-3 propoziții)
2. Informații cheie identificate (nume, date, sume, termene)
3. Dacă documentul pare să corespundă tipului "%s" (da/nu și motivul)
4. Orice probleme sau neclarități observate

Sa tii cont si de faptul ca data curenta este %v.

Răspunde în română, concis și structurat.`, documentTitle, extractedText, documentTitle, time.Now())

	body, err := json.Marshal(map[string]any{
		"messages": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{"text": prompt},
				},
			},
		},
		"inferenceConfig": map[string]any{
			"maxTokens": 1024,
		},
	})
	if err != nil {
		return "", err
	}

	output, err := s.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		s.logger.Error("bedrock invocation failed",
			slog.String("model", modelID),
			slog.Any("error", err),
		)
		return "", err
	}

	var resp novaResponse
	if err := json.Unmarshal(output.Body, &resp); err != nil {
		return "", err
	}

	if len(resp.Output.Message.Content) == 0 {
		return "", fmt.Errorf("empty response from bedrock")
	}

	return resp.Output.Message.Content[0].Text, nil
}
