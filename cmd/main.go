package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsnuyen/go-zebee-worker/internal"
)

func main() {
	// Start the client
	fmt.Println("Start Client on port 9999")
	router := gin.Default()
	router.GET("/ping", internal.PingHandler)
	credit := router.Group("/api/v1")
	{
		credit.POST("/credit/start", internal.StartCreditProcess)
	}
	http.ListenAndServe(":9999", router)
}
