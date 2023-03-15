package main

import (
	"fmt"

	"github.com/venjiang/openai-server/server"
)

func main() {
	addr := ":9080"
	fmt.Printf("OpenAI GPT Mock Server running: %s\n", addr)
	if err := server.Run(addr); err != nil {
		panic(err)
	}
}
