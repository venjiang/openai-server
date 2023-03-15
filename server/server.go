package server

import "github.com/gin-gonic/gin"

// create new gin web server
func Run() error {
	e := gin.Default()
	v1 := e.Group("/v1")
	v1.POST("/chat/completions", ChatCompletionHandler)
	return e.Run()
}
