package services

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
)

const (
	pollyVoice    = types.VoiceIdCarmen
	pollyLanguage = types.LanguageCodeRoRo
	pollyEngine   = types.EngineStandard
	pollyFormat   = types.OutputFormatMp3
	maxChunkSize  = 2900
)

type IPollyService interface {
	SynthesizeSpeech(ctx context.Context, text string) ([]byte, error)
}

type PollyService struct {
	client *polly.Client
	logger *slog.Logger
}

func NewPollyService(client *polly.Client, logger *slog.Logger) *PollyService {
	return &PollyService{
		client: client,
		logger: logger,
	}
}

func (s *PollyService) SynthesizeSpeech(ctx context.Context, text string) ([]byte, error) {
	chunks := chunkText(text, maxChunkSize)
	var result bytes.Buffer

	for _, chunk := range chunks {
		output, err := s.client.SynthesizeSpeech(ctx, &polly.SynthesizeSpeechInput{
			Text:         aws.String(chunk),
			VoiceId:      pollyVoice,
			LanguageCode: pollyLanguage,
			Engine:       pollyEngine,
			OutputFormat: pollyFormat,
		})
		if err != nil {
			s.logger.Error("polly synthesis failed",
				slog.Any("error", err),
			)
			return nil, err
		}
		defer output.AudioStream.Close()

		if _, err := io.Copy(&result, output.AudioStream); err != nil {
			return nil, err
		}
	}

	s.logger.Info("speech synthesized",
		slog.Int("chunks", len(chunks)),
		slog.Int("bytes", result.Len()),
	)
	return result.Bytes(), nil
}

func chunkText(text string, size int) []string {
	runes := []rune(text)
	var chunks []string

	for len(runes) > 0 {
		if len(runes) <= size {
			chunks = append(chunks, string(runes))
			break
		}

		// try to cut at sentence boundary
		cutAt := size
		for i := size; i > size-200; i-- {
			if runes[i] == '.' || runes[i] == '!' || runes[i] == '?' || runes[i] == '\n' {
				cutAt = i + 1
				break
			}
		}

		chunks = append(chunks, string(runes[:cutAt]))
		runes = runes[cutAt:]
	}

	return chunks
}
