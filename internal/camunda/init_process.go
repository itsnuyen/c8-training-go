package camunda

import (
	"context"
	"github.com/google/uuid"
	"github.com/itsnuyen/go-zebee-worker/internal/model"
	"log"
)

func (c *Client) StartCreditProcessWithMessage(credit model.Credit) error {
	client := *c.Client

	correlationKey, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
	}
	credit.OrderId = correlationKey.String()

	request, err := client.NewPublishMessageCommand().
		MessageName("orderReceivedMessage").
		CorrelationKey(correlationKey.String()).
		VariablesFromObject(credit)

	if err != nil {
		return err
	}

	ctx := context.Background()

	msg, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}
	log.Println(msg)
	return nil
}

func (c *Client) StartProcessWithProcessId(credit model.Credit) error {
	client := *c.Client
	request, err := client.NewCreateInstanceCommand().BPMNProcessId("Process_OrderProcess").
		LatestVersion().
		VariablesFromObject(credit)
	if err != nil {
		return err
	}
	ctx := context.Background()
	msg, err := request.Send(ctx)
	if err != nil {
		return err
	}
	log.Println(msg)
	return nil
}
