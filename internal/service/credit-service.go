package service

import (
	"context"
	"fmt"
	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"log"
)

func getCustomerCredit(customerId string) float64 {
	log.Println("Get Customer id")

	balance := float64(len(customerId)) * 10.5
	log.Println("Customer with with balance ", balance)

	return balance
}

func DeductCredit(customerId string, amount float64) float64 {
	balance := getCustomerCredit(customerId)
	result := balance - amount
	return result
}

func ChargeAmount(client worker.JobClient, job entities.Job, cardNumber string, cvc string, expiryDate string, amount float64) {
	log.Println(fmt.Sprintf("Credit Card infos are [CarNumber: %s] [CVC: %s] [ExpiryDate: %s] [Amount: %f]", cardNumber, cvc, expiryDate, amount))
	if len(expiryDate) != 5 {
		log.Println("Error happens because this and stuff check the expiry date", expiryDate)
		//client.NewFailJobCommand().JobKey(job.Key).Retries(0).Send(context.Background())
		client.NewThrowErrorCommand().
			JobKey(job.Key).ErrorCode("errorInvalidData").Send(context.Background())
		return
	}
}
