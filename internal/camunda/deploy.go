package camunda

import (
	"context"
	"errors"
	"fmt"
	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
	"log"
	"os"
	"time"
)

const CorrectDataForm = "correctData.form"
const OrderDiscountDmn = "orderDiscount.dmn"
const OrderProcessBpmn = "orderProcess.bpmn"
const PaymentProcessBpmn = "paymentProcess.bpmn"

func (c *Client) DeployFormDefinition(definitionName string) {
	client := *c.Client
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	resource, err := client.NewDeployResourceCommand().
		AddResourceFile(fmt.Sprintf("./bpmn/%s", definitionName)).
		Send(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resource)
}

func (c *Client) DeployProcessDefinition(definitionName string) *pb.ProcessMetadata {
	client := *c.Client
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	resource, err := client.NewDeployResourceCommand().
		AddResourceFile(fmt.Sprintf("./bpmn/%s", definitionName)).
		Send(ctx)

	if err != nil {
		panic(err)
	}

	if len(resource.GetDeployments()) < 0 {
		panic(errors.New("failed to deploy send-email model; nothing was deployed"))
	}

	deployment := resource.GetDeployments()[0]
	process := deployment.GetProcess()
	if process == nil {
		panic(errors.New("failed to deploy process; the deployment was successful, but no process was returned"))
	}

	log.Printf("deployed BPMN process [%s] with key [%d]", process.GetBpmnProcessId(), process.GetProcessDefinitionKey())
	return process
}

func ReadBPMNFile(resourceFile string) []byte {
	pwd, _ := os.Getwd()
	path := fmt.Sprintf("%s/bpmn/%s", pwd, resourceFile)
	fmt.Println(path)
	file, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	fmt.Println("The size of the file is ", len(file))

	return file
}
