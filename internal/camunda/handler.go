package camunda

import (
	"context"
	"fmt"
	"github.com/itsnuyen/go-zebee-worker/internal/service"
	"log"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
)

func HandleChargeCredit(client worker.JobClient, job entities.Job) {
	vars, err := job.GetVariablesAsMap()
	if err != nil {
		log.Printf("failed to get variables for job %d: [%s]", job.Key, err)
		return
	}
	openAmount := vars["openAmount"].(float64)
	creditCard := vars["creditCard"].(map[string]interface{})
	service.ChargeAmount(client, job, creditCard["cardNumber"].(string), creditCard["cvc"].(string), creditCard["expiryDate"].(string), openAmount)
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	newCommand, err := client.NewCompleteJobCommand().
		JobKey(job.Key).
		VariablesFromMap(map[string]interface{}{
			"credit": "charged",
		})
	if err != nil {
		log.Printf("failed to complete job with key %d: [%s]", job.Key, err)
	}
	// Send the command to the broker, need to go to the next step
	newCommand.Send(ctx)
	log.Printf("HandleChargeCredit completed job %d successfully", job.Key)
}

func HandleDeductCredit(client worker.JobClient, job entities.Job) {
	vars, err := job.GetVariablesAsMap()
	if err != nil {
		log.Printf("failed to get variables for job %d: [%s]", job.Key, err)
		return
	}
	tmp := ""
	for k, v := range vars {
		tmp = tmp + fmt.Sprintf("[%s: %s]", k, v)
	}
	log.Printf("Charge Credit Card with: %s", tmp)
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	customerId := vars["customerId"].(string)
	//orderTotal := vars["orderTotal"].(float64)
	discountedAmount := vars["discountedAmount"].(float64)
	balance := 0.0
	result := service.DeductCredit(customerId, discountedAmount)
	status := "sufficient"
	if result < 0 {
		status = "insufficient"
		log.Println("Credit result is ", result)
	} else {
		balance = result
	}
	newCommand, err := client.NewCompleteJobCommand().
		JobKey(job.Key).
		VariablesFromMap(map[string]interface{}{
			"status":     status,
			"balance":    balance,
			"openAmount": result * -1,
		})
	if err != nil {
		log.Printf("failed to complete job with key %d: [%s]", job.Key, err)
	}
	// Send the command to the broker, need to go to the next step
	newCommand.Send(ctx)
	log.Printf("HandleDeductCredit completed job %d successfully", job.Key)
}

func (c *Client) paymentHandler(workerClient worker.JobClient, job entities.Job) {
	log.Println("Payment Worker start with ", job.Key)
	client := *c.Client
	vars, err := job.GetVariablesAsMap()
	if err != nil {
		log.Println(err)
	}
	publish, err := client.NewPublishMessageCommand().
		MessageName("paymentReceivedMessage").
		CorrelationKey(vars["orderId"].(string)).
		VariablesFromMap(vars)
	if err != nil {
		log.Println(err)
	}
	send, err := publish.Send(context.Background())
	if err != nil {
		log.Println(err)
	}
	log.Println(fmt.Sprintf("Send payment infos %s with maps %v", send, vars))

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	newCommand, err := workerClient.NewCompleteJobCommand().
		JobKey(job.Key).
		VariablesFromMap(map[string]interface{}{})
	if err != nil {
		log.Printf("failed to complete job with key %d: [%s]", job.Key, err)
	}
	// Send the command to the broker, need to go to the next step
	newCommand.Send(ctx)
	log.Printf("paymentHandler completed job %d successfully and send out message for paymentReceivedMessage", job.Key)
}
