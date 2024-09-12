package camunda

import (
	"context"
	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"log"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
)

const (
	jobWorkerDeductCredit     = "deductCredit"
	jobWorkerChargeCreditCard = "chargeCreditCard"
	jobWorkerInvokePayment    = "invokePayment"
	jobWorkerPaymentDone      = "paymentDone"
	kWorkerName               = "go-get-started"
)

func (c *Client) StartFetchDeductCreditWorker() worker.JobWorker {
	client := *c.Client
	w := client.NewJobWorker().
		JobType(jobWorkerDeductCredit).
		Handler(HandleDeductCredit).
		Concurrency(1).
		MaxJobsActive(10).
		RequestTimeout(1 * time.Second).
		PollInterval(1 * time.Second).
		Name(kWorkerName).
		Open()
	log.Printf("started worker [%s] for jobs of type [%s]", kWorkerName, jobWorkerDeductCredit)
	return w
}

func (c *Client) StartFetchPaymentDonetWorker() worker.JobWorker {
	client := *c.Client
	w := client.NewJobWorker().
		JobType(jobWorkerPaymentDone).
		Handler(c.paymentHandler).
		Concurrency(1).
		MaxJobsActive(10).
		RequestTimeout(1 * time.Second).
		PollInterval(1 * time.Second).
		Name(kWorkerName).
		Open()
	log.Printf("started worker [%s] for jobs of type [%s]", kWorkerName, jobWorkerPaymentDone)
	return w
}

func (c *Client) StartFetchChargeCreditWorker() worker.JobWorker {
	client := *c.Client
	w := client.NewJobWorker().
		JobType(jobWorkerChargeCreditCard).
		Handler(HandleChargeCredit).
		Concurrency(1).
		MaxJobsActive(10).
		RequestTimeout(1 * time.Second).
		PollInterval(1 * time.Second).
		Name(kWorkerName).
		Open()
	log.Printf("started worker [%s] for jobs of type [%s]", kWorkerName, jobWorkerChargeCreditCard)
	return w
}

func (c *Client) StartFetchInvokePaymentWorker() worker.JobWorker {
	client := *c.Client
	handleInvokePayment := func(workerClient worker.JobClient, job entities.Job) {
		log.Println("Invoke Payment Worker start with ", job.Key)
		vars, err := job.GetVariablesAsMap()
		if err != nil {
			log.Println(err)
		}
		delete(vars, "token")
		publish, err := client.NewPublishMessageCommand().
			MessageName("paymentRequestedMsg").
			CorrelationKey(vars["orderId"].(string)).
			VariablesFromMap(vars)
		if err != nil {
			log.Println(err)
		}
		send, err := publish.Send(context.Background())
		if err != nil {
			log.Println(err)
		}
		log.Println(send)
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
		log.Printf("handleInvokePayment completed job %d successfully", job.Key)
	}

	w := client.NewJobWorker().
		JobType(jobWorkerInvokePayment).
		Handler(handleInvokePayment).
		Concurrency(1).
		MaxJobsActive(10).
		RequestTimeout(1 * time.Second).
		PollInterval(1 * time.Second).
		Name(kWorkerName).
		Open()

	log.Printf("started worker [%s] for jobs of type [%s]", kWorkerName, jobWorkerInvokePayment)
	return w
}
