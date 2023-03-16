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

// create new gin web server
func Run() error {
	e := gin.Default()
	v1 := e.Group("/v1")
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

func init() {
	if err := fake.SetLang("en"); err != nil {
		log.Println(err)
	}
}
