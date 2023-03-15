package server

import (
	"net/http"
	"strconv"
	"strings"
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
		ID:      strconv.Itoa(int(time.Now().Unix())),
		Object:  "test-object",
		Created: time.Now().Unix(),
		Model:   req.Model,
	}
	// create completions
	for i := 0; i < req.N; i++ {
		// generate a random string of length completionReq.Length
		completionStr := strings.Repeat("a", req.MaxTokens)

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
	c.Header("Content-Type", "text/event-stream")
	// Send test responses
	dataBytes := []byte{}
	dataBytes = append(dataBytes, []byte("event: message\n")...)
	//nolint:lll
	data := `{"id":"1","object":"completion","created":1598069254,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"content":"response1"},"finish_reason":"max_tokens"}]}`
	dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

	dataBytes = append(dataBytes, []byte("event: message\n")...)
	//nolint:lll
	data = `{"id":"2","object":"completion","created":1598069255,"model":"gpt-3.5-turbo","choices":[{"index":0,"delta":{"content":"response2"},"finish_reason":"max_tokens"}]}`
	dataBytes = append(dataBytes, []byte("data: "+data+"\n\n")...)

	dataBytes = append(dataBytes, []byte("event: done\n")...)
	dataBytes = append(dataBytes, []byte("data: [DONE]\n\n")...)

	_, err := c.Writer.Write(dataBytes)
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

// numTokens Returns the number of GPT-3 encoded tokens in the given text.
// This function approximates based on the rule of thumb stated by OpenAI:
// https://beta.openai.com/tokenizer
//
// TODO: implement an actual tokenizer for GPT-3 and Codex (once available)
func numTokens(s string) int {
	return int(float32(len(s)) / 4)
}
