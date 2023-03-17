package server

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/icrowley/fake"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var fns = []func() string{
	fake.UserName,
	fake.Word,
	fake.Words,
	fake.Sentence,
	fake.Sentences,
	fake.ProductName,
	fake.Brand,
	fake.City,
	fake.Color,
	fake.Title,
	fake.Street,
	fake.Country,
	fake.City,
	fake.Company,
}

func init() {
	if err := fake.SetLang("en"); err != nil {
		log.Println(err)
	}
}

// create new gin web server
func Run() error {
	e := gin.Default()
	v1 := e.Group("/v1")
	v1.POST("/completions", CompletionHandler)
	v1.POST("/chat/completions", ChatCompletionHandler)
	return e.Run()
}

func RandomContent() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(fns))
	return fns[n]()
}

func ID(length ...int) string {
	return PrefixID("", length...)
}

func PrefixID(prefix string, length ...int) string {
	l := 29
	if len(length) > 0 {
		l = length[0]
	}
	id, err := gonanoid.Generate(alphabet, l)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s%s", prefix, id)
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
