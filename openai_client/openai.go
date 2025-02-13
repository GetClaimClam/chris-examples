package openaiClient

import (
	"context"
	"log"
	"time"

	"github.com/chr1sbest/api.mobl.ai/internal/util/errors"

	openai "github.com/sashabaranov/go-openai"
)

const (
	// OpenAI Config
	sleepDuration = 1 * time.Second
	maxAttempts   = 10

	// OpenAI Thread Status
	completedStatus = "completed"
	cancelledStatus = "cancelled"
	failedStatus    = "failed"
	expiredStatus   = "expired"
)

var client *openai.Client

func NewOpenAIClient(apiKey string) *openai.Client {
	return openai.NewClient(apiKey)
}

func GenerateStatelessResponse(ctx context.Context, prompt string, agentID string) (string, error) {
	assistant, err := client.RetrieveAssistant(ctx, agentID)
	if err != nil {
		return "", errors.FailedSetup
	}

	thread, err := client.CreateThread(ctx, openai.ThreadRequest{})
	if err != nil {
		return "", errors.FailedSetup
	}

	_, err = client.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    "user",
		Content: prompt,
	})
	if err != nil {
		return "", errors.FailedSetup
	}

	run, err := client.CreateRun(ctx, thread.ID, openai.RunRequest{AssistantID: assistant.ID})
	if err != nil {
		return "", errors.FailedSetup
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		time.Sleep(sleepDuration)
		runStatus, err := client.RetrieveRun(ctx, thread.ID, run.ID)
		if err != nil {
			return "", err
		}
		status := string(runStatus.Status)
		if status == completedStatus {
			messages, err := client.ListMessage(ctx, thread.ID, nil, nil, nil, nil)
			if err != nil {
				return "", err
			}
			// Success
			if len(messages.Messages) > 0 && len(messages.Messages[0].Content) > 0 {
				response := messages.Messages[0].Content[0].Text.Value
				return response, nil
			}
			return "", errors.NoResponse
		} else if status == cancelledStatus || status == failedStatus || status == expiredStatus {
			log.Errorf("Failed run: %s", status)
			return "", errors.FailedRun
		}
	}

	return "", errors.AITimeout
}

func GenerateStatefulResponse(ctx context.Context, threadID string, prompt string, agentID string) (string, string, error) {
	var thread openai.Thread
	var err error

	if threadID == "" {
		// Create a new thread if threadID is blank
		thread, err = client.CreateThread(ctx, openai.ThreadRequest{})
		if err != nil {
			log.Errorf("Failed create thread: %s", err.Error())
			return "", "", errors.FailedSetup
		}
		threadID = thread.ID
	} else {
		// Retrieve the existing thread
		thread, err = client.RetrieveThread(ctx, threadID)
		if err != nil {
			log.Errorf("Failed retrieve thread: %s", threadID)
			return "", "", errors.FailedSetup
		}
	}

	// Create a new message in the thread
	_, err = client.CreateMessage(ctx, thread.ID, openai.MessageRequest{
		Role:    "user",
		Content: prompt,
	})
	if err != nil {
		log.Errorf("Failed create message: %s", threadID)
		return "", "", errors.FailedSetup
	}

	// Retrieve the assistant
	assistant, err := client.RetrieveAssistant(ctx, agentID)
	if err != nil {
		log.Errorf("Failed retrieve agentID: %s", agentID)
		return "", "", errors.FailedSetup
	}

	// Create a new run
	run, err := client.CreateRun(ctx, thread.ID, openai.RunRequest{AssistantID: assistant.ID})
	if err != nil {
		log.Errorf("Failed create run: %s", err.Error())
		return "", "", errors.FailedSetup
	}

	// Poll for the completion status
	for attempt := 0; attempt < maxAttempts; attempt++ {
		time.Sleep(sleepDuration)
		runStatus, err := client.RetrieveRun(ctx, thread.ID, run.ID)
		if err != nil {
			return threadID, "", err
		}
		status := string(runStatus.Status)
		if status == completedStatus {
			messages, err := client.ListMessage(ctx, thread.ID, nil, nil, nil, nil)
			if err != nil {
				return threadID, "", err
			}
			// Success
			if len(messages.Messages) > 0 && len(messages.Messages[0].Content) > 0 {
				response := messages.Messages[0].Content[0].Text.Value
				return threadID, response, nil
			}
			return threadID, "", errors.NoResponse
		} else if status == cancelledStatus || status == failedStatus || status == expiredStatus {
			log.Errorf("Failed run: %s", status)
			return threadID, "", errors.FailedRun
		}
	}

	// Failed to poll response
	return threadID, "", errors.AITimeout
}
