package server

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
)

func ChatCompletionHandler(c *gin.Context) {
	var req openai.ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// stream
	if req.Stream {
		handleChatCompletionStream(&req, c)
		return
	}
	handleChatCompletion(&req, c)
}

func handleChatCompletion(req *openai.ChatCompletionRequest, c *gin.Context) {
	res := openai.ChatCompletionResponse{
		ID:      PrefixID("chatcmpl-"),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
	}
	rand.Seed(time.Now().UnixNano())
	if req.N == 0 {
		req.N = rand.Intn(3)
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = rand.Intn(64)
	}
	for i := 0; i < req.N; i++ {
		completionStr := RandomContent()

		res.Choices = append(res.Choices, openai.ChatCompletionChoice{
			Message: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completionStr,
			},
			Index: i,
		})
	}
	inputTokens := numTokens(req.Messages[0].Content) * req.N
	completionTokens := req.MaxTokens * req.N
	res.Usage = openai.Usage{
		PromptTokens:     inputTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      inputTokens + completionTokens,
	}
	c.JSON(http.StatusOK, res)
}

func handleChatCompletionStream(req *openai.ChatCompletionRequest, c *gin.Context) {
	resChan := make(chan *openai.ChatCompletionStreamResponse)
	go func() {
		defer close(resChan)
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(20)
		id := PrefixID("chatcmpl-")
		for i := 0; i < n; i++ {
			delay := rand.Int63n(300) * 1e6
			time.Sleep(time.Duration(delay))
			res := openai.ChatCompletionStreamResponse{
				ID:      id,
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   req.Model,
			}
			content := RandomContent()
			choise := openai.ChatCompletionStreamChoice{
				Delta: openai.ChatCompletionStreamChoiceDelta{
					Content: content,
				},
			}
			if i == n-1 {
				choise.Delta = openai.ChatCompletionStreamChoiceDelta{}
				choise.FinishReason = "stop"
			}
			res.Choices = append(res.Choices, choise)
			resChan <- &res
		}
	}()
	// stream
	SetStreamHeaders(c)
	c.Stream(func(w io.Writer) bool {
		data := []byte("data: ")
		// chunk data
		if res, ok := <-resChan; ok {
			chunk, err := json.Marshal(res)
			if err != nil {
				w.Write([]byte("data: [ERROR]\n\n"))
				return false
			}
			// write
			data = append(data, chunk...)
			data = append(data, []byte("\n\n")...)
			_, err = w.Write(data)
			if err != nil {
				w.Write([]byte("data: [ERROR]\n\n"))
				return false
			}
			return true
		}
		// done
		w.Write([]byte("data: [DONE]\n\n"))
		return false
	})
}

func SetStreamHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	// c.Header("openai-model", req.Model)
	// c.Header("openai-organization", "yomo")
	// c.Header("openai-processing-ms", ms)
	// c.Header("openai-version", "mock")
	// c.Header("x-request-id", req.ID)
}

// numTokens Returns the number of GPT-3 encoded tokens in the given text.
// This function approximates based on the rule of thumb stated by OpenAI:
// https://beta.openai.com/tokenizer
//
// TODO: implement an actual tokenizer for GPT-3 and Codex (once available)
func numTokens(s string) int {
	return int(float32(len(s)) / 4)
}
