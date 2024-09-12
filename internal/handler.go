package internal

import (
	"github.com/itsnuyen/go-zebee-worker/internal/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartCreditProcess(context *gin.Context) {
	log.Print("Start Credit Process")
	var creditModel model.Credit
	if err := context.BindJSON(&creditModel); err != nil {
		log.Println("Could not decompose json because of err", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	log.Println("current model", creditModel)
	err := camundaClient.StartCreditProcessWithMessage(creditModel)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	context.JSON(http.StatusOK, gin.H{"message": "Credit Process Started"})
}

func PingHandler(context *gin.Context) {
	context.JSON(200, gin.H{"ping": "pong"})
}
